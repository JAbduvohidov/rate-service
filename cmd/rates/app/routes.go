package app

import (
	"github.com/JAbduvohidov/mux/middleware/authenticated"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"rate-service/pkg/core/token"
	"reflect"
)

func (s *Server) InitRoutes() {
	s.router.GET("/api/rates",
		s.handleGetRates(),
		logger.Logger("RATES"),
	)

	s.router.GET("/api/rates/{id}",
		s.handleGetRate(),
		logger.Logger( "RATE"),
	)

	s.router.DELETE("/api/rates/{id}",
		s.handleDeleteRate(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("RATE"),
	)

	s.router.POST("/api/rates/",
		s.handleNewRate(),
		logger.Logger("RATE"),
	)
}
