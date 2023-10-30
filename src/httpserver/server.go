package httpserver

import (
	"fmt"
	"net/http"

	userhttp "github.com/apm-dev/oha/src/user/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
	router *mux.Router

	userCtrl *userhttp.UserController
}

func NewServer(
	uc *userhttp.UserController,
) *Server {
	return &Server{
		server:   nil,
		router:   nil,
		userCtrl: uc,
	}
}

func (s *Server) Start(port uint, pathPrefix string) {
	s.router = mux.NewRouter()
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}

	s.RegisterRoutes(pathPrefix)

	go func() {
		log.Printf("http server is listening on port %d", port)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start http server: %s", err)
		}
	}()
}

func (s *Server) Stop() {
	err := s.server.Close()
	if err != nil {
		log.Infof("failed to stop http server: %s", err)
	}
}
