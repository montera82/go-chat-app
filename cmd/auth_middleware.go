package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"io"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

const (
	googleClientId    = "732340674204-57qqs14p51hjehfnb23132nclgdrqft9.apps.googleusercontent.com"
	googleSecret      = "ivTJaDVemMR06idXkcMHriHC"
	googleRedirectURL = "http://localhost:8083/auth/callback/google"
)

type authMiddlewareHandler struct {
	next http.Handler
}

func (a *authMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")

	if err == http.ErrNoCookie || cookie.Value == "" {
		// not logged in
		w.Header().Add("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(w, "an error has occured", http.StatusInternalServerError)
		return
	}

	a.next.ServeHTTP(w, r)

}

// Custom handler to handle the request
// /auth/login/*  and associated callbacks
func handleAuthenticationProcess(w http.ResponseWriter, r *http.Request) {

	segs := strings.Split(r.URL.Path, "/")

	action := segs[2]
	provider := segs[3]

	gomniauth.SetSecurityKey("very_secret")
	gomniauth.WithProviders(
		facebook.New("", "", "/"),
		google.New(googleClientId, googleSecret, googleRedirectURL),
		github.New("clientId", "clientSecret", "redirect"),
	)

	switch action {

	case "login":
		provider, err := gomniauth.Provider(provider)

		if err != nil {

			http.Error(w,
				fmt.Sprintf("No provider registered for %s %s", provider, err),
				http.StatusInternalServerError)

			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)

		if err != nil {

			http.Error(w,
				fmt.Sprintf("Error when attempting to fetch %s %s", provider, err),
				http.StatusInternalServerError)

			return
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("Error on callback from %s, reason %s ", provider, err),
				http.StatusInternalServerError)
			return
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(
				w, fmt.Sprintf("Error attempting to decode grant reason %s ", err),
				http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(creds)

		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		userId := fmt.Sprintf("%x", m.Sum(nil))
		authCookie := objx.New(
			map[string]interface{}{
				"name":       user.Name(),
				"avatar_url": user.AvatarURL(),
				"email":      user.Email(),
				"userId":     userId,
			}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookie,
			Path:  "/",
		})

		w.Header().Set("location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No handler found for %s", provider)
	}
}

func MustAuth(next http.Handler) http.Handler {
	return &authMiddlewareHandler{next: next}
}
