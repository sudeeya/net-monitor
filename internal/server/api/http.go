package api

import (
	"html/template"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/server/handlers"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

// Endpoints.
const (
	defaultEndpoint       = "/"
	getTimestampsEndpoint = "/timestamps"
	getSnapshotEndpoint   = "/snapshots"
)

// snapshotsHTTPServer defines object to interact with the server using HTTP.
type snapshotsHTTPServer struct {
	*chi.Mux
	logger  *zap.Logger
	service services.SnapshotsService
}

// Paths to HTML files.
var (
	commonPath     = filepath.Join("assets", "html", "common.html")
	indexPath      = filepath.Join("assets", "html", "index.html")
	timestampsPath = filepath.Join("assets", "html", "timestamps.html")
	snapshotsPath  = filepath.Join("assets", "html", "snapshots.html")
)

// NewSnapshotsHTTPServer returns snapshotsHTTPServer object.
func NewSnapshotsHTTPServer(logger *zap.Logger, service services.SnapshotsService) (*snapshotsHTTPServer, error) {
	mux := chi.NewRouter()

	tmpls, err := parseHTMLFiles()
	if err != nil {
		return nil, err
	}

	registerEndpoints(mux, logger, service, tmpls)

	return &snapshotsHTTPServer{
		Mux:     mux,
		logger:  logger,
		service: service,
	}, nil
}

// parseHTMLFiles parses HTML files for enpoints.
func parseHTMLFiles() (map[string]*template.Template, error) {
	indexTmpl, err := template.ParseFiles(indexPath, commonPath)
	if err != nil {
		return nil, err
	}

	timestampsTmpl, err := template.ParseFiles(timestampsPath, commonPath)
	if err != nil {
		return nil, err
	}

	snapshotsTmpl, err := template.ParseFiles(snapshotsPath, commonPath)
	if err != nil {
		return nil, err
	}

	return map[string]*template.Template{
		defaultEndpoint:       indexTmpl,
		getTimestampsEndpoint: timestampsTmpl,
		getSnapshotEndpoint:   snapshotsTmpl,
	}, nil
}

// registerEndpoints registers enpoints for HTTP requests.
func registerEndpoints(
	mux *chi.Mux,
	logger *zap.Logger,
	service services.SnapshotsService,
	tmpls map[string]*template.Template,
) {
	mux.Get(defaultEndpoint, handlers.DefaultHandler(logger, tmpls[defaultEndpoint]))
	mux.Get(getTimestampsEndpoint, handlers.GetTimestampsHandler(logger, service, tmpls[getTimestampsEndpoint]))
	mux.Get(getSnapshotEndpoint, handlers.GetSnapshotHandler(logger, service, tmpls[getSnapshotEndpoint]))
}
