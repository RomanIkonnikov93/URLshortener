package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/RomanIkonnikov93/URLshortner/internal/handlers/gzipmid"
	"github.com/RomanIkonnikov93/URLshortner/internal/model"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
	"github.com/RomanIkonnikov93/URLshortner/logging"
)

func GzipResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") == "gzip" {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipmid.GzipWriter{ResponseWriter: w, Writer: gz}, r)
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

func GzipRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		logger := logging.GetLogger()

		if r.Header.Get("Content-Encoding") == "gzip" {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err := gzipmid.DecompressGZIP(b)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(strings.NewReader(string(data)))
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func UserValidation(rep repository.Pool, logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("UserTokenID")
			if err != nil {
				ctx, err := SetUserCtx(w, r, model.Key, rep)
				if err != nil {
					logger.Error(err)
				}
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			c := cookie.Value
			userID, err := CheckUSerCookies(c, model.Key, rep)
			if err != nil || userID == "" {
				ctx, err := SetUserCtx(w, r, model.Key, rep)
				if err != nil {
					logger.Error(err)
				}
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserCtx("userID"), userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
