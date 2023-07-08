package server

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/notifire/internal/integrations/mailgun"
	"github.com/svilenkomitov/notifire/internal/notification"
	"net/http"
	"strconv"
)

type Server struct {
	server *http.Server
}

func New(c *Config) *Server {
	server := setUpServer(c)
	return &Server{
		server: server,
	}
}

func initRoutes(router *chi.Mux) {
	notificationHandler := notification.Handler{
		EmailService: mailgun.New(mailgun.LoadConfig()),
	}
	notificationHandler.Routes(router)
}

func setUpServer(c *Config) *http.Server {
	router := chi.NewRouter()
	initRoutes(router)

	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(int(c.Port)),
		Handler: router,
	}
	return server
}

func (s *Server) Start() error {
	log.Info("starting the HTTP server at addr: ", s.server.Addr)
	if err := s.server.ListenAndServe(); nil != err && err != http.ErrServerClosed {
		log.Errorf("failed to start server: %v", err)
	}
	return nil
}
