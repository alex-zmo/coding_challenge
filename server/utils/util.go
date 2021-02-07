package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var l = log.New(os.Stderr, "", log.Ldate)

// Logs error to Stderr.
func LogError(err error) {
	l.Println(err)
}

// Generates Random string using crypto/rand.
func CryptoRandomString(len int) (string, error) {
	random := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, random)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(random), nil
}

// Enables content security policy which only allows self scripts as well as specific nonces for inline scripts.
func EnableContentSecurityPolicy(w http.ResponseWriter) {
	ContentSecurityPolicyString := strings.Join([]string{
		"script-src 'self' 'nonce-EDNnf03nceIOfn39fn3e9h3sdfa' 'nonce-afds3h9e3nf93nf0Iecn30fnNDE'",
		"style-src 'self' 'unsafe-inline' fonts.googleapis.com ",
		"object-src 'none'"}, ";")
	w.Header().Set("Content-Security-Policy", ContentSecurityPolicyString)
}

// Extracts a token from a cookie based on token name.
func ExtractTokenFromCookie(tokenName string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(tokenName)
	if err != nil {
		return "", errors.New(tokenName + " not found")
	}
	return cookie.Value, nil
}

// Saves a token cookie with httpOnly and Secure flag.
func SaveCookie(cookieName, token string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		MaxAge:   0,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// Serves json response.
func ServeJson(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		LogError(err)
		ServeInternalServerError(w)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		LogError(err)
		ServeInternalServerError(w)
		return
	}
}

func ServeBadRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func ServeUnauthorized(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func ServeForbidden(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func ServeMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func ServeInternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
