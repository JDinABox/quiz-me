package quizme

import (
	_ "embed"
	"net/http"

	"github.com/JDinABox/quiz-me/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newApp(conf *Config) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	if conf.Logging {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)

	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "public, max-age=15778800, immutable")
				h.ServeHTTP(w, r)
			})
		})

		//layouts["resume"] = buildPage(templatesFs, "resume.html")
		r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.FS(web.AssetsFs))))
	})

	return r, nil
}
