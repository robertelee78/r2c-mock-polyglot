package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"example.com/snapstore/gateway/internal/native"
	"example.com/snapstore/gateway/internal/store"

	"golang.org/x/text/language"
)

// preferredLocale picks the best-matching supported locale for a request from
// its Accept-Language header, defaulting to en-US.
func preferredLocale(header string) string {
	tags, _, err := language.ParseAcceptLanguage(header)
	if err != nil || len(tags) == 0 {
		return "en-US"
	}
	return tags[0].String()
}

// Server holds the request-handling dependencies for the gateway.
type Server struct {
	store   *store.Store
	artRoot string
}

// NewServer wires a Server to its backing store and artifact directory.
func NewServer(s *store.Store, artRoot string) *Server {
	return &Server{store: s, artRoot: artRoot}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// CreateSnapshot records metadata for a newly uploaded artifact and, when an
// owner filter is supplied, returns that owner's existing snapshots so the
// client can de-duplicate.
func (s *Server) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	name := r.URL.Query().Get("name")
	if owner == "" || name == "" {
		http.Error(w, "owner and name required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<20))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	id, err := s.store.Insert(name, owner, int64(len(body)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	locale := preferredLocale(r.Header.Get("Accept-Language"))
	w.Header().Set("Content-Language", locale)

	// Return the owner's catalogue alongside the new id.
	existing, err := s.store.ListByOwner(owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "existing": existing})
}

// GetSnapshot streams a stored artifact back to the client.
func (s *Server) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	data, err := s.store.ReadArtifact(name)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(data)
}

// TrimMetadata hands a stored artifact to the native exiftrim helper, which
// strips EXIF/comment blocks in place and returns the reclaimed byte count.
func (s *Server) TrimMetadata(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	path := s.store.ArtifactPath(name)

	saved, err := native.RunExiftrim(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"reclaimed": saved})
}

// JobRequest is the body of POST /jobs.
type JobRequest struct {
	Input   string `json:"input"`   // artifact name to process
	Output  string `json:"output"`  // desired output artifact name
	Filter  string `json:"filter"`  // imagemagick-style convert filter chain
	Quality int    `json:"quality"` // output quality 1..100
}

// EnqueueJob forwards a processing request to the Rust worker unit.
func (s *Server) EnqueueJob(w http.ResponseWriter, r *http.Request) {
	var jr JobRequest
	if err := json.NewDecoder(r.Body).Decode(&jr); err != nil {
		http.Error(w, "bad job", http.StatusBadRequest)
		return
	}

	spec := native.JobSpec{
		Input:   s.store.ArtifactPath(jr.Input),
		Output:  jr.Output,
		Filter:  jr.Filter,
		Quality: strconv.Itoa(jr.Quality),
	}
	out, err := native.RunWorker(spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"result": out})
}
