// Package handlers provides a collection of HTTP handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/server/services"
)

const limitInSeconds = 5

func DefaultHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("assets", "html", "index.html")
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetNTimestampsHandler returns an http.HandlerFunc that requests a list of
// snapshot ids and timestamps from the service and writes them to the response.
// If an error occurs, it logs the error and returns an appropriate HTTP status code.
func GetTimestampsHandler(logger *zap.Logger, service services.SnapshotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		n, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		timestamps, err := service.GetNTimestamps(ctx, n)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response strings.Builder
		for _, timestamp := range timestamps {
			if _, err := response.Write([]byte(fmt.Sprintf("%d: %v\n", timestamp.ID, timestamp.Timestamp))); err != nil {
				logger.Error(err.Error())
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

// GetSnapshotHandler returns an http.HandlerFunc that requests a snapshot
// from the service and writes it to the response in json format.
// If an error occurs, it logs the error and returns an appropriate HTTP status code.
func GetSnapshotHandler(logger *zap.Logger, service services.SnapshotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
		defer cancel()

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		snapshot, err := service.GetSnapshot(ctx, id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.MarshalIndent(snapshot, "", "\t")
		if err != nil {
			logger.Error(err.Error())
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
