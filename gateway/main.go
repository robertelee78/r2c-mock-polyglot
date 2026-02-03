// snapstore gateway — the front-end HTTP unit of the polyglot ingest service.
//
// Request lifecycle:
//   POST /snapshots?owner=<o>   -> api.CreateSnapshot -> store.Insert (+ owner lookup)
//   GET  /snapshots/{name}      -> api.GetSnapshot    -> store.ReadArtifact (serve file)
//   POST /snapshots/{name}/trim -> api.TrimMetadata   -> native.RunExiftrim (C helper)
//   POST /jobs                  -> api.EnqueueJob     -> native.RunWorker   (Rust worker)
//
// The gateway is the PRIMARY unit; it spawns the sibling first-party units
// `exiftrim` (C) and `worker` (Rust) as child processes — two cross-unit seams.
package main

import (
	"log"
	"net/http"

	"example.com/snapstore/gateway/internal/api"
	"example.com/snapstore/gateway/internal/store"
)

func main() {
	db, err := store.Open("snapstore.db")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer db.Close()

	srv := api.NewServer(db, "/var/lib/snapstore/artifacts")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /snapshots", srv.CreateSnapshot)
	mux.HandleFunc("GET /snapshots/{name}", srv.GetSnapshot)
	mux.HandleFunc("POST /snapshots/{name}/trim", srv.TrimMetadata)
	mux.HandleFunc("POST /jobs", srv.EnqueueJob)

	log.Println("snapstore gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
