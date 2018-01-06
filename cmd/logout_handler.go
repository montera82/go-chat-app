package main

import "net/http"

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header().Set("location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
