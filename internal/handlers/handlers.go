package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/RomanIkonnikov93/URLshortner/internal/model"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
	"github.com/RomanIkonnikov93/URLshortner/logging"
)

// GetHandler get long URL by short URL
func GetHandler(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b := strings.Trim(r.URL.Path, "/")
		resp, err := rep.Storage.Get(r.Context(), b)
		if err != nil {
			if err == model.ErrDelFlag {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusGone)
				return
			} else {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		w.Header().Set("Location", resp)
		http.Redirect(w, r, resp, http.StatusTemporaryRedirect)
	}
}

// PostHandler get long URL and return short URL
func PostHandler(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// validate URL function
		_, err = url.ParseRequestURI(string(b))
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// generate short URL
		sURL := Short()
		// make response
		short := "http://" + r.Host + "/" + sURL

		// add userID and URLs in repository
		userID, ok := r.Context().Value(UserCtx("userID")).(string)
		if !ok || userID == "" {
			logger.Printf("userID not exist %v", http.StatusInternalServerError)
			http.Error(w, "userID not exist", http.StatusInternalServerError)
			return
		}
		if err := rep.Storage.Add(r.Context(), sURL, string(b), userID); err != nil {
			if errors.Is(err, model.ErrConflict) {
				s, err := rep.Storage.GetShort(r.Context(), string(b))
				if err != nil {
					logger.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				logger.Printf("%v", http.StatusConflict)
				w.WriteHeader(http.StatusConflict)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, _ = w.Write([]byte("http://" + r.Host + "/" + s))
				return
			}
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(short))
	}
}

// PostJSONHandler get long URL in JSON format, return short URL in JSON format
func PostJSONHandler(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// unmarshal request
		data := model.URLRequest{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set that problem HTML characters should not be escaped
		buf := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		err = encoder.Encode(data.URL)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// validate URL function
		_, err = url.ParseRequestURI(data.URL)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set a value to response
		res := model.URLResponse{}
		res.Result = buf.String()

		// generate short URL
		sURL := Short()
		// make response
		short := "http://" + r.Host + "/" + sURL

		// add userID and URLs in repository
		userID, ok := r.Context().Value(UserCtx("userID")).(string)
		if !ok || userID == "" {
			logger.Printf("userID not exist %v", http.StatusInternalServerError)
			http.Error(w, "userID not exist", http.StatusInternalServerError)
			return
		}
		if err := rep.Storage.Add(r.Context(), sURL, data.URL, userID); err != nil {
			if errors.Is(err, model.ErrConflict) {
				s, err := rep.Storage.GetShort(r.Context(), data.URL)
				if err != nil {
					logger.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				res.Result = "http://" + r.Host + "/" + s
				j, err := json.Marshal(&res)
				if err != nil {
					logger.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				logger.Printf("%v", http.StatusConflict)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write(j)
				return
			}
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// marshal response
		res.Result = short
		j, err := json.Marshal(&res)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(j)
	}
}

// GetAllUserURLs get userID, return all User short and long URLs in JSON format
func GetAllUserURLs(repository repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(UserCtx("userID")).(string)
		if !ok || userID == "" {
			logger.Printf("userID not exist %v", http.StatusInternalServerError)
			http.Error(w, "userID not exist", http.StatusInternalServerError)
			return
		}
		data, err := repository.Storage.GetByUserID(r.Context(), userID)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// make response structure
		if len(data) == 0 {
			logger.Printf("%v", http.StatusNoContent)
			w.WriteHeader(http.StatusNoContent)
		} else {
			var arr []*model.URLsJSONResponse
			for s, l := range data {
				res := new(model.URLsJSONResponse)
				res.Short = "http://" + r.Host + "/" + s
				res.Long = l
				arr = append(arr, res)
			}

			// marshal response
			j, err := json.Marshal(&arr)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(j)
		}
	}
}

// PostBatchHandler get batch URLs in JSON format, return many short URLs in JSON format
func PostBatchHandler(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// unmarshal request
		data := model.BatchRequest{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// range data response
		var arr []*model.BatchResponse
		for _, val := range data {

			// set that problem HTML characters should not be escaped
			buf := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buf)
			encoder.SetEscapeHTML(false)
			err = encoder.Encode(val.OriginalURL)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// validate URL function
			_, err = url.ParseRequestURI(val.OriginalURL)
			if err != nil {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// set a value to response
			val.OriginalURL = buf.String()
			long := strings.Trim(val.OriginalURL, "\"\n")
			// generate short URL
			sURL := Short()

			// add userID and URLs in repository
			userID, ok := r.Context().Value(UserCtx("userID")).(string)
			if !ok || userID == "" {
				logger.Printf("userID not exist %v", http.StatusInternalServerError)
				http.Error(w, "userID not exist", http.StatusInternalServerError)
				return
			}
			if err := rep.Storage.Add(r.Context(), sURL, long, userID); err != nil {
				if errors.Is(err, model.ErrConflict) {
					s, err := rep.Storage.GetShort(r.Context(), long)
					if err != nil {
						logger.Error(err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					sURL = s
				}
			}

			// make response
			short := "http://" + r.Host + "/" + sURL
			res := new(model.BatchResponse)
			res.ShortURL = short
			res.CorrelationID = val.CorrelationID
			arr = append(arr, res)
		}

		// marshal response
		j, err := json.Marshal(&arr)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(j)
	}
}

// DeleteUserURLs get batch short URLs ID in JSON format, changes the status in the database to deleted
func DeleteUserURLs(rep repository.Pool, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// add userID and URLs in repository
		userID, ok := r.Context().Value(UserCtx("userID")).(string)
		if !ok || userID == "" {
			logger.Printf("userID not exist %v", http.StatusInternalServerError)
			http.Error(w, "userID not exist", http.StatusInternalServerError)
			return
		}

		//unmarshal request
		data := model.UserRequest{
			UserID: userID,
		}

		err = json.Unmarshal(b, &data.UserUrls)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			err = rep.Storage.BatchDelete(data)
			logger.Error(err)
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}
