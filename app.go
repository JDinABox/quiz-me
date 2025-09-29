package quizme

import (
	_ "embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/JDinABox/quiz-me/web"
	"github.com/JDinABox/quiz-me/web/templates"
	"github.com/JDinABox/quiz-me/web/templates/home"
	"github.com/a-h/templ"
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

	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			linkPreload, err := web.GetLinkPreload()
			if err != nil {
				slog.Error("failed to get link preload", "error", err)
				os.Exit(1)
			}
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if !conf.Dev {
					w.Header().Set("Link", linkPreload)
					w.Header().Set("Cache-Control", "private, max-age=300")
				}
				h.ServeHTTP(w, r)
			})
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			err := execPage(w, r, home.Head(), home.Body())
			if err != nil {
				slog.Error("failed to render home page", "error", err)
			}
		})
	})
	return r, nil
}

func execPage(w http.ResponseWriter, r *http.Request, head, body templ.Component) error {
	/*if r.Header.Get("x-alpine-request") == "true" {
		return template.Main(body).Render(r.Context(), w)
	}*/
	return templates.Layout(head, body, r.RequestURI).Render(r.Context(), w)
}
