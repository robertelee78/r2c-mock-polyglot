# snapstore — intentionally vulnerable polyglot target (Go gateway + C/Rust/Python helpers).
# Authored for DAST testing; upstream ships no Dockerfile. The gateway invokes exiftrim/worker by
# bare name (must be on PATH) and report/report.py by a path relative to the working dir (/app).

# --- Go gateway (static) + C exiftrim ---
FROM golang:1.22-bookworm AS gobuild
WORKDIR /src
COPY . .
# The gateway binds 127.0.0.1, which is container-local and unreachable through Docker port-mapping. Rewrite
# it to 0.0.0.0 for the image only, so the repo's gateway/main.go stays untouched (this Dockerfile is additive).
RUN sed -i 's/127\.0\.0\.1:8080/0.0.0.0:8080/' gateway/main.go \
 && mkdir -p /out \
 && CGO_ENABLED=0 go build -o /out/gateway ./gateway \
 && make -C native/exiftrim \
 && cp native/exiftrim/exiftrim /out/exiftrim

# --- Rust worker ---
FROM rust:1-bookworm AS rustbuild
WORKDIR /src/worker
COPY worker/ ./
RUN cargo build --release && mkdir -p /out && cp target/release/worker /out/worker

# --- runtime ---
FROM python:3.12-slim-bookworm
WORKDIR /app
COPY report/ /app/report/
COPY --from=gobuild /out/gateway  /usr/local/bin/gateway
COPY --from=gobuild /out/exiftrim /usr/local/bin/exiftrim
COPY --from=rustbuild /out/worker /usr/local/bin/worker
# Seed the sqlite DB the Python /reports helper queries (so SQLi is reachable + boolean-differential works),
# and create the Rust worker's spool dir (so /jobs path-traversal resolves through an existing directory).
RUN python3 - <<'PY'
import sqlite3
c = sqlite3.connect("/app/snapstore.db")
c.execute("CREATE TABLE snapshots(name TEXT, created_at TEXT, owner TEXT)")
c.executemany("INSERT INTO snapshots VALUES (?,?,?)",
    [("alpha","2024-01-01","alice"),("beta","2024-02-02","bob"),("secret-key","2024-03-03","admin")])
c.commit(); c.close()
PY
RUN mkdir -p /var/lib/snapstore/jobs && printf 'demo job\n' > /var/lib/snapstore/jobs/pending.job
EXPOSE 8080
CMD ["gateway"]
