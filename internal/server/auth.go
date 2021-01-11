package server

import (
	"encoding/json"
	"net/http"

	"github.com/imarrche/jwt-auth-example/internal/model"
)

// signUp signs up a user.
func (s *Server) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.service.Auth().SignUp(u)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

// signIn returns access and refresh JWTs for user.
func (s *Server) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		accessJWT, refreshJWT, err := s.service.Auth().SignIn(req.Email, req.Password)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		res := signInResponse{AccessToken: accessJWT, RefreshToken: refreshJWT}
		s.respond(w, r, http.StatusOK, res)
	}
}

type refreshRequest struct {
	RefreshToken string `json:"refresh"`
}

type refreshResponse struct {
	AccessToken string `json:"access"`
}

// refresh returns new access JWT for user if valid refresh JWT is provided.
func (s *Server) refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		accessJWT, err := s.service.Auth().RefreshAccessJWT(req.RefreshToken)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		res := refreshResponse{AccessToken: accessJWT}
		s.respond(w, r, http.StatusOK, res)
	}
}

// public is a public route to test JWT authorization.
func (s *Server) public() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, nil)
	}
}

// private is a private route to test JWT authorization.
func (s *Server) private() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, nil)
	}
}
