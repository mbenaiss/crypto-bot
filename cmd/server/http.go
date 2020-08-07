package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mbenaiss/crypto-bot/config"
	"github.com/mbenaiss/crypto-bot/internal/service"
)

type Server struct {
	httpServer      *http.Server
	technicalServer *http.Server
	cfg             *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) StartHTTP(svc *service.Service) error {
	r := gin.Default()
	r.POST("/upload/:provider", uploadCsv(svc))
	r.POST("/provider", addProvider(svc))

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HttpPort),
		Handler: r,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) StartHealthz() error {
	r := gin.Default()
	ok := func(c *gin.Context) {
		c.Status(200)
	}
	r.GET("/readiness", ok)
	r.GET("/liveness", ok)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	s.technicalServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HealthzPort),
		Handler: r,
	}
	return s.technicalServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	errHTTP := s.httpServer.Shutdown(ctx)
	if errHTTP != nil {
		return errHTTP
	}
	errTech := s.technicalServer.Shutdown(ctx)
	if errTech != nil {
		return errTech
	}
	return nil
}
