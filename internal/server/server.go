package server

import (
	"fmt"

	"otus-hl-network/internal/auth"
	"otus-hl-network/internal/user/delivery"

	"github.com/labstack/echo"
)

type Server struct {
	user           *delivery.UserHander
	e              *echo.Echo
	authMiddleware *AuthMiddleware
}

func New(jwtManager *auth.JWTManager, user *delivery.UserHander) *Server {
	e := echo.New()
	e.HideBanner = true

	return &Server{
		authMiddleware: &AuthMiddleware{jwtManager: jwtManager},
		user:           user,
		e:              e,
	}
}

func (s *Server) Run(port int) error {

	apiGroup := s.e.Group("/api/v1")

	// s.e.GET("/auth", s.user.Login)
	// s.e.POST("/auth", s.user.Authorize)
	//
	// s.e.GET("/register", s.user.Authorize)

	// API endpoints
	// apiGroup.POST("/register", nil)
	//
	// apiGroup.GET("/profile/{id}", nil)
	// apiGroup.GET("/profile", nil)

	apiGroup.GET("/users", s.authMiddleware.Do(s.user.Users))

	// apiGroup.GET("/friends", nil)
	// apiGroup.POST("/friends", nil)
	// apiGroup.DELETE("/friends", nil)

	// render endpoints
	// s.e.GET("/auth", nil)
	// s.e.GET("/register", nil)
	// s.e.GET("/profile", nil)
	// s.e.GET("/users", nil)

	return s.e.Start(fmt.Sprintf(":%d", port))
}
