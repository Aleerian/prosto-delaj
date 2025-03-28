package server

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	ReadTimeout = 60
	WriteTimeout
	MinHeaderBytes = 1
	MaxHeaderBytes = 20
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: MinHeaderBytes << MaxHeaderBytes,
		ReadTimeout:    ReadTimeout * time.Second,
		WriteTimeout:   WriteTimeout * time.Second,
	}
	logrus.Info("server started successfully")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context, postgres *sql.DB) {
	logrus.Info("server shutdown process started")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Info("http listener shutdown successfully")
	}

	if err := postgres.Close(); err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Info("business database connection closed successfully")
	}

	logrus.Info("server shutdown process completed successfully")
}
