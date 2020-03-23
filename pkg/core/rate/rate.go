package rate

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

var ErrInvalidRate = errors.New("invalid rate")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type ResponseDTO struct {
	Id          int64    `json:"id"`
	MovieId     int64    `json:"movie_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Creators    []string `json:"creators"`
	Studio      string   `json:"studio"`
	ExtLink     string   `json:"ext_link"`
	UserId      int64    `json:"user_id"`
	UserName    string   `json:"user_name"`
	Rate        int64    `json:"rate"`
}

type MovieRateDTO struct {
	MovieId     int64    `json:"movie_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Creators    []string `json:"creators"`
	Studio      string   `json:"studio"`
	ExtLink     string   `json:"ext_link"`
	Rate        float64  `json:"rate"`
}

func (s *Service) AddRate(ctx context.Context, request ResponseDTO) (err error) {
	if request.Rate < 0 {
		return ErrInvalidRate
	}
	if request.Rate > 10 {
		return ErrInvalidRate
	}
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("unable to acquire ctx: %w", err)
	}

	_, err = conn.Exec(ctx, addRateDML,
		request.MovieId,
		request.Title,
		request.Description,
		request.Image,
		request.Year,
		request.Country,
		request.Actors,
		request.Genres,
		request.Creators,
		request.Studio,
		request.ExtLink,
		request.UserId,
		request.UserName,
		request.Rate,
	)
	if err != nil {
		return fmt.Errorf("unable to add rate: %w", err)
	}
	return nil
}

func (s *Service) GetAllRates() (rates []MovieRateDTO, err error) {
	rows, err := s.pool.Query(context.Background(), getAllRatesDML)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ra := MovieRateDTO{}
		err := rows.Scan(
			&ra.MovieId,
			&ra.Title,
			&ra.Description,
			&ra.Image,
			&ra.Year,
			&ra.Country,
			&ra.Actors,
			&ra.Genres,
			&ra.Creators,
			&ra.Studio,
			&ra.ExtLink,
			&ra.Rate,
		)
		if err != nil {
			return nil, err
		}

		rates = append(rates, ra)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rates, err
}

func (s *Service) GetRate(request *http.Request) (rates MovieRateDTO, err error) {
	row := s.pool.QueryRow(context.Background(), getRateDML, request.Context().Value("id"))

	err = row.Scan(
		&rates.MovieId,
		&rates.Title,
		&rates.Description,
		&rates.Image,
		&rates.Year,
		&rates.Country,
		&rates.Actors,
		&rates.Genres,
		&rates.Creators,
		&rates.Studio,
		&rates.ExtLink,
		&rates.Rate,
	)

	if err != nil {
		return rates, err
	}

	return rates, nil
}

func (s *Service) DeleteRate(ctx context.Context, id int64) (err error) {
	_, err = s.pool.Exec(ctx, deleteRateDML, ctx.Value("id"), id)
	if err != nil {
		return err
	}
	return nil
}
