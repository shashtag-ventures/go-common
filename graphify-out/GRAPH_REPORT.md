# Graph Report - /Users/shashwatguta/Desktop/go projects/go-common  (2026-04-18)

## Corpus Check
- Corpus is ~25,981 words - fits in a single context window. You may not need a graph.

## Summary
- 385 nodes · 374 edges · 61 communities detected
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## God Nodes (most connected - your core abstractions)
1. `mockFieldError` - 13 edges
2. `Repository[T]` - 10 edges
3. `gcpService` - 8 edges
4. `MockService` - 8 edges
5. `InMemFile` - 7 edges
6. `GormLogger` - 6 edges
7. `RealGothProvider` - 5 edges
8. `mockTracker` - 5 edges
9. `PostHogTracker` - 5 edges
10. `MultiTracker` - 5 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Communities

### Community 0 - "Module 0"
Cohesion: 0.07
Nodes (9): GetAuthenticatedUser(), GetAuthenticatedUserID(), MembershipStore, Repository, ErrorResponse, SendAutoErrorResponse(), SendErrorResponse(), Claims (+1 more)

### Community 1 - "Module 1"
Cohesion: 0.13
Nodes (3): mockFieldError, TestJsonError(), TestSendErrorResponse()

### Community 2 - "Module 2"
Cohesion: 0.12
Nodes (14): BitbucketConfig, CorsConfig, DatabaseConfig, GitHubConfig, GitLabConfig, GoogleConfig, JWTConfig, Level (+6 more)

### Community 3 - "Module 3"
Cohesion: 0.19
Nodes (6): EnrichLogger(), GetLoggerFromContext(), GetTraceID(), CtxKey, LogState, TraceIDExtractor

### Community 4 - "Module 4"
Cohesion: 0.18
Nodes (1): DBConfig

### Community 5 - "Module 5"
Cohesion: 0.17
Nodes (2): BufferFile, InMemFile

### Community 6 - "Module 6"
Cohesion: 0.17
Nodes (5): GothConfig, GothInitializer, GothProvider, OAuthProviderInitializer, RealGothProvider

### Community 7 - "Module 7"
Cohesion: 0.17
Nodes (2): gcpService, Service

### Community 8 - "Module 8"
Cohesion: 0.27
Nodes (1): Repository[T]

### Community 9 - "Module 9"
Cohesion: 0.22
Nodes (4): AddGracefulShutdownGoroutine(), DoneGracefulShutdownGoroutine(), Server, Server

### Community 10 - "Module 10"
Cohesion: 0.24
Nodes (4): getRealIP(), RequestLogger(), scrubPayload(), responseWriter

### Community 11 - "Module 11"
Cohesion: 0.2
Nodes (1): MockService

### Community 12 - "Module 12"
Cohesion: 0.31
Nodes (2): GormLogger, Config

### Community 13 - "Module 13"
Cohesion: 0.25
Nodes (2): GetEnvAsHexBytes(), GetRequiredEnv()

### Community 14 - "Module 14"
Cohesion: 0.22
Nodes (3): Event, MultiTracker, Tracker

### Community 15 - "Module 15"
Cohesion: 0.25
Nodes (3): ga4Event, ga4Payload, GA4Tracker

### Community 16 - "Module 16"
Cohesion: 0.25
Nodes (2): Config, Router

### Community 17 - "Module 17"
Cohesion: 0.25
Nodes (0): 

### Community 18 - "Module 18"
Cohesion: 0.25
Nodes (0): 

### Community 19 - "Module 19"
Cohesion: 0.32
Nodes (1): CSRFConfig

### Community 20 - "Module 20"
Cohesion: 0.32
Nodes (3): mockTracker, TestGA4Tracker(), TestMultiTracker()

### Community 21 - "Module 21"
Cohesion: 0.29
Nodes (2): ResponseWriterWrapper, RetryRoundTripper

### Community 22 - "Module 22"
Cohesion: 0.43
Nodes (2): DecodeAndValidate(), Binder

### Community 23 - "Module 23"
Cohesion: 0.33
Nodes (1): PostHogTracker

### Community 24 - "Module 24"
Cohesion: 0.38
Nodes (3): RunSync(), SafeGo(), Task

### Community 25 - "Module 25"
Cohesion: 0.4
Nodes (4): RateLimitConfig, RateLimitStore, RateLimitMiddleware(), RateLimitMiddlewareWithContext()

### Community 26 - "Module 26"
Cohesion: 0.4
Nodes (1): loggingResponseWriter

### Community 27 - "Module 27"
Cohesion: 0.6
Nodes (2): ETag(), responseBuffer

### Community 28 - "Module 28"
Cohesion: 0.4
Nodes (3): DiscordEmbed, DiscordEmbedField, DiscordWebhook

### Community 29 - "Module 29"
Cohesion: 0.4
Nodes (0): 

### Community 30 - "Module 30"
Cohesion: 0.4
Nodes (2): DefaultGitCloner, GitCloner

### Community 31 - "Module 31"
Cohesion: 0.67
Nodes (2): mockRoundTripper, TestRetryRoundTripper()

### Community 32 - "Module 32"
Cohesion: 0.5
Nodes (1): txKey

### Community 33 - "Module 33"
Cohesion: 0.5
Nodes (1): TestUser

### Community 34 - "Module 34"
Cohesion: 0.5
Nodes (0): 

### Community 35 - "Module 35"
Cohesion: 0.5
Nodes (0): 

### Community 36 - "Module 36"
Cohesion: 0.67
Nodes (2): PaginatedResponse, PaginationParams

### Community 37 - "Module 37"
Cohesion: 1.0
Nodes (1): CorsConfig

### Community 38 - "Module 38"
Cohesion: 0.67
Nodes (0): 

### Community 39 - "Module 39"
Cohesion: 0.67
Nodes (0): 

### Community 40 - "Module 40"
Cohesion: 0.67
Nodes (0): 

### Community 41 - "Module 41"
Cohesion: 0.67
Nodes (1): TestModel

### Community 42 - "Module 42"
Cohesion: 0.67
Nodes (1): BaseModel

### Community 43 - "Module 43"
Cohesion: 0.67
Nodes (0): 

### Community 44 - "Module 44"
Cohesion: 0.67
Nodes (0): 

### Community 45 - "Module 45"
Cohesion: 1.0
Nodes (0): 

### Community 46 - "Module 46"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "Module 47"
Cohesion: 1.0
Nodes (0): 

### Community 48 - "Module 48"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "Module 49"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "Module 50"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "Module 51"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "Module 52"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "Module 53"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "Module 54"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "Module 55"
Cohesion: 1.0
Nodes (0): 

### Community 56 - "Module 56"
Cohesion: 1.0
Nodes (0): 

### Community 57 - "Module 57"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "Module 58"
Cohesion: 1.0
Nodes (0): 

### Community 59 - "Module 59"
Cohesion: 1.0
Nodes (0): 

### Community 60 - "Module 60"
Cohesion: 1.0
Nodes (0): 

## Knowledge Gaps
- **42 isolated node(s):** `ErrorResponse`, `PaginationParams`, `PaginatedResponse`, `Claims`, `CSRFConfig` (+37 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Module 45`** (2 nodes): `recovery.go`, `Recovery()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 46`** (2 nodes): `request_id_test.go`, `TestRequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 47`** (2 nodes): `logging_test.go`, `TestScrubPayload()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 48`** (2 nodes): `trailing_slash_test.go`, `TestTrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 49`** (2 nodes): `csrf_test.go`, `TestCSRFMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 50`** (2 nodes): `request_id.go`, `RequestIDMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 51`** (2 nodes): `trailing_slash.go`, `TrailingSlashMiddleware()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 52`** (2 nodes): `goth_initializer_test.go`, `TestGothInitializer()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 53`** (2 nodes): `random.go`, `RandomString()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 54`** (2 nodes): `slug_test.go`, `TestSlugify()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 55`** (2 nodes): `server_test.go`, `TestServer_StartAndShutdown()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 56`** (2 nodes): `service_test.go`, `TestGCPService_CoveragePlaceholder()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 57`** (2 nodes): `domain.go`, `IsValidSubdomain()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 58`** (2 nodes): `db_test.go`, `TestSetupTestDatabase()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 59`** (1 nodes): `run.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Module 60`** (1 nodes): `build.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `ErrorResponse`, `PaginationParams`, `PaginatedResponse` to the rest of the system?**
  _42 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Module 0` be split into smaller, more focused modules?**
  _Cohesion score 0.07 - nodes in this community are weakly interconnected._
- **Should `Module 1` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Module 2` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._