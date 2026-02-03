# snapstore

A small artifact-ingest service. Clients upload binary snapshots; the service
records their metadata, serves them back, strips embedded EXIF/comment blocks,
and runs an image-processing pipeline over them.

## Components

| Unit | Language | Role |
|------|----------|------|
| `gateway/` | Go | HTTP front end — metadata store (sqlite), artifact serving, job intake. Spawns the two helper units below. |
| `native/exiftrim/` | C | Reads a snapshot's on-disk header and strips its metadata tags/comment in place. Invoked per artifact by the gateway. |
| `worker/` | Rust | Applies an operator-supplied convert filter chain to an artifact. Invoked per job by the gateway. |

## HTTP surface

```
POST /snapshots?owner=<o>&name=<n>   upload + record metadata
GET  /snapshots/{name}               fetch a stored artifact
POST /snapshots/{name}/trim          strip metadata via exiftrim
POST /jobs                           enqueue a convert job for the worker
```

## Building

```sh
# gateway (needs cgo for the sqlite driver)
cd gateway && go build ./...

# native helper
cd native/exiftrim && make && make install

# worker
cd worker && cargo build --release
```

The gateway locates the compiled helpers via `$SNAPSTORE_HELPERS`
(default `/usr/local/libexec/snapstore`).

## Layout

```
gateway/
  main.go
  internal/api/handlers.go     request handlers
  internal/store/db.go         sqlite metadata + artifact blobs
  internal/native/invoke.go    child-process bridge to the helper units
native/exiftrim/
  src/exiftrim.c               entry point + file IO + audit log
  src/parse.c                  header / tag / comment parsing
  include/parse.h
worker/
  src/main.rs                  arg parsing
  src/jobs.rs                  job execution
```
