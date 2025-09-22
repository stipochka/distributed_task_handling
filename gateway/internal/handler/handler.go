package handler

import (
	"context"
	"encoding/json"
	"gateway/internal/models"
	"gateway/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
)

var validate = validator.New()

func writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{
		"status":  "error",
		"message": msg,
	})
}

type Handler struct {
	service service.TaskService
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func NewHandler(service service.TaskService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoute() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(h.LogMiddleware)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Logger)

	router.Post("/api/v1/task", h.CreateTask)

	return router
}

func (h *Handler) LogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		t1 := time.Now()

		next.ServeHTTP(sw, r)

		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": sw.status,
			"duration":    time.Since(t1),
		}).Info("request completed")
	}

	return http.HandlerFunc(fn)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	logrus.Info("received create task request")

	task := models.Task{}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to decode request body")

		writeError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validate.Struct(task); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("invalid request body")

		writeError(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.service.SendTask(context.Background(), &task)
	if err != nil { // can add error handling
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to send task")

		writeError(w, r, http.StatusInternalServerError, "internal error")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status":  "created",
		"task_id": id.String(),
	})

}
