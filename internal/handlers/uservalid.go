package handlers

import (
	"context"
	"crypto/aes"
	"encoding/hex"
	"math/rand"
	"net/http"
	"time"

	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
)

type UserCtx string

func GenerateUserID() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Encrypt(src, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, aes.BlockSize)
	aesblock.Encrypt(dst, src)
	return dst, nil
}

func Decrypt(src, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, aes.BlockSize)
	aesblock.Decrypt(dst, src)
	return dst, nil
}

func CheckUSerCookies(c string, key []byte, rep repository.Pool) (string, error) {
	src, err := hex.DecodeString(c)
	if err != nil {
		return "", err
	}
	res, err := Decrypt(src, key)
	if err != nil {
		return "", err
	}
	k := string(res)
	exist, err := rep.Users.CheckUserID(k)
	if err != nil {
		return "", err
	}
	if exist {
		return k, nil
	} else {
		return "", err
	}
}

func CreateUserCookiesAndUserID(key []byte, rep repository.Pool) (cookie string, ID string, err error) {
	ID = GenerateUserID()
	err = rep.Users.AddUserID(ID)
	if err != nil {
		return "", "", err
	}
	c, err := Encrypt([]byte(ID), key)
	if err != nil {
		return "", "", err
	}
	cookie = hex.EncodeToString(c)
	return
}

func SetCookie(w http.ResponseWriter, c string) {
	cookie := &http.Cookie{
		Name:    "UserTokenID",
		Value:   c,
		Path:    "/",
		Domain:  "",
		Expires: time.Now().Add(time.Hour * 24),
	}
	w.Header().Set("Set-Cookie", cookie.String())
}

func SetUserCtx(w http.ResponseWriter, r *http.Request, key []byte, rep repository.Pool) (context.Context, error) {
	c, ID, err := CreateUserCookiesAndUserID(key, rep)
	if err != nil {
		return nil, err
	}
	SetCookie(w, c)
	return context.WithValue(r.Context(), UserCtx("userID"), ID), nil
}
