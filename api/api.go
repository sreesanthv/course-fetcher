package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sreesanthv/course-fetcher/database"
	"github.com/sreesanthv/course-fetcher/logging"
)

func New() (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := database.DBConn()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	store := database.NewStore(db, logger)

	handler := NewHandler(logger, store)
	course := NewCourseHandler(handler)

	r := chi.NewRouter()
	r.Use(middleware.Timeout(30 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/course", course.Router())
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("status"))
	})

	return r, nil
}
