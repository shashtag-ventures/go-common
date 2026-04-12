# Graph Report - .  (2026-04-12)

## Corpus Check
- Corpus is ~25,981 words - fits in a single context window. You may not need a graph.

## Summary
- 435 nodes · 453 edges · 61 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## God Nodes (most connected - your core abstractions)
1. `mockFieldError` - 13 edges
2. `IntegrationService` - 12 edges
3. `GitHubClient` - 11 edges
4. `Repository[T]` - 10 edges
5. `gcpService` - 8 edges
6. `MockService` - 8 edges
7. `InMemFile` - 7 edges
8. `integrationRepository` - 7 edges
9. `GormLogger` - 6 edges
10. `RealGothProvider` - 5 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Communities

### Community 0 - "gcpService"
Cohesion: 0.07
Nodes (8): gcpService, Service, IntegrationStorage, ContentItem, IntegrationClient, Namespace, Repository, TokenRefreshResponse

### Community 1 - "errors.go"
Cohesion: 0.1
Nodes (7): GetAuthenticatedUser(), GetAuthenticatedUserID(), MembershipStore, GormLogger, mockRoundTripper, Config, TestRetryRoundTripper()

### Community 2 - "jsonResponse.go"
Cohesion: 0.09
Nodes (9): ErrorResponse, SendAutoErrorResponse(), SendErrorResponse(), Claims, AuthenticatedUser, RateLimitConfig, RateLimitStore, RateLimitMiddleware() (+1 more)

### Community 3 - "GitHubClient"
Cohesion: 0.22
Nodes (5): GitHubClient, githubRepo, githubRepoWrapper, extractNextPageURL(), isInstallationContext()

### Community 4 - "mockFieldError"
Cohesion: 0.13
Nodes (3): mockFieldError, TestJsonError(), TestSendErrorResponse()

### Community 5 - "structs.go"
Cohesion: 0.12
Nodes (14): BitbucketConfig, CorsConfig, DatabaseConfig, GitHubConfig, GitLabConfig, GoogleConfig, JWTConfig, Level (+6 more)

### Community 6 - "postgres.go"
Cohesion: 0.18
Nodes (1): DBConfig

### Community 7 - "InMemFile"
Cohesion: 0.17
Nodes (2): BufferFile, InMemFile

### Community 8 - "goth_initializer.go"
Cohesion: 0.17
Nodes (5): GothConfig, GothInitializer, GothProvider, OAuthProviderInitializer, RealGothProvider

### Community 9 - "IntegrationService"
Cohesion: 0.3
Nodes (1): IntegrationService

### Community 10 - "Repository[T]"
Cohesion: 0.27
Nodes (1): Repository[T]

### Community 11 - "integrationRepository"
Cohesion: 0.18
Nodes (2): Repository, integrationRepository

### Community 12 - "server.go"
Cohesion: 0.22
Nodes (4): AddGracefulShutdownGoroutine(), DoneGracefulShutdownGoroutine(), Server, Server

### Community 13 - "logging.go"
Cohesion: 0.24
Nodes (4): getRealIP(), RequestLogger(), scrubPayload(), responseWriter

### Community 14 - "context_keys.go"
Cohesion: 0.24
Nodes (4): EnrichLogger(), GetLoggerFromContext(), CtxKey, LogState

### Community 15 - "MockService"
Cohesion: 0.2
Nodes (1): MockService

### Community 16 - "env.go"
Cohesion: 0.25
Nodes (2): GetEnvAsHexBytes(), GetRequiredEnv()

### Community 17 - "MultiTracker"
Cohesion: 0.22
Nodes (3): Event, MultiTracker, Tracker

### Community 18 - "GA4Tracker"
Cohesion: 0.25
Nodes (3): ga4Event, ga4Payload, GA4Tracker

### Community 19 - "router.go"
Cohesion: 0.25
Nodes (2): Config, Router

### Community 20 - "url.go"
Cohesion: 0.32
Nodes (1): CSRFConfig

### Community 21 - "middleware_test.go"
Cohesion: 0.25
Nodes (0): 

### Community 22 - "env_test.go"
Cohesion: 0.25
Nodes (0): 

### Community 23 - "mockTracker"
Cohesion: 0.32
Nodes (3): mockTracker, TestGA4Tracker(), TestMultiTracker()

### Community 24 - "ResponseWriterWrapper"
Cohesion: 0.29
Nodes (2): ResponseWriterWrapper, RetryRoundTripper

### Community 25 - "Binder"
Cohesion: 0.43
Nodes (2): DecodeAndValidate(), Binder

### Community 26 - "PostHogTracker"
Cohesion: 0.33
Nodes (1): PostHogTracker

### Community 27 - "worker.go"
Cohesion: 0.38
Nodes (3): RunSync(), SafeGo(), Task

### Community 28 - "loggingResponseWriter"
Cohesion: 0.4
Nodes (1): loggingResponseWriter

### Community 29 - "ETag()"
Cohesion: 0.6
Nodes (2): ETag(), responseBuffer

### Community 30 - "discord.go"
Cohesion: 0.4
Nodes (3): DiscordEmbed, DiscordEmbedField, DiscordWebhook

### Community 31 - "cloner.go"
Cohesion: 0.4
Nodes (2): DefaultGitCloner, GitCloner

### Community 32 - "repository_test.go"
Cohesion: 0.5
Nodes (1): TestModel

### Community 33 - "tx.go"
Cohesion: 0.5
Nodes (1): txKey

### Community 34 - "BaseModel"
Cohesion: 0.5
Nodes (2): BaseModel, ExternalConnection

### Community 35 - "github_test.go"
Cohesion: 0.5
Nodes (0): 

### Community 36 - "binder_test.go"
Cohesion: 0.5
Nodes (1): TestUser

### Community 37 - "helpers_test.go"
Cohesion: 0.5
Nodes (0): 

### Community 38 - "helpers.go"
Cohesion: 0.5
Nodes (0): 

### Community 39 - "pagination.go"
Cohesion: 0.67
Nodes (2): PaginatedResponse, PaginationParams

### Community 40 - "cors.go"
Cohesion: 1.0
Nodes (1): CorsConfig

### Community 41 - "middleware_extra_test.go"
Cohesion: 0.67
Nodes (0): 

### Community 42 - "cookies.go"
Cohesion: 0.67
Nodes (0): 

### Community 43 - "authz_test.go"
Cohesion: 0.67
Nodes (0): 

### Community 44 - "slug.go"
Cohesion: 0.67
Nodes (0): 

### Community 45 - "url_test.go"
Cohesion: 0.67
Nodes (0): 

### Community 46 - "recovery.go"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "request_id_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 48 - "logging_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "trailing_slash_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "csrf_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "request_id.go"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "trailing_slash.go"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "goth_initializer_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "random.go"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "slug_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 56 - "server_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 57 - "domain.go"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "db_test.go"
Cohesion: 1.0
Nodes (0): 

### Community 59 - "run.go"
Cohesion: 1.0
Nodes (0): 

### Community 60 - "build.go"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **50 isolated node(s):** `ErrorResponse`, `PaginationParams`, `PaginatedResponse`, `Claims`, `CSRFConfig` (+45 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `recovery.go`** (2 nodes): `recovery.go`, `Recovery()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `request_id_test.go`** (2 nodes): `request_id_test.go`, `TestRequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `logging_test.go`** (2 nodes): `logging_test.go`, `TestScrubPayload()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `trailing_slash_test.go`** (2 nodes): `trailing_slash_test.go`, `TestTrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `csrf_test.go`** (2 nodes): `csrf_test.go`, `TestCSRFMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `request_id.go`** (2 nodes): `request_id.go`, `RequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `trailing_slash.go`** (2 nodes): `trailing_slash.go`, `TrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `goth_initializer_test.go`** (2 nodes): `goth_initializer_test.go`, `TestGothInitializer()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `random.go`** (2 nodes): `random.go`, `RandomString()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `slug_test.go`** (2 nodes): `slug_test.go`, `TestSlugify()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `server_test.go`** (2 nodes): `server_test.go`, `TestServer_StartAndShutdown()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `domain.go`** (2 nodes): `domain.go`, `IsValidSubdomain()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `db_test.go`** (2 nodes): `db_test.go`, `TestSetupTestDatabase()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `run.go`** (1 nodes): `run.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `build.go`** (1 nodes): `build.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `IntegrationService` connect `IntegrationService` to `gcpService`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `ErrorResponse`, `PaginationParams`, `PaginatedResponse` to the rest of the system?**
  _50 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `gcpService` be split into smaller, more focused modules?**
  _Cohesion score 0.07 - nodes in this community are weakly interconnected._
- **Should `errors.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `jsonResponse.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09 - nodes in this community are weakly interconnected._
- **Should `mockFieldError` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `structs.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._