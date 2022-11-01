package server

import (
	"context"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/redirectBL"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/starter"
	"log"
	"net/http"
	"time"
)

var _ starter.APIServer = &Server{}

type Server struct {
	srv        http.Server
	redirectBL *redirectBL.Redirect
}

func NewServer(addr string, h http.Handler) *Server {
	server := &Server{}

	server.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return server
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Println("server shutdown error: ", err)
	}
	cancel()
}

func (s *Server) Start(redirectBL *redirectBL.Redirect) {
	s.redirectBL = redirectBL
	go func(*Server) {
		err := s.srv.ListenAndServe()
		if err != nil {
			log.Println("server error: ", err)
		}
	}(s)
}
