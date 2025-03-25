// Package handlers provides a collection of HTTP handlers.
package handlers

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/server/services"
)

const limitInSeconds = 5

// DefaultHandler returns an http.HandlerFunc that writes default page to the response.
// If an error occurs, it logs the error and returns an appropriate HTTP status code.
func DefaultHandler(logger *zap.Logger, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetTimestampsHandler returns an http.HandlerFunc that requests a list of
// snapshot ids and timestamps from the service and writes them to the response.
// If an error occurs, it logs the error and returns an appropriate HTTP status code.
func GetTimestampsHandler(logger *zap.Logger, service services.SnapshotsService, tmpl *template.Template) http.HandlerFunc {
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

		if err = tmpl.Execute(w, timestamps); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetSnapshotHandler returns an http.HandlerFunc that requests a snapshot
// from the service and writes it to the response in json format.
// If an error occurs, it logs the error and returns an appropriate HTTP status code.
func GetSnapshotHandler(logger *zap.Logger, service services.SnapshotsService, tmpl *template.Template) http.HandlerFunc {
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

		if err = tmpl.Execute(w, snapshot); err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
