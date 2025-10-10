package handler

import (
	"api/internal/models"
	"api/internal/repository"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{
		"status":  "error",
		"message": msg,
	})
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	router.Get("/api/tasks/{task_id}", h.GetTaskByID)
	router.Get("/api/tasks/{task_type}", h.GetTasksByType)

	return router
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "task_id")

	logrus.WithField("task_id", taskIDStr).Info(
		"received request",
	)

	taskID, err := uuid.FromBytes([]byte(taskIDStr))
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("not valid taskID: %s", err.Error()))

		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"taskID": taskID,
		}).Error("failed to parse task id")
		return
	}

	task, err := h.repo.GetTask(r.Context(), taskID)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "failed to get task")

		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"taskID": taskID,
		}).Error("failed to get task")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]models.Task{
		"task": task,
	})
}

func (h *Handler) GetTasksByType(w http.ResponseWriter, r *http.Request) {
	taskType := chi.URLParam(r, "task_type")

	logrus.WithField("task_type", taskType).Info("received request")

	if taskType == "" {
		writeError(w, r, http.StatusBadRequest, "empty task type not allowed")

		logrus.Info("given empty task_type")
		return
	}

	tasks, err := h.repo.GetAllTaskByType(r.Context(), taskType)
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "failed to get tasks")

		logrus.WithFields(logrus.Fields{
			"error":     err.Error(),
			"task_type": taskType,
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string][]models.Task{
		"tasks": tasks,
	})
}
