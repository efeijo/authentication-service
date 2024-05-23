package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"authservice/internal/authservice"
	"authservice/internal/config"
	"authservice/internal/model"
)

type Server struct {
	authService authservice.AuthService
	mux         *chi.Mux
	port        string
}

func NewServer(service authservice.AuthService, config *config.AuthServiceConfig) *Server {
	s := &Server{authService: service, mux: chi.NewRouter()}

	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	s.addTokensEndpoint()
	s.addUsersEndpoint()

	return s
}

// addUsersEndpoint adds the user endpoints to the server.
func (s *Server) addUsersEndpoint() {
	subRouter := chi.NewRouter()

	s.mux.Mount("/user", subRouter)

	subRouter.Post("/", s.createUserHandler)
	subRouter.Delete("/", s.deleteUser)
	subRouter.Get("/", s.listUsersHandler)
}

// addTokensEndpoint adds the token endpoints to the server.
func (s *Server) addTokensEndpoint() {
	subRouter := chi.NewRouter()

	s.mux.Mount("/token", subRouter)

	subRouter.Post("/", s.createToken)
	subRouter.Delete("/", s.invalidateToken)
	subRouter.Get("/{jwt_token}", s.validateToken)
}

// create is the handler for the POST /token endpoint.
func (s *Server) createToken(w http.ResponseWriter, req *http.Request) {
	log.Println("creating token")
	var userAuth UserAuth
	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}
	token, err := s.authService.GetUserToken(req.Context(), userAuth.Username, userAuth.Password)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}

	WriteJson(w, http.StatusCreated, GetTokenResponse{JwtToken: token})

}

// createUserHandler is the handler for the POST /user endpoint.
func (s *Server) createUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("creating user")
	var userAuth UserAuth
	err := json.NewDecoder(req.Body).Decode(&userAuth)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}

	err = s.authService.CreateUser(req.Context(), userAuth.Username, userAuth.Password)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}

	WriteJson(w, http.StatusCreated, nil)
}

// deleteUser is the handler for the DELETE /user endpoint.
func (s *Server) deleteUser(w http.ResponseWriter, req *http.Request) {
	log.Println("deleting user")
	var deleteRequest DeleteUserRequest
	err := json.NewDecoder(req.Body).Decode(&deleteRequest)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}
	err = s.authService.DeleteUser(req.Context(), deleteRequest.Username)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusInternalServerError})
		return
	}
	WriteJson(w, http.StatusOK, []byte{})
}

// invalidateToken is the handler for the DELETE /token endpoint.
func (s *Server) invalidateToken(w http.ResponseWriter, req *http.Request) {
	log.Println("invalidating token")
	var invalidateRequest InvalidateToken
	err := json.NewDecoder(req.Body).Decode(&invalidateRequest)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusBadRequest})
		return
	}
	err = s.authService.InvalidateToken(req.Context(), invalidateRequest.Username)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusInternalServerError})
		return
	}
}

// listUsersHandler is the handler for the GET /user endpoint.
func (s *Server) listUsersHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("listing users")
	users, err := s.authService.ListUsers(req.Context())
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusInternalServerError})
		return
	}

	if users == nil {
		users = []*model.User{}
	}
	WriteJson(w, http.StatusOK, users)
}

// validateToken is the handler for the GET /token/{jwt_token} endpoint.
func (s *Server) validateToken(w http.ResponseWriter, req *http.Request) {
	log.Println("validating token")
	token := chi.URLParam(req, "jwt_token")

	isValid, err := s.authService.ValidateToken(req.Context(), token)
	fmt.Println(isValid, err)
	if err != nil {
		WriteError(w, ApiError{Err: err, HttpStatusCode: http.StatusInternalServerError})
		return
	}

	WriteJson(w, http.StatusOK, ValidateToken{Valid: isValid})
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.mux)
}
