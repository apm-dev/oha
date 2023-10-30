package httpserver

import "net/http"

func (s *Server) RegisterRoutes(pathPrefix string) {
	insecureRouter := s.router.PathPrefix(pathPrefix).Subrouter()

	insecureRouter.Methods(http.MethodPost).
		Path("/users").HandlerFunc(s.userCtrl.CreateUser)

	insecureRouter.Methods(http.MethodGet).
		Path("/users/{id}").HandlerFunc(s.userCtrl.GetUserByID)
}
