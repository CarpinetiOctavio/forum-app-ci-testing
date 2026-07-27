# Context — why this report exists and what's expected from this chat

This analysis didn't originate from work done in forum-app-ci-testing. It came
up while auditing forum-app-qa-pipeline-legacy (mirror of TP7, a different repo):
while investigating a specific coverage gap in DeleteComment, something more
fundamental showed up — the application's authentication model is spoofable by
design (no real token/session, any request authenticates with an X-User-ID
header that the client sets with no verification). It was checked whether this
was specific to that old mirror or also present here, in forum-app-ci-testing
(already closed, tag v1.0.0) — it's in both, with the same code, almost line
for line. It's not something introduced by qa-pipeline; it's in the app's
source code that both repos share.

Why it's brought to this chat instead of being resolved directly in
qa-pipeline: the new qa-pipeline repo is going to be built as a copy of this
repo at its v1.0.0 tag. If the problem is fixed over there without touching it
here, the series ends up with a strange inconsistency — the earlier, "closed"
repo showing the flaw, a later one showing the fixed version. The series'
correction methodology already establishes that each repo is grounded in its
own scope and fixes travel forward, never backward — it makes sense for this
to be resolved here first, with its own ADR, so that the copy made afterward
already inherits it fixed.

This document is the second audit pass (after an initial brief) run directly
against this repo's code — it confirms the 4 original findings with its own
evidence, corrects two points that were actually already well resolved, and
adds two new findings. It's diagnostic input, not a decision made: how each
thing gets fixed (real authentication, payload limits, pagination, etc.), with
what scope and what alternatives were evaluated, gets decided in this chat and
documented with the same rigor as the rest of the repo.

# Audit results — app design/security in forum-app-ci-testing

Inventory produced following `app-security-audit-brief.md`. Everything here is input for Octavio to decide on in the ci-testing chat; nothing was applied, no file in `forum-app-ci-testing` was modified. Verified directly against that repo's working tree (confirmed identical in content to tag `v1.0.0`: `git diff v1.0.0 HEAD --stat` shows no differences, though it diverges in history).

---

## 0. Executive summary

- The 4 findings already confirmed in the brief (no real auth, plaintext passwords, CORS `*`, exposed internal errors) were re-verified and **remain accurate**, identical code to `qa-pipeline-legacy`.
- **One correction to the brief**: point 2 (is `Password` exposed in the `User` JSON?) — no, it's already protected with `json:"-"`. Not a new finding, it's a confirmation that this specific point is already well resolved.
- **The highest-weight new finding**: no endpoint that receives a body (`CreatePost`, `CreateComment`, `Register`, `Login`) limits payload size (`r.Body` is never wrapped in `http.MaxBytesReader`) — a request with a body of several GB is attempted to be decoded in full. It's in the same style as the 4 already confirmed (code that works but wasn't designed with judgment about what it allows), though of lower severity than authentication.
- The database schema (cascades, indexes) is **well designed** — not a finding, it's a correction to item 1 of the brief's pending list (no data inconsistency from deletion without cascade).
- The frontend doesn't persist anything sensitive in `localStorage`/`sessionStorage` — everything lives in transient React state, lost on refresh. Not a finding.
- Classic CSRF doesn't apply (the reasoning was confirmed), but not because the design prevents it — it's made irrelevant by an already-worse problem: CSRF isn't needed when anyone can call the API directly thanks to CORS `*` + an unverified `X-User-ID`.
- Dependency vulnerabilities: 55 in the frontend (2 critical, both in Create React App dev tooling, not in the production bundle) — same picture already seen in `qa-pipeline-legacy` (58). In the backend, `govulncheck` flags 18 "reachable," but all 18 are from the Go standard library tied to the local toolchain version (go1.25.0), not from the app's code or its 3 direct dependencies (`gorilla/mux`, `go-sqlite3`, `testify`, with 0 reachable vulnerabilities).

---

## 1. Re-verification of the 4 already confirmed findings

All re-confirmed, identical code to `qa-pipeline-legacy` (same `forum-app-ci-testing` module, same relative paths):

| # | Finding | Evidence in ci-testing |
|---|---|---|
| 1 | No real authentication | `router.go:16` (only middleware, CORS); `auth_service.go` `Login` doesn't issue a token; the 4 mutation handlers in `post_handler.go` read `X-User-ID` from the header with no verification (`CreatePost:36`, `DeletePost:99`, `CreateComment:139`, `DeleteComment:194`) |
| 2 | Plaintext passwords | `auth_service.go` — direct comparison/storage, comments "in production: hash with bcrypt" unimplemented |
| 3 | Maximally permissive CORS | `router.go:40-42` — `Access-Control-Allow-Origin: *`, `X-User-ID` explicit in `Access-Control-Allow-Headers` |
| 4 | Internal errors exposed | `post_handler.go:62` — `GetAllPosts` responds with `http.StatusInternalServerError, err.Error()` raw |

---

## 2. Sweep of the 6 pending points

### 2.1 Database schema (`database.go`) — **no finding, correct design**

```sql
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE   -- posts
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE   -- comments
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE   -- comments
CREATE INDEX idx_posts_user_id, idx_comments_post_id, idx_comments_user_id
```

Complete cascades on both child tables, indexes on all 3 FK columns. Deleting a user or a post leaves no orphaned comments/posts. **Correction to the brief**: there's no design finding here, the schema is well thought out.

### 2.2 Models (`post.go`, `users.go`) — **no finding, `Password` already protected**

```go
Password  string    `json:"-"` // Not serialized in JSON (for security)
```
(`users.go:9`). The explicit comment ("for security") suggests it was a conscious decision, not an accident. **Correction to the brief**: point 2 was left as an open question, the answer is that it's already well resolved — the password is not leaked (neither plaintext nor hashed, since it isn't hashed either) in any JSON response.

### 2.3 Dependencies with known vulnerabilities

**Backend (`govulncheck`)**: installed the tool (`go install golang.org/x/vuln/cmd/govulncheck@latest`) and ran `govulncheck ./...`. Result: **18 "reachable" vulnerabilities** from the code, but all 18 are from the Go **standard library** (`crypto/tls`, `crypto/x509`, `net/url`, `encoding/asn1`, `net/textproto`), not from `gorilla/mux`, `go-sqlite3`, or `testify` (0 reachable vulnerabilities in the 3 direct dependencies). Each one is tied to the locally installed Go toolchain version (`go1.25.0`) and is resolved simply by compiling with a newer patch (`go1.25.2`/`go1.25.3`) — it's not a design problem in the app or its chosen dependencies, it's compiler freshness on the machine/CI that builds it. Severity: low, informational — resolved with `go install`/updating the CI runner, not with a code change.

**Frontend (`npm audit`)**: 55 total vulnerabilities (11 low, 14 moderate, 28 high, **2 critical**: `shell-quote`, `websocket-driver`). Confirmed the dependency chain of the 2 critical ones — both are **transitive**, buried in Create React App's dev tooling (`react-scripts`/`webpack-dev-server`), not in the bundle served in production. Same picture already seen in `qa-pipeline-legacy` (58 vulnerabilities) — confirms this isn't something that changed when it was rebuilt, it's inherent to Create React App being deprecated. Severity: low for a portfolio (dev-tooling surface, not production), but 2 critical + 28 high is a long list if someone looks at it without that context.

### 2.4 Frontend — storage of sensitive data — **no finding**

`grep` for `localStorage`/`sessionStorage` across all of `frontend/src` (excluding tests): **zero results**. `App.tsx` stores `currentUser` in root-level component `useState` — it lives only in memory, lost on page refresh (requires logging in again). `Login.tsx` stores `password` in local form `useState` — never persisted to any storage, and the component unmounts (losing that state) as soon as login succeeds. `X-User-ID` is built on each request from `currentUser.id` (a prop passed down from `App.tsx`, `postService.ts`), never saved to `localStorage`. **No finding here** — nothing sensitive survives longer than necessary on the client. (Side effect unrelated to security: the session doesn't persist across refreshes, which is more of a UX limitation than a design one.)

### 2.5 Missing pagination / uncapped payloads — **finding confirmed and expanded**

- **No pagination on 2 endpoints, not just 1**: `GetAllPosts` (`post_handler.go:59-67`) and `GetComments` (`post_handler.go:162-178`) return everything with no limit. Also confirmed there's no `LIMIT`/`OFFSET` in any repository query (`grep` over `post_repository.go` with no results) — not even a hard cap at the data layer.
- **New finding, not anticipated by the brief**: no handler that decodes a body (`CreatePost`, `CreateComment`, `Register`, `Login`) limits the request size. All 4 use `json.NewDecoder(r.Body).Decode(&req)` directly on `r.Body` without wrapping it in `http.MaxBytesReader(w, r.Body, limit)` (confirmed with `grep -rn "MaxBytesReader\|ContentLength"` over all of `backend/`, no results). A request with a several-GB body in a post or comment's `content` field is attempted to be read/parsed in full before any field-length validation fails. It's the same pattern as "no pagination" (uncapped resources) but on the input side instead of the output side.

**Severity**: lower than authentication/passwords, but it's a concrete finding with code evidence, not hypothetical — classified alongside point 4 of the original brief's logging issue (same caliber: functional code that wasn't designed with judgment about what it allows).

### 2.6 CSRF — reasoning confirmed, with a nuance

Confirmed: the classic CSRF vector (a `<form>` or `<img>` on a malicious site that triggers a request and the browser automatically attaches a session cookie) **does not apply** here because there are no session cookies — `X-User-ID` has to be explicitly set by JavaScript (`postService.ts`), and a "passive" CSRF request (form submit, img tag) can't set custom headers.

**Nuance to add** (as the brief asks): it's not that the design defends against CSRF — it's that CSRF becomes unnecessary because the real problem is already worse. Since `Access-Control-Allow-Origin: *` allows JavaScript from **any origin** to call the API and read the response, and `X-User-ID` isn't verified against anything, an attacker doesn't need the victim at all: they can call the API directly with any `user_id`. And that `user_id` doesn't even need to be guessed — it's in plaintext in the `user_id` field of every post returned by `GET /api/posts` (`post.go:10`, unprotected). This nuance isn't an independent new finding — it's a direct extension of the already-confirmed finding #1, showing that it makes worrying about CSRF as a separate vector unnecessary.

### 2.7 Other findings that came up along the way

- **Complete and well-placed FK indexes** (see 2.1) — mentioned as a correction, not as a negative finding.
- No other pattern from the same family (functional code without professional judgment) showed up beyond what's already listed in 2.5. Reviewed in full: `auth_handler.go`, `post_handler.go`, `router.go`, `database.go`, both models, `postService.ts`, `authService.ts`, `Login.tsx`, `App.tsx` — found nothing additional to report with concrete evidence.

---

## 3. Relative severity table

| Finding | Severity vs. the 4 confirmed |
|---|---|
| No real authentication (original #1) | Reference point — the most severe |
| Plaintext passwords (original #2) | Same caliber as #1 |
| CORS `*` (original #3) | Minor, aggravates #1 but isn't independent |
| Internal errors exposed (original #4) | Minor, informational |
| **Uncapped payload size (new)** | Similar caliber to #4 — concrete, with evidence, but exploitation requires effort/volume |
| **No pagination on GetAllPosts/GetComments (new, expanded)** | Similar to uncapped payloads — same type of risk (resource exhaustion), lower urgency given the academic data volume |
| Frontend dependencies (2 critical, transitive in dev-tooling) | Low — doesn't touch the production bundle |
| Backend dependencies via `govulncheck` (18, all stdlib/toolchain) | Informational — not a problem with the code or its 3 direct dependencies |
| DB schema, `User.Password` model, frontend storage | **No finding** — verified as correct, no action needed |

---

## 4. Corrections to the brief's "already verified as correct" list

To add to that list (not confirmed before this sweep):
- The `Password` field of the `User` model has `json:"-"` — not leaked in any response.
- The database schema has complete cascades and indexes — no risk of orphaned data.
- The frontend doesn't persist anything sensitive in `localStorage`/`sessionStorage` — all auth state is transient in React memory.
