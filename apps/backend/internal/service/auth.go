package service

import (
	"github.com/sriniously/go-resumify/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server: s,
	}
}
