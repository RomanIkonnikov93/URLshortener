package handlers

import (
	"net/http"

	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
	"github.com/RomanIkonnikov93/URLshortner/logging"
)

func PingDataBase(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rep.Ping.PingDB()
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
