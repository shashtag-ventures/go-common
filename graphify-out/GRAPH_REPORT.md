# Graph Report - .  (2026-04-12)

## Corpus Check
- Corpus is ~24,540 words - fits in a single context window. You may not need a graph.

## Summary
- 427 nodes · 433 edges · 60 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## God Nodes (most connected - your core abstractions)
1. `mockFieldError` - 13 edges
2. `Repository[T]` - 10 edges
3. `IntegrationService` - 10 edges
4. `GitHubClient` - 10 edges
5. `gcpService` - 8 edges
6. `MockService` - 8 edges
7. `InMemFile` - 7 edges
8. `GormLogger` - 6 edges
9. `integrationRepository` - 6 edges
10. `RealGothProvider` - 5 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Communities

### Community 0 - "HTTP Middleware"
Cohesion: 0.09
Nodes (9): ErrorResponse, SendAutoErrorResponse(), SendErrorResponse(), Claims, AuthenticatedUser, RateLimitConfig, RateLimitStore, RateLimitMiddleware() (+1 more)

### Community 1 - "Authentication & Security"
Cohesion: 0.1
Nodes (8): GetAuthenticatedUser(), GetAuthenticatedUserID(), MembershipStore, mockRoundTripper, TestRetryRoundTripper(), RunSync(), SafeGo(), Task

### Community 2 - "Cloud Services (GCP)"
Cohesion: 0.1
Nodes (4): gcpService, Service, IntegrationService, IntegrationStorage

### Community 3 - "HTTP Middleware"
Cohesion: 0.15
Nodes (3): GitHubClient, extractNextPageURL(), CSRFConfig

### Community 4 - "Internal Testing Utilities"
Cohesion: 0.13
Nodes (3): mockFieldError, TestJsonError(), TestSendErrorResponse()

### Community 5 - "Database & Persistence"
Cohesion: 0.12
Nodes (14): BitbucketConfig, CorsConfig, DatabaseConfig, GitHubConfig, GitLabConfig, GoogleConfig, JWTConfig, Level (+6 more)

### Community 6 - "Cloud Services (GCP)"
Cohesion: 0.14
Nodes (5): ContentItem, IntegrationClient, Namespace, Repository, TokenRefreshResponse

### Community 7 - "Database & Persistence"
Cohesion: 0.18
Nodes (1): DBConfig

### Community 8 - "Community 8"
Cohesion: 0.17
Nodes (2): BufferFile, InMemFile

### Community 9 - "Authentication & Security"
Cohesion: 0.17
Nodes (5): GothConfig, GothInitializer, GothProvider, OAuthProviderInitializer, RealGothProvider

### Community 10 - "Database & Persistence"
Cohesion: 0.27
Nodes (1): Repository[T]

### Community 11 - "Community 11"
Cohesion: 0.22
Nodes (4): AddGracefulShutdownGoroutine(), DoneGracefulShutdownGoroutine(), Server, Server

### Community 12 - "HTTP Middleware"
Cohesion: 0.24
Nodes (4): getRealIP(), RequestLogger(), scrubPayload(), responseWriter

### Community 13 - "HTTP Middleware"
Cohesion: 0.24
Nodes (4): EnrichLogger(), GetLoggerFromContext(), CtxKey, LogState

### Community 14 - "Database & Persistence"
Cohesion: 0.2
Nodes (2): Repository, integrationRepository

### Community 15 - "Cloud Services (GCP)"
Cohesion: 0.2
Nodes (1): MockService

### Community 16 - "Database & Persistence"
Cohesion: 0.31
Nodes (2): GormLogger, Config

### Community 17 - "Community 17"
Cohesion: 0.25
Nodes (2): GetEnvAsHexBytes(), GetRequiredEnv()

### Community 18 - "Observability & Logging"
Cohesion: 0.22
Nodes (3): Event, MultiTracker, Tracker

### Community 19 - "Observability & Logging"
Cohesion: 0.25
Nodes (3): ga4Event, ga4Payload, GA4Tracker

### Community 20 - "Internal Testing Utilities"
Cohesion: 0.25
Nodes (2): Config, Router

### Community 21 - "HTTP Middleware"
Cohesion: 0.25
Nodes (0): 

### Community 22 - "Internal Testing Utilities"
Cohesion: 0.25
Nodes (0): 

### Community 23 - "Observability & Logging"
Cohesion: 0.32
Nodes (3): mockTracker, TestGA4Tracker(), TestMultiTracker()

### Community 24 - "Community 24"
Cohesion: 0.29
Nodes (2): ResponseWriterWrapper, RetryRoundTripper

### Community 25 - "Community 25"
Cohesion: 0.43
Nodes (2): DecodeAndValidate(), Binder

### Community 26 - "Observability & Logging"
Cohesion: 0.33
Nodes (1): PostHogTracker

### Community 27 - "HTTP Middleware"
Cohesion: 0.4
Nodes (1): loggingResponseWriter

### Community 28 - "HTTP Middleware"
Cohesion: 0.6
Nodes (2): ETag(), responseBuffer

### Community 29 - "Community 29"
Cohesion: 0.4
Nodes (3): DiscordEmbed, DiscordEmbedField, DiscordWebhook

### Community 30 - "Community 30"
Cohesion: 0.4
Nodes (2): DefaultGitCloner, GitCloner

### Community 31 - "Database & Persistence"
Cohesion: 0.5
Nodes (1): TestModel

### Community 32 - "Database & Persistence"
Cohesion: 0.5
Nodes (2): BaseModel, ExternalConnection

### Community 33 - "Database & Persistence"
Cohesion: 0.5
Nodes (1): txKey

### Community 34 - "Internal Testing Utilities"
Cohesion: 0.5
Nodes (0): 

### Community 35 - "Internal Testing Utilities"
Cohesion: 0.5
Nodes (1): TestUser

### Community 36 - "Authentication & Security"
Cohesion: 0.5
Nodes (0): 

### Community 37 - "Authentication & Security"
Cohesion: 0.5
Nodes (0): 

### Community 38 - "Community 38"
Cohesion: 0.67
Nodes (2): PaginatedResponse, PaginationParams

### Community 39 - "HTTP Middleware"
Cohesion: 1.0
Nodes (1): CorsConfig

### Community 40 - "HTTP Middleware"
Cohesion: 0.67
Nodes (0): 

### Community 41 - "Authentication & Security"
Cohesion: 0.67
Nodes (0): 

### Community 42 - "Authentication & Security"
Cohesion: 0.67
Nodes (0): 

### Community 43 - "Community 43"
Cohesion: 0.67
Nodes (0): 

### Community 44 - "Internal Testing Utilities"
Cohesion: 0.67
Nodes (0): 

### Community 45 - "Community 45"
Cohesion: 1.0
Nodes (0): 

### Community 46 - "HTTP Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "Observability & Logging"
Cohesion: 1.0
Nodes (0): 

### Community 48 - "HTTP Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "HTTP Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "HTTP Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "HTTP Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "Internal Testing Utilities"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "Community 53"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "Internal Testing Utilities"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "Internal Testing Utilities"
Cohesion: 1.0
Nodes (0): 

### Community 56 - "Community 56"
Cohesion: 1.0
Nodes (0): 

### Community 57 - "Database & Persistence"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "Community 58"
Cohesion: 1.0
Nodes (0): 

### Community 59 - "Community 59"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **49 isolated node(s):** `ErrorResponse`, `PaginationParams`, `PaginatedResponse`, `Claims`, `CSRFConfig` (+44 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 45`** (2 nodes): `recovery.go`, `Recovery()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Middleware`** (2 nodes): `request_id_test.go`, `TestRequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Observability & Logging`** (2 nodes): `logging_test.go`, `TestScrubPayload()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Middleware`** (2 nodes): `trailing_slash_test.go`, `TestTrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Middleware`** (2 nodes): `csrf_test.go`, `TestCSRFMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Middleware`** (2 nodes): `request_id.go`, `RequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `HTTP Middleware`** (2 nodes): `trailing_slash.go`, `TrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Internal Testing Utilities`** (2 nodes): `goth_initializer_test.go`, `TestGothInitializer()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 53`** (2 nodes): `random.go`, `RandomString()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Internal Testing Utilities`** (2 nodes): `slug_test.go`, `TestSlugify()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Internal Testing Utilities`** (2 nodes): `server_test.go`, `TestServer_StartAndShutdown()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 56`** (2 nodes): `domain.go`, `IsValidSubdomain()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Database & Persistence`** (2 nodes): `db_test.go`, `TestSetupTestDatabase()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 58`** (1 nodes): `run.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 59`** (1 nodes): `build.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `ErrorResponse`, `PaginationParams`, `PaginatedResponse` to the rest of the system?**
  _49 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Middleware` be split into smaller, more focused modules?**
  _Cohesion score 0.09 - nodes in this community are weakly interconnected._
- **Should `Authentication & Security` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Cloud Services (GCP)` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `Internal Testing Utilities` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Database & Persistence` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._
- **Should `Cloud Services (GCP)` be split into smaller, more focused modules?**
  _Cohesion score 0.14 - nodes in this community are weakly interconnected._