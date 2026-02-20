# Snapstore security ground truth

This repository is an intentionally vulnerable test target. It is meant to be
scanned as a standalone client repository; none of its weaknesses should be
treated as examples, tests, or out-of-scope code.

| ID | Unit | Defect | CWE | Attacker path |
|---|---|---|---|---|
| SNAP-001 | Go gateway | Server-side request forgery in `handleRequest` | CWE-918 | `GET /import?url=...` to `http.Get` |
| SNAP-002 | C exiftrim | Stack buffer overflow in `trim_metadata` | CWE-120 | `GET /trim?file=...` to `strcpy` through the child-process boundary |
| SNAP-003 | Rust worker | Path traversal in `process_job` | CWE-22 | `POST /jobs?name=...` to `fs::read_to_string` through the child-process boundary |
| SNAP-004 | Python report | SQL injection in `lookup_owner` | CWE-89 | `GET /reports?owner=...` to `sqlite3.Connection.execute` through the child-process boundary |

Expected analysis units:

1. `primary` — Go, repository root (including `gateway/`)
2. `unit--native--exiftrim` — C, `native/exiftrim/`
3. `unit--report` — Python, `report/`
4. `unit--worker` — Rust, `worker/`

The three Go `exec.Command` calls deliberately use direct argument vectors.
They are cross-unit boundaries, not command-injection findings.

