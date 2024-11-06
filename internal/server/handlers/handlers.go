package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/services"
	"go.uber.org/zap"
)

const limitInSeconds = 5

func GetNTimestampsHandler(logger *zap.Logger, service services.SnapshotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		n, err := strconv.Atoi(chi.URLParam(r, "timestampsCount"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		timestamps, err := service.GetNTimestamps(ctx, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response strings.Builder
		for id, timestamp := range timestamps {
			if _, err := response.Write([]byte(fmt.Sprintf("%d: %v\n", id, timestamp))); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(response.String())); err != nil {
			logger.Error(err.Error())
		}
	}
}

func GetSnapshotHandler(logger *zap.Logger, service services.SnapshotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		id, err := strconv.Atoi(chi.URLParam(r, "snapshotID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		snapshot, err := service.GetSnapshot(ctx, model.ID(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.MarshalIndent(snapshot, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response); err != nil {
			logger.Error(err.Error())
		}
	}
}

func DeleteSnapshotHandler(logger *zap.Logger, service services.SnapshotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		id, err := strconv.Atoi(chi.URLParam(r, "snapshotID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := service.DeleteSnapshot(ctx, model.ID(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Snapshot is deleted")); err != nil {
			logger.Error(err.Error())
		}
	}
}
