# Snapstore security ground truth

This repository is an intentionally vulnerable test target. It is meant to be
scanned as a standalone client repository; none of its weaknesses should be
treated as examples, tests, or out-of-scope code.

| ID | Unit | Defect | CWE | Expected nomination | Attacker path |
|---|---|---|---|---|---|
| SNAP-001 | Go gateway | Server-side request forgery in `handleRequest` | CWE-918 | Stage 3 static SAST | `GET /import?url=...` to `http.Get` |
| SNAP-002 | C exiftrim | Stack buffer overflow in `trim_metadata` | CWE-120 | Stage 3 static SAST; Stage 5 may add `userland-c-memory` context | `GET /trim?file=...` to `strcpy` through the child-process boundary |
| SNAP-003 | Rust worker | Path traversal in `process_job` | CWE-22 | Stage 3 static SAST | `POST /jobs?name=...` to `fs::read_to_string` through the child-process boundary |
| SNAP-004 | Python report | SQL injection in `lookup_owner` | CWE-89 | Stage 3 static SAST | `GET /reports?owner=...` to `sqlite3.Connection.execute` through the child-process boundary |
| SNAP-005 | C exiftrim | libcurl easy-handle use after cleanup | CWE-416 | Stage 5 `backend-network-crypto` | `GET /trim?file=retry:...` through cleanup, cross-function transfer, and reuse |
| SNAP-006 | Rust worker | ABBA lock-order deadlock | CWE-833 | Stage 5 concurrency seam | `POST /jobs?name=maintenance:...` starts both lock-taking threads and blocks the worker |
| SNAP-007 | C exiftrim | file-backed page aliased into an in-place crypto destination | CWE-73 / CWE-669 | Future kernel/crypto page-provenance plugin; current product gap | `GET /trim?file=pagecache:<offset>:<four-bytes>:<path>` selects the file, offset, and scratch-write value |

Expected analysis units:

1. `primary` — Go, repository root (including `gateway/`)
2. `unit--native--exiftrim` — C, `native/exiftrim/`
3. `unit--report` — Python, `report/`
4. `unit--worker` — Rust, `worker/`

The three Go `exec.Command` calls deliberately use direct argument vectors.
They are cross-unit boundaries, not command-injection findings.

## SNAP-005 — libcurl handle reused across cleanup boundary

- **Files and functions:** `native/exiftrim/src/main.c`:
  `retry_metadata_transfer`, `restart_transfer`, and `perform_transfer`.
  `main` is dispatch context.
- **Defective anti-pattern:** `restart_transfer` calls
  `curl_easy_cleanup(handle)` and then passes that same retired `handle` to
  `perform_transfer`, which calls `curl_easy_perform(handle)`. This is a
  real use of a libcurl easy handle after its lifetime ended. It is inert
  unless the helper receives the explicit `retry:` prefix.
- **Why isolated review is insufficient:** the source function contains the
  lifetime-ending call but not the libcurl use, while the sink function
  contains the use but not the cleanup. The defect exists only in the ordered
  relation between them.
- **Expected Stage 5 recognition:** built-in plugin
  `backend-network-crypto`; applicability comes from the real
  `curl/curl.h` integration. The exact shape is
  `backend-network-crypto:libcurl-easy-handle-lifecycle`, confirmed as a
  `FreeUse` engine path from `curl_easy_cleanup(handle)` through the resolved
  call to `curl_easy_perform(handle)`.
- **Expected participants:** `restart_transfer` and `perform_transfer`.
  `retry_metadata_transfer` and `main` are
  reachability/dispatch context only and must not replace either participant.
- **Expected Stage 6 review unit:** one ordered lifetime review unit containing
  both participants with a `JointReviewRequired`-style relationship. The
  engine-confirmed path must preserve cleanup, call, and use topology.
- **Known current handoff mismatch:** the current Stage 5 substrate carries
  both functions in `engine_confirmed_paths`, but exact candidate projection
  promotes only a source-anchored `RecognizedSite`. Consequently current
  `plugin_candidates.json` is expected to mark only
  `restart_transfer` as context and omit the joint relationship. That is a
  product gap exposed by this fixture, not an alternate expected result.
- **Why it earns its cost:** two small functions exercise real domain ownership
  semantics that ordinary generic SAST is unlikely to model faithfully. It
  also tests both typestate confirmation and lossless Stage 6 handoff. Its only
  new build cost is the system libcurl development package.

## SNAP-006 — cross-function lock-order inversion

- **Files and functions:** `worker/src/main.rs`: `refresh_index` and
  `rotate_spool`. `run_parallel_maintenance` and `process_job` are
  launch/reachability context.
- **Defective anti-pattern:** `refresh_index` acquires `INDEX_LOCK` before
  `SPOOL_LOCK`; `rotate_spool` acquires the same locks in reverse order.
  A two-party barrier makes the explicit `maintenance:` path reliably hold each
  first lock before either requests the second, blocking both threads and the
  worker process.
- **Why isolated review is insufficient:** each function is internally
  consistent and has no missing unlock. The deadlock is visible only by
  comparing both functions' lock-order edges over the same lock pair.
- **Expected Stage 5 recognition:** the core Stage 5 concurrency-seam detector
  recognizes the `.lock().unwrap()` acquisitions. The exact structural shape
  is two `ConcurrencySeam` items, one for `INDEX_LOCK` and one for
  `SPOOL_LOCK`, with the same actor pair and opposing lock-order edges.
- **Expected participants:** `refresh_index` and `rotate_spool`. There
  are no context-only members in the concurrency relationship;
  `run_parallel_maintenance` and `process_job` are upstream context only.
- **Expected Stage 6 review unit:** one transitive review unit containing both
  participant functions. Each seam must carry a `JointReviewRequired`
  concurrency relationship; their shared participant pair should merge into
  the same Stage 6 unit without dropping either function.
- **Why it earns its cost:** it uses only the Rust standard library, is inert
  without the explicit prefix, and proves the multi-function participant and
  topology behavior that distinguishes Stage 5 from isolated function review.

## SNAP-007 — file-backed page transferred into a writable crypto sphere

- **Files and functions:** `native/exiftrim/src/main.c`:
  `map_file_crypto_source`, `choose_in_place_destination`, and
  `finish_authenticated_decrypt`. `decrypt_mapped_metadata`, `main`, and the
  Go gateway's `handleRequest` are reachability and parsing context.
- **Defective anti-pattern:** the unauthenticated caller selects a service-
  writable file, byte offset, and four-byte value. `map_file_crypto_source`
  creates a real `MAP_SHARED` file mapping. The zero-copy optimization in
  `choose_in_place_destination` aliases those file-backed source pages as the
  crypto destination without checking their provenance. Crypto completion
  then uses that destination as scratch space and flushes the four-byte write
  into the backing file. This is a deterministic integrity violation across
  the caller/service privilege boundary, not a simulated marker.
- **Why isolated review is insufficient:** mapping a file as a crypto source,
  selecting an in-place destination, and writing a four-byte sequence scratch
  value are individually plausible operations. The bug exists only when the
  source's file-backed provenance is carried through the alias decision to the
  write sink.
- **Expected Stage 5 recognition:** no current built-in plugin recognizes this
  relationship. It is an intentional future kernel/crypto domain-plugin target
  with the exact structural shape `SharedFilePageSource` to
  `InPlaceDestinationAlias` to `CryptoScratchWrite`, retaining the selected
  offset and the alias edge. This expected gap must not be reported as a
  successful current nomination.
- **Expected participants:** `map_file_crypto_source` is the resource-
  provenance producer, `choose_in_place_destination` is the unsafe transfer/
  alias participant, and `finish_authenticated_decrypt` is the write
  participant. `decrypt_mapped_metadata`, `main`, and `handleRequest` are
  context only.
- **Expected Stage 6 review unit:** one ordered joint unit containing all three
  participants, with the shared-file mapping, source/destination alias, and
  four-byte write connected in the handoff. None of the three may be replaced
  by a context-only caller.
- **Public-vulnerability influence:** this is a safe userland analogue of the
  resource-provenance failure in [Copy Fail (CVE-2026-31431)](https://xint.io/blog/copy-fail-linux-distributions):
  a zero-copy source retains file-page identity, an in-place crypto design
  moves it into a writable destination sphere, and a four-byte scratch write
  corrupts the page. It deliberately does not contain AF_ALG exploit code or
  target privileged system binaries.
- **Why it earns its cost:** it adds a distinct domain relationship—page
  provenance crossing a crypto-buffer ownership boundary—that neither the
  libcurl lifetime case nor the lock topology covers. It is portable POSIX C,
  adds no dependency, and is safely demonstrable using a disposable file.

## Independent reality check

On 2026-08-27, OpenCode 1.18.23 with `opencode/kimi-k3` reviewed only the C
and Rust implementation files; this document and the defect descriptions were
withheld. It independently identified SNAP-005 as a cross-function CWE-416
use-after-cleanup and SNAP-006 as a deterministic cross-function CWE-833/CWE-667
deadlock, rating both as real advanced security defects. It also distinguished
the pre-existing simple `strcpy` and path-traversal cases and rejected unrelated
code as false positives.

After SNAP-007 was added, a fresh blind review gave it the firm verdict **REAL
SECURITY DEFECT**. Kimi traced `argv` through the attacker-selected path,
`O_RDWR` open, `MAP_SHARED` source, in-place destination alias, four-byte write,
and `msync`. It classified the concrete arbitrary-file-write interface as
CWE-73 while confirming that the security property is visible only through the
cross-function backing-store and alias topology.
