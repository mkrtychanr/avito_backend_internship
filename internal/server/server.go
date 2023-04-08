package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type Server struct {
	db     *sql.DB
	config *config
	router *chi.Mux
	mu     sync.Mutex
}

func NewServer(configPath string) (*Server, error) {
	config, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	// connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=http://postgres sslmode=disable", config.User, config.Password, config.DBname)
	connStr := fmt.Sprintf("postgresql://%s:%s@postgres/%s?sslmode=disable", config.User, config.Password, config.DBname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Server{
		db:     db,
		config: config,
		router: chi.NewMux(),
		mu:     sync.Mutex{},
	}, nil
}

func (s *Server) setRouter() {
	s.router.Use(middleware.Logger)
	s.router.Get("/balance/{client_id}", s.GetClientBalance)
	s.router.Get("/reports/{file}", s.GetReport)
	s.router.Post("/add_money", s.AddMoney)
	s.router.Post("/reserve_money", s.ReserveMoney)
	s.router.Post("/unreserve_money", s.UnreserveMoney)
	s.router.Post("/allow_transaction", s.AllowTransaction)
	s.router.Post("/generate_report", s.GenerateReport)
}

func (s *Server) Up() {
	address := s.GetAddres()
	s.setRouter()
	fmt.Printf("Server is up on %s\n", address)
	http.ListenAndServe(address, s.router)
}

func (s *Server) GetAddres() string {
	return fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
}
