package app

import (
	"errors"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	jwt2 "github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/rest"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"rate-service/pkg/core/rate"
	"rate-service/pkg/core/token"
)

type Server struct {
	router   *mux.ExactMux
	pool     *pgxpool.Pool
	secret   jwt.Secret
	movieSvc *rate.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, userSvc *rate.Service) *Server {
	return &Server{router: router, pool: pool, secret: secret, movieSvc: userSvc}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleGetRates() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		movies, err := s.movieSvc.GetAllRates()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
		}

		if len(movies) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		err = rest.WriteJSONBody(writer, &movies)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleNewRate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dto rate.ResponseDTO
		err := rest.ReadJSONBody(request, &dto)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		movi := rate.ResponseDTO{
			Id:          dto.Id,
			MovieId:     dto.MovieId,
			Title:       dto.Title,
			Description: dto.Description,
			Image:       dto.Image,
			Year:        dto.Year,
			Country:     dto.Country,
			Actors:      dto.Actors,
			Genres:      dto.Genres,
			Creators:    dto.Creators,
			Studio:      dto.Studio,
			ExtLink:     dto.ExtLink,
			UserId:      dto.UserId,
			UserName:    dto.UserName,
			Rate:        dto.Rate,
		}

		err = s.movieSvc.AddRate(request.Context(), movi)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}
	}
}

func (s *Server) handleGetRate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ra, err := s.movieSvc.GetRate(request)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
		}

		err = rest.WriteJSONBody(writer, &ra)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleDeleteRate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)

		err := s.movieSvc.DeleteRate(request.Context(), payload.Id)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}