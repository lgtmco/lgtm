package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Server struct {
	Addr string
	Cert string
	Key  string
}

func (s *Server) Run(handler http.Handler) {
	log.Infof("starting server %s", s.Addr)

	if len(s.Cert) != 0 {
		log.Fatal(
			http.ListenAndServeTLS(s.Addr, s.Cert, s.Key, handler),
		)
	} else {
		log.Fatal(
			http.ListenAndServe(s.Addr, handler),
		)
	}
}
