package server

import (
	"fmt"

	"otus-hl-network/internal/auth"
	authdl "otus-hl-network/internal/auth/delivery"
	userdl "otus-hl-network/internal/user/delivery"

	"github.com/labstack/echo"
)

type Server struct {
	user           *userdl.UserHander
	auth           *authdl.AuthHandler
	e              *echo.Echo
	authMiddleware *AuthMiddleware
}

func New(authManager auth.AuthManager, user *userdl.UserHander, auth *authdl.AuthHandler) *Server {
	e := echo.New()
	e.HideBanner = true

	return &Server{
		authMiddleware: &AuthMiddleware{authManager: authManager},
		user:           user,
		auth:           auth,
		e:              e,
	}
}

func (s *Server) Run(port int) error {

	apiGroup := s.e.Group("/api/v1")

	// API endpoints
	apiGroup.POST("/auth", s.auth.Authorize)
	apiGroup.POST("/register", s.user.Register)
	apiGroup.GET("/profile/:id", s.authMiddleware.Do(s.user.Profile))
	apiGroup.GET("/profile", s.authMiddleware.Do(s.user.Profile))
	apiGroup.POST("/profile", s.authMiddleware.Do(s.user.UpdateProfile))
	apiGroup.GET("/users", s.authMiddleware.Do(s.user.Users))
	apiGroup.GET("/friends", s.authMiddleware.Do(s.user.Friends))
	apiGroup.POST("/friends/:id", s.authMiddleware.Do(s.user.AddFriend))
	apiGroup.DELETE("/friends/:id", s.authMiddleware.Do(s.user.DeleteFriend))
	apiGroup.GET("/search", s.authMiddleware.Do(s.user.Search))

	s.e.Static("/assets", "/frontend/assets")
	s.e.File("/", "/frontend/index.html")
	s.e.File("/favicon.ico", "/frontend/favicon.ico")

	// s.e.Static("/assets", "/app/frontend/dist/assets")
	// s.e.File("/", "/app/frontend/dist/index.html")
	// s.e.File("/favicon.ico", "/app/frontend/dist/favicon.ico")

	// render endpoints
	// s.e.GET("/auth", nil)
	// s.e.GET("/register", nil)
	// s.e.GET("/profile", nil)
	// s.e.GET("/users", nil)

	return s.e.Start(fmt.Sprintf(":%d", port))
}
