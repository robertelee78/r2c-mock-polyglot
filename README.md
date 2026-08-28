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
(cd native/exiftrim && make)             # requires libcurl headers and library
(cd worker && cargo build)              # no external deps
```

There are no downloaded language-package dependencies. The native helper uses
the system C compiler and libcurl development package.

## Deliberately dangerous fixture modes

Normal helper invocations do not enter the advanced Stage 5 cases. The
`retry:` file prefix enters an interprocedural libcurl use-after-cleanup, and
the `maintenance:` job prefix enters a two-thread lock-order deadlock. The
`pagecache:<offset>:<four-bytes>:<path>` mode aliases a shared file mapping into
an in-place crypto destination and overwrites four bytes of that file. These
modes exist for scanner evaluation and should only be invoked in a disposable,
time-bounded process; the page-cache mode should target only a disposable file.
