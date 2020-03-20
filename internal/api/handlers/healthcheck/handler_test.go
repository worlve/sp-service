package healthcheckhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/worlve/sp-service/internal/api"
	"github.com/worlve/sp-service/internal/api/handlers/handlertestutils"
	"github.com/worlve/sp-service/internal/api/handlers/healthcheck/mocks"
	"github.com/stretchr/testify/mock"
)

// @NOTE: I'm doing this here since we need a handler to test it with and the healthcheck one
// seems like the best choice.
func TestMethodNotAllowed(t *testing.T) {
	cases := []struct {
		name                 string
		method               string
		endpoint             string
		headers              map[string]string
		authN                api.AuthN
		authZ                api.AuthZ
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name:     "method not allowed",
			method:   http.MethodPost,
			endpoint: "healthcheck",
			authN: api.AuthN{
				Datacenter:      "LOCAL",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"405 - Method Not Allowed\",\"message\":\"method not allowed\"}}\n",
			expectedStatusCode:   405,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			healthcheckService := new(mocks.HealthcheckService)
			routerHandlers := HealthcheckRouterHandlers(tc.authZ.APIPath, healthcheckService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         tc.method,
				Endpoint:       tc.endpoint,
				Headers:        tc.headers,
				Body:           nil,
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}

type isHealthyCall struct {
	returnIsHealthy bool
	returnErr       error
}

func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name                 string
		headers              map[string]string
		authN                api.AuthN
		authZ                api.AuthZ
		expectedResponseBody string
		expectedStatusCode   int
		isHealthyCalls       []isHealthyCall
	}{
		{
			name:                 "not authenticated",
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name: "bad admin secret",
			headers: map[string]string{
				"X-ADMIN-AUTH-SECRET": "BAD_SECRET",
			},
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name:                 "happy healthy healthcheck, local",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			isHealthyCalls:       []isHealthyCall{{returnIsHealthy: true}},
		},
		{
			name: "happy healthy healthcheck, prod",
			headers: map[string]string{
				"X-ADMIN-AUTH-SECRET": "SECRET",
			},
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			isHealthyCalls:       []isHealthyCall{{returnIsHealthy: true}},
		},
		{
			name:                 "bad healthcheck, local",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"status\":\"error\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			isHealthyCalls:       []isHealthyCall{{returnIsHealthy: false}},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckService := new(mocks.HealthcheckService)
			for index := range tc.isHealthyCalls {
				healthcheckService.On("IsHealthy", mock.Anything).Return(tc.isHealthyCalls[index].returnIsHealthy, tc.isHealthyCalls[index].returnErr)
			}
			routerHandlers := HealthcheckRouterHandlers(tc.authZ.APIPath, healthcheckService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodGet,
				Endpoint:       "healthcheck",
				Headers:        tc.headers,
				Body:           nil,
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			healthcheckService.AssertNumberOfCalls(t, "IsHealthy", len(tc.isHealthyCalls))
		})
	}
}
