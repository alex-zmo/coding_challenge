package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var l = log.New(os.Stderr, "", 1)

// Logs error to Stderr.
func LogError(a ...interface{}) {
	l.Println(a...)
}

// Generates Random string using crypto/rand.
func CryptoRandomString(len int) (string, error) {
	random := make([]byte, len)
	_, err := rand.Reader.Read(random)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(random), nil
}

// Enables content security policy which only allows self scripts as well as specific nonces for inline scripts.
func EnableContentSecurityPolicy(w *http.ResponseWriter) {
	ContentSecurityPolicyString := strings.Join([]string{
		"script-src 'self' 'nonce-EDNnf03nceIOfn39fn3e9h3sdfa' 'nonce-afds3h9e3nf93nf0Iecn30fnNDE'",
		"style-src 'self' 'unsafe-inline' fonts.googleapis.com ",
		"object-src 'none'"}, ";")
	(*w).Header().Set("Content-Security-Policy", ContentSecurityPolicyString)
}

// Gets header from request based on key.
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
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
	w.Header().Add("Vary", "Cookie")
}

// Gets body from request.
func GetBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			return nil, err
	}
	return body, nil
}

// Serves Json Response.
func ServeJson(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
			LogError(err)
			return
	}
	_, err = w.Write(data)
	if err != nil {
			LogError(err)
			return
	}
}

func ServeCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func ServeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func ServeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func ServeForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func ServeNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func ServeMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func ServeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
