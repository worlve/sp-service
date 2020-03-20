package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/worlve/sp-service/internal/api"
	healthcheckhandler "github.com/worlve/sp-service/internal/api/handlers/healthcheck"
	pagehandler "github.com/worlve/sp-service/internal/api/handlers/page"
	pagedetailhandler "github.com/worlve/sp-service/internal/api/handlers/pagedetail"
	healthcheckservice "github.com/worlve/sp-service/internal/services/healthcheck"
	pageservice "github.com/worlve/sp-service/internal/services/page"
	pagedetailservice "github.com/worlve/sp-service/internal/services/pagedetail"
	"github.com/worlve/sp-service/internal/stores/mysqlstore"
	"github.com/worlve/sp-service/internal/util/env"
)

const localUIURL = "http://127.0.0.1:8081"

const (
	defaultAdminAuthSecret = "DEFAULT_SECRET"
	defaultPort            = "8782"
	defaultStaticPath      = "../../static"
	defaultDatacenter      = "LOCAL"
)

func getHTTPServerAddr() string {
	port := env.Get("PORT", defaultPort)
	return ":" + port
}

func getHTTPServerReadTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerWriteTimeout() time.Duration {
	return 10 * time.Second
}

func getHTTPServerMaxHeaderBytes() int {
	return 1 << 20
}

func getAPIPath() string {
	return "api"
}

func getStaticPath() string {
	return env.Get("STATIC_PATH", defaultStaticPath)
}

func getDatacenter() string {
	return env.Get("DATACENTER", api.LocalDatacenterEnv)
}

func main() {
	mysqldb, err := mysqlstore.SetupMySQL("")
	if err != nil {
		fmt.Printf("Failed to connect to MySQL db.\nIf connecting locally, follow https://github.com/worlve/sp-database/blob/master/README.md to get the local db running.\n")
		log.Fatal(err)
	}
	defer mysqldb.Close()
	apiPath := getAPIPath()
	staticPath := getStaticPath()
	datacenter := getDatacenter()
	handler, err := setupHandler(apiPath, staticPath, datacenter, mysqldb)
	if err != nil {
		log.Fatal(err)
	}
	handler, err = setupCors(datacenter, handler)
	if err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Addr:           getHTTPServerAddr(),
		Handler:        handler,
		ReadTimeout:    getHTTPServerReadTimeout(),
		WriteTimeout:   getHTTPServerWriteTimeout(),
		MaxHeaderBytes: getHTTPServerMaxHeaderBytes(),
	}
	fmt.Printf("Starting server at http://localhost%v\nVerify locally by running:\ncurl -X GET http://localhost%v/%v/healthcheck\nAPI docs: http://localhost%v/%v/docs\n", getHTTPServerAddr(), getHTTPServerAddr(), getAPIPath(), getHTTPServerAddr(), getAPIPath())
	log.Fatal(s.ListenAndServe())
}

func setupHandler(apiPath, staticPath, datacenter string, mysqldb *sql.DB) (http.Handler, error) {
	var handler http.Handler
	pageStore := mysqlstore.NewPageStore(mysqldb)
	userStore := mysqlstore.NewUserStore(mysqldb)
	healthcheckStore := mysqlstore.NewHealthcheckStore(mysqldb)
	pageTemplateStore := mysqlstore.NewPageTemplateStore(mysqldb)
	versionStore := mysqlstore.NewVersionStore(mysqldb)
	pageService := pageservice.PageService{
		PageStore:         pageStore,
		PageTemplateStore: pageTemplateStore,
		VersionStore:      versionStore,
		UserStore:         userStore,
	}
	pageDetailService := pagedetailservice.PageDetailService{
		PageDetailStore: pageDetailStore,
	}
	healthcheckService := healthcheckservice.HealthcheckService{
		HealthcheckStore: healthcheckStore,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, pagehandler.PageRouterHandlers(apiPath, pageService)...)
	routerHandlers = append(routerHandlers, pagedetailhandler.PageDetailRouterHandlers(apiPath, pageDetailService)...)
	routerHandlers = append(routerHandlers, healthcheckhandler.HealthcheckRouterHandlers(apiPath, healthcheckService)...)
	router := api.NewRouter(apiPath, staticPath, routerHandlers)
	authN, authZ, err := getAuths(apiPath, datacenter)
	if err != nil {
		return handler, err
	}
	return &api.Handler{
		AuthN:      authN,
		AuthZ:      authZ,
		Router:     router,
		Datacenter: datacenter,
		APIPath:    apiPath,
	}, nil
}

func getAuths(apiPath, datacenter string) (api.AuthN, api.AuthZ, error) {
	adminAuthSecret, err := getAdminAuthSecret(datacenter)
	if err != nil {
		return api.AuthN{}, api.AuthZ{}, err
	}
	authN := api.AuthN{
		Datacenter:      datacenter,
		AdminAuthSecret: adminAuthSecret,
	}
	authZ := api.AuthZ{
		APIPath: apiPath,
	}
	return authN, authZ, nil
}

func getAdminAuthSecret(datacenter string) (string, error) {
	if datacenter != api.LocalDatacenterEnv {
		return env.Require("ADMIN_AUTH_SECRET")
	}
	return env.Get("ADMIN_AUTH_SECRET", defaultAdminAuthSecret), nil
}

func setupCors(datacenter string, handler http.Handler) (http.Handler, error) {
	if datacenter != api.LocalDatacenterEnv {
		return handler, nil
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{localUIURL},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"X-AUTH-TOKEN", "Content-Type", "X-USER-ID"},
	})
	return c.Handler(handler), nil
}
