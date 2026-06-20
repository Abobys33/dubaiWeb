package feedback

import (
	"context"
	"dubaiWeb/internal/lib/sl"
	"dubaiWeb/internal/mailer"
	"dubaiWeb/internal/storage/postgres"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

type Request struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type Storage interface {
	CreateFeedbackRequest(ctx context.Context, req postgres.FeedbackRequest) (int64, error)
}

func New(log *slog.Logger, storage Storage, m *mailer.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Phone = strings.TrimSpace(req.Phone)

		if req.Name == "" || req.Phone == "" {
			respondError(w, http.StatusBadRequest, "name and phone are required")
			return
		}

		id, err := storage.CreateFeedbackRequest(r.Context(), postgres.FeedbackRequest{
			Name:  req.Name,
			Phone: req.Phone,
		})
		if err != nil {
			log.Error("failed to save feedback request", sl.Err(err))
			respondError(w, http.StatusInternalServerError, "internal error")
			return
		}

		log.Info("feedback request saved", slog.Int64("id", id))

		// Письмо отправляем асинхронно, чтобы не задерживать ответ пользователю.
		go func() {
			if err := m.SendFeedbackNotification(req.Name, req.Phone); err != nil {
				log.Error("failed to send email notification", sl.Err(err))
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Success: true})
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Success: false, Error: msg})
}
