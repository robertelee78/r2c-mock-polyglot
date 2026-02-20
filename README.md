# snapstore

A small artifact-ingest service. The Go gateway is the front end; it delegates
work to three sibling helper units, each spawned as a child process with a
direct argument vector.

## Components

| Unit | Language | Entry point | Invoked by |
|------|----------|-------------|------------|
| gateway (primary) | Go | `gateway/main.go` | — |
| exiftrim | C | `native/exiftrim/src/main.c` | `GET /trim?file=` |
| worker | Rust | `worker/src/main.rs` | `POST /jobs?name=` |
| report | Python | `report/report.py` | `GET /reports?owner=` |

## HTTP surface

```
GET  /import?url=<u>      fetch and stream a remote artifact
GET  /trim?file=<f>       strip metadata via the exiftrim helper
POST /jobs?name=<n>       run a worker job
GET  /reports?owner=<o>   owner report via the python helper
```

## Building

```sh
(cd gateway && go build ./...)          # stdlib only, no cgo
(cd native/exiftrim && make)
(cd worker && cargo build)              # no external deps
```

All four units build in well under a second; there are no downloaded
dependencies and no native toolchain beyond a C compiler.
