package pagehandler

import (
	"fmt"
	"net/http"

	"github.com/worlve/sp-service/internal/api"
)

// HTTP path fragments keys
const (
	PageIDRouteKey = "pageID"
)

// PageRouterHandlers returns the requests for the associated routes.
func PageRouterHandlers(apiPath string, pageService PageService) []api.RouterHandler {
	handler := PageHandler{
		PageService: pageService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPost,
		Endpoint: fmt.Sprintf("/%v/pages", apiPath),
		Handle:   handler.CreatePage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPatch,
		Endpoint: fmt.Sprintf("/%v/pages/:%v", apiPath, PageIDRouteKey),
		Handle:   handler.UpdatePage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodDelete,
		Endpoint: fmt.Sprintf("/%v/pages/:%v", apiPath, PageIDRouteKey),
		Handle:   handler.DeletePage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/pages", apiPath),
		Handle:   handler.GetPages,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/pages/:%v", apiPath, PageIDRouteKey),
		Handle:   handler.GetPage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/pages/:%v/full", apiPath, PageIDRouteKey),
		Handle:   handler.GetEntirePage,
	})
	return routerHandlers
}
