# Graph Report - .  (2026-05-02)

## Corpus Check
- Corpus is ~26,368 words - fits in a single context window. You may not need a graph.

## Summary
- 447 nodes · 464 edges · 63 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS · INFERRED: 1 edges (avg confidence: 0.6)
- Token cost: 0 input · 0 output

## God Nodes (most connected - your core abstractions)
1. `mockFieldError` - 13 edges
2. `ConnectionService` - 12 edges
3. `GitHubClient` - 11 edges
4. `Repository[T]` - 10 edges
5. `gcpService` - 8 edges
6. `MockService` - 8 edges
7. `InMemFile` - 7 edges
8. `connectionRepository` - 7 edges
9. `GormLogger` - 6 edges
10. `RealGothProvider` - 5 edges

## Surprising Connections (you probably didn't know these)
- `Go Common README` --conceptually_related_to--> `Graph Report`  [INFERRED]
  README.md → graphify-out/GRAPH_REPORT.md

## Communities

### Community 0 - "Connection Service"
Cohesion: 0.1
Nodes (4): ConnectionService, ConnectionStorage, gcpService, Service

### Community 1 - "Authorization & JSON Response"
Cohesion: 0.09
Nodes (9): ErrorResponse, SendAutoErrorResponse(), SendErrorResponse(), Claims, AuthenticatedUser, RateLimitConfig, RateLimitStore, RateLimitMiddleware() (+1 more)

### Community 2 - "Auth Helpers & Errors"
Cohesion: 0.1
Nodes (8): GetAuthenticatedUser(), GetAuthenticatedUserID(), MembershipStore, mockRoundTripper, TestRetryRoundTripper(), RunSync(), SafeGo(), Task

### Community 3 - "GitHub Client"
Cohesion: 0.22
Nodes (5): GitHubClient, githubRepo, githubRepoWrapper, extractNextPageURL(), isInstallationContext()

### Community 4 - "JSON Validation Mocks"
Cohesion: 0.13
Nodes (3): mockFieldError, TestJsonError(), TestSendErrorResponse()

### Community 5 - "Config Structs"
Cohesion: 0.12
Nodes (14): BitbucketConfig, CorsConfig, DatabaseConfig, GitHubConfig, GitLabConfig, GoogleConfig, JWTConfig, Level (+6 more)

### Community 6 - "Database & OTel"
Cohesion: 0.14
Nodes (1): DBConfig

### Community 7 - "Crypto & Service Tests"
Cohesion: 0.14
Nodes (5): ContentItem, Namespace, ProviderClient, Repository, TokenRefreshResponse

### Community 8 - "Context Keys & Logging"
Cohesion: 0.19
Nodes (6): EnrichLogger(), GetLoggerFromContext(), GetTraceID(), CtxKey, LogState, TraceIDExtractor

### Community 9 - "Archive Buffering"
Cohesion: 0.17
Nodes (2): BufferFile, InMemFile

### Community 10 - "OAuth Providers"
Cohesion: 0.17
Nodes (5): GothConfig, GothInitializer, GothProvider, OAuthProviderInitializer, RealGothProvider

### Community 11 - "Connection Repository"
Cohesion: 0.18
Nodes (2): connectionRepository, Repository

### Community 12 - "Generic GORM Repository"
Cohesion: 0.27
Nodes (1): Repository[T]

### Community 13 - "Server Lifecycle"
Cohesion: 0.22
Nodes (4): AddGracefulShutdownGoroutine(), DoneGracefulShutdownGoroutine(), Server, Server

### Community 14 - "Logger & Request Logging"
Cohesion: 0.24
Nodes (4): getRealIP(), RequestLogger(), scrubPayload(), responseWriter

### Community 15 - "GCP Mock Service"
Cohesion: 0.2
Nodes (1): MockService

### Community 16 - "GORM Logger"
Cohesion: 0.31
Nodes (2): GormLogger, Config

### Community 17 - "Env Config Loading"
Cohesion: 0.25
Nodes (2): GetEnvAsHexBytes(), GetRequiredEnv()

### Community 18 - "Telemetry Multi-Tracker"
Cohesion: 0.22
Nodes (3): Event, MultiTracker, Tracker

### Community 19 - "GA4 Analytics"
Cohesion: 0.25
Nodes (3): ga4Event, ga4Payload, GA4Tracker

### Community 20 - "Router"
Cohesion: 0.25
Nodes (2): Config, Router

### Community 21 - "CSRF & URL Utils"
Cohesion: 0.32
Nodes (1): CSRFConfig

### Community 22 - "Middleware Tests"
Cohesion: 0.25
Nodes (0): 

### Community 23 - "Env Config Tests"
Cohesion: 0.25
Nodes (0): 

### Community 24 - "Telemetry Mocks & Tests"
Cohesion: 0.32
Nodes (3): mockTracker, TestGA4Tracker(), TestMultiTracker()

### Community 25 - "HTTP Transport"
Cohesion: 0.29
Nodes (2): ResponseWriterWrapper, RetryRoundTripper

### Community 26 - "Request Binder"
Cohesion: 0.43
Nodes (2): DecodeAndValidate(), Binder

### Community 27 - "Git Cloner & Exec"
Cohesion: 0.29
Nodes (2): DefaultGitCloner, GitCloner

### Community 28 - "PostHog Tracker"
Cohesion: 0.33
Nodes (1): PostHogTracker

### Community 29 - "Metrics Middleware"
Cohesion: 0.4
Nodes (1): loggingResponseWriter

### Community 30 - "Connection Storage"
Cohesion: 0.6
Nodes (2): ETag(), responseBuffer

### Community 31 - "Tracing"
Cohesion: 0.4
Nodes (3): DiscordEmbed, DiscordEmbedField, DiscordWebhook

### Community 32 - "Worker Queue"
Cohesion: 0.5
Nodes (1): TestModel

### Community 33 - "JWT Middleware"
Cohesion: 0.5
Nodes (2): ExternalConnection, BaseModel

### Community 34 - "Slug Utils"
Cohesion: 0.5
Nodes (1): txKey

### Community 35 - "Request ID"
Cohesion: 0.5
Nodes (0): 

### Community 36 - "Rate Limiting"
Cohesion: 0.5
Nodes (1): TestUser

### Community 37 - "ETag Middleware"
Cohesion: 0.5
Nodes (0): 

### Community 38 - "Trailing Slash"
Cohesion: 0.5
Nodes (0): 

### Community 39 - "Cookie Middleware"
Cohesion: 0.67
Nodes (2): PaginatedResponse, PaginationParams

### Community 40 - "CORS Middleware"
Cohesion: 1.0
Nodes (1): CorsConfig

### Community 41 - "Recovery Middleware"
Cohesion: 0.67
Nodes (0): 

### Community 42 - "JWT Token"
Cohesion: 0.67
Nodes (0): 

### Community 43 - "Domain Utils"
Cohesion: 0.67
Nodes (0): 

### Community 44 - "GCP Storage"
Cohesion: 0.67
Nodes (0): 

### Community 45 - "GCP Cloud Run"
Cohesion: 0.67
Nodes (0): 

### Community 46 - "GCP Cloud Build"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "Connection Types"
Cohesion: 1.0
Nodes (0): 

### Community 48 - "Discord Notify"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "Pagination DTO"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "GORM Transaction"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "GORM Model"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "Test DB Setup"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "Test Helpers"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "Random String"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "Exec Tests"
Cohesion: 1.0
Nodes (0): 

### Community 56 - "Binder Tests"
Cohesion: 1.0
Nodes (0): 

### Community 57 - "HTTP Transport Tests"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "URL Tests"
Cohesion: 1.0
Nodes (0): 

### Community 59 - "Postgres Tests"
Cohesion: 1.0
Nodes (0): 

### Community 60 - "GitHub Client Tests"
Cohesion: 1.0
Nodes (2): Graph Report, Go Common README

### Community 61 - "Repository Tests"
Cohesion: 1.0
Nodes (0): 

### Community 62 - "Connection Repo Tests"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **52 isolated node(s):** `ErrorResponse`, `PaginationParams`, `PaginatedResponse`, `Claims`, `CSRFConfig` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `GCP Cloud Build`** (2 nodes): `recovery.go`, `Recovery()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Connection Types`** (2 nodes): `request_id_test.go`, `TestRequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Discord Notify`** (2 nodes): `logging_test.go`, `TestScrubPayload()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Pagination DTO`** (2 nodes): `trailing_slash_test.go`, `TestTrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `GORM Transaction`** (2 nodes): `csrf_test.go`, `TestCSRFMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `GORM Model`** (2 nodes): `request_id.go`, `RequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Test DB Setup`** (2 nodes): `trailing_slash.go`, `TrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Test Helpers`** (2 nodes): `goth_initializer_test.go`, `TestGothInitializer()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Random String`** (2 nodes): `random.go`, `RandomString()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Exec Tests`** (2 nodes): `slug_test.go`, `TestSlugify()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Binder Tests`** (2 nodes): `server_test.go`, `TestServer_StartAndShutdown()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Transport Tests`** (2 nodes): `executil_test.go`, `TestCommandContext()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `URL Tests`** (2 nodes): `domain.go`, `IsValidSubdomain()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Postgres Tests`** (2 nodes): `db_test.go`, `TestSetupTestDatabase()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `GitHub Client Tests`** (2 nodes): `Graph Report`, `Go Common README`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Repository Tests`** (1 nodes): `run.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Connection Repo Tests`** (1 nodes): `build.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `ErrorResponse`, `PaginationParams`, `PaginatedResponse` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Connection Service` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Authorization & JSON Response` be split into smaller, more focused modules?**
  _Cohesion score 0.09 - nodes in this community are weakly interconnected._
- **Should `Auth Helpers & Errors` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `JSON Validation Mocks` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Config Structs` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._
- **Should `Database & OTel` be split into smaller, more focused modules?**
  _Cohesion score 0.14 - nodes in this community are weakly interconnected._