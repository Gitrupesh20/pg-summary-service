package handler

import (
	"fmt"
	"net/http"
	service2 "pg-summary-service/internal/service"
)

type AuthType string

const (
	AuthTypeNone  AuthType = "none"
	AuthTypeBasic AuthType = "basic"
	AuthTypeToken AuthType = "token"
)

type Route struct {
	Path     string
	Method   string
	Handler  func(http.ResponseWriter, *http.Request)
	AuthType AuthType
}

var service service2.SummaryService

func RegisterRoutes(handler func(pattern string, handler func(http.ResponseWriter, *http.Request)), s service2.SummaryService) {
	service = s
	for i, route := range routes {
		fmt.Println("ith ", i)
		handler(route.Path, func(method string, handler http.HandlerFunc, authType AuthType) http.HandlerFunc {
			fmt.Println("Registering route:", route.Path, "with method:", method, "and auth type:", authType)

			// Here you can add middleware for authentication based on authType

			return ApplyMiddlewares(route.Method, route.AuthType, handler)
		}(route.Method, route.Handler, route.AuthType))
	}
	// Catch-all 404 handler for unknown paths
	handler("/", NotFoundHandler)
}

var routes = []Route{
	{
		Path:    "/summary/sync",
		Method:  http.MethodPost,
		Handler: SyncSummaryHandler,
	},
	{
		Path:    "/summaries",
		Method:  http.MethodGet,
		Handler: GetSummariesHandler,
	},
	{
		Path:    "/summaries/",
		Method:  http.MethodGet,
		Handler: GetSummaryByIDHandler,
	},
}
