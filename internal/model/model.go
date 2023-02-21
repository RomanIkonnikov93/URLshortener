package model

import (
	"errors"
	"time"
)

// Key for Encrypt and Decrypt func
var Key = []byte{42, 166, 161, 254, 64, 82, 58, 84, 65, 191, 62, 76, 181, 24, 199, 74}

var (
	ErrConflict = errors.New("conflict on insert")
	ErrDelFlag  = errors.New("url is deleted")
)

const TimeOut = time.Second * 5

// URLRequest structure for func PostJSONHandler
type URLRequest struct {
	URL string `json:"url"`
}

// URLResponse structure for func PostJSONHandler
type URLResponse struct {
	Result string `json:"result"`
}

// URLsJSONResponse structure for func GetAllUserURLs
type URLsJSONResponse struct {
	Short string `json:"short_url"`
	Long  string `json:"original_url"`
}

// BatchRequest structure for func PostBatchHandler
type BatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse structure for func PostBatchHandler
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserRequest structure for func DeleteUserUrls
type UserRequest struct {
	UserID   string
	UserUrls []string
}
