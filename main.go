package main

import (
	"net/http"
	"strings"

	"github.com/sybrenstuvel/site-soundabout/views"
)

func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	index := views.NewView("master", "views/index.gohtml")
	// contact := views.NewView("bootstrap", "views/contact.gohtml")
	static := noDirListing(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index.Render(w, r, nil)
	})
	// http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
	// 	contact.Render(w, r, nil)
	// })
	http.Handle("/static/", static)
	http.ListenAndServe(":3000", nil)
}
