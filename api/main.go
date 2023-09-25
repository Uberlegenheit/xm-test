package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"xm-task/conf"
	"xm-task/log"
	"xm-task/services"
)

type (
	API struct {
		router   *gin.Engine
		server   *http.Server
		cfg      conf.Config
		services services.Service
	}

	// Route stores an API route data.
	Route struct {
		Path   string
		Method string
		Func   func(http.ResponseWriter, *http.Request)
	}
)

func NewAPI(cfg conf.Config, s services.Service) (*API, error) {
	api := &API{
		cfg:      cfg,
		services: s,
	}

	api.initialize()
	return api, nil
}

// Run starts the http server and binds the handlers.
func (api *API) Run() error {
	return api.startServe()
}

func (api *API) Stop() error {
	return api.server.Shutdown(context.Background())
}

func (api *API) Title() string {
	return "API"
}

func (api *API) initialize() {
	api.router = gin.Default()

	// api.router.Use(gin.Logger())
	api.router.Use(gin.Recovery())

	api.router.Use(cors.New(cors.Config{
		AllowOrigins:     api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodPost, http.MethodHead, http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete,
		},
		AllowHeaders: []string{
			"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "User-Env", "Access-Control-Request-Headers", "Access-Control-Request-Method",
		},
	}))

	// public routes
	api.router.GET("/", api.Index)
	api.router.GET("/health", api.Health)

	api.router.GET("/companies/:id", api.GetCompany)
	api.router.POST("/sign-in", api.SignIn)
	api.router.POST("/refresh", api.Refresh)

	// protected routes
	authGroup := api.router.Group("/auth")
	authGroup.Use(api.AuthMiddleware())
	{
		authGroup.POST("/companies", api.CreateCompany)
		authGroup.PATCH("/companies/:id", api.PatchCompany)
		authGroup.DELETE("/companies/:id", api.DeleteCompany)

		authGroup.POST("/logout", api.LogOut)
	}

	api.server = &http.Server{Addr: fmt.Sprintf(":%d", api.cfg.API.ListenOnPort), Handler: api.router}
}

func (api *API) startServe() error {
	log.Info("Start listening server on port", zap.Uint64("port", api.cfg.API.ListenOnPort))
	err := api.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Warn("API server was closed")
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot run API service: %s", err.Error())
	}
	return nil
}
