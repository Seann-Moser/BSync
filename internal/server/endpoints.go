package server

import (
	"context"
	"github.com/Seann-Moser/BaseGoAPI/internal/configuration"
	"github.com/Seann-Moser/BaseGoAPI/internal/middleware"
	"github.com/Seann-Moser/BaseGoAPI/pkg/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type EndpointManager struct {
	conf            *configuration.Config
	router          *mux.Router
	ctx             context.Context
	responseManager *response.Response
}

func NewEndpoints(ctx context.Context, conf *configuration.Config) *EndpointManager {
	em := &EndpointManager{
		conf:            conf,
		router:          mux.NewRouter(),
		ctx:             ctx,
		responseManager: response.NewResponse(conf.Logger),
	}
	em.AddMiddleware()
	em.HealthCheck()

	return em
}

func (e *EndpointManager) AddEndpoints() {

}

func (e *EndpointManager) AddMiddleware() {
	e.router.Use(middleware.NewCorsMiddleware().Cors)
}

func (e *EndpointManager) HealthCheck() {
	e.router.HandleFunc("/health_check", func(w http.ResponseWriter, _ *http.Request) {
		e.responseManager.Message(w, "system is healthy")
	})
}

func (e *EndpointManager) StartServer() error {
	server := &http.Server{
		Addr:    ":" + e.conf.Port,
		Handler: e.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.conf.Logger.Error("failed creating server", zap.Error(err))
		}
	}()
	e.conf.Logger.Info("server started")
	<-e.ctx.Done()
	e.conf.Logger.Info("server stopped")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctxShutDown); err != nil {
		e.conf.Logger.Error("server Shutdown Failed", zap.Error(err))
		return err
	}
	e.conf.Logger.Info("server exited properly")
	return nil
}
