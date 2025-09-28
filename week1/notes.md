# Week 1 Session Notes - HLD & LLd for API gateway

## Kickoff & Context

- This 5-week challenge = building a scalable API Gateway in Go.
- Why? To learn real-world system design & coding practices.
- End Goal → Gateway that handles 100–1000 RPS, with rate limiting, service discovery, and API validation.

## High level Design (HLD)

### Problem Context

- Old monoliths bundled auth, rate limiting, monitoring → scaling was hard.
- Microservices solved some issues but introduced:
  - Clients juggling multiple URLs.
  - Duplication of cross-cutting concerns.

### Why API Gateway?

- Single entry point for clients.
- Centralizes common logic (auth, logging, monitoring, rate limiting).
- Simplifies scaling and service management.

### Core Responsibilities

- Request routing.
- Logging & monitoring.
- Future: auth, security hooks, rate limiting.

### Scalability Concerns

- Stateless design → horizontal scaling.
- Load balancing in front of gateways.
- Caching: DNS/service discovery, request-level caching.

### Basic Flow

Client → API Gateway → Backend Services

![Diagram](hld.drawio.svg)

## Low-Level Design (LLD)

### Modules

- Logger & monitoring → log and add monitoring for all incoming/outgoing requests.
- Security → Blocking blacklisted IPs and basic security to validate request data.
- Rate Limiting → IP based rate limiting (configurable from the config at startup)
- Authorization → vailidate requests, if API that need Authorization.
- Rate limiting → User based rate limiting (configurable from the config at startup)
  Request Handler → parse requests.
- Router → decide backend service.
- Forwarder → perform backend HTTP calls using Service discovery and Load balancing.
- Response Handler → format & return response.

Note: currently we will not focus on API validation, but just have placeholder for that and enhance that in future

![Diagram](basic_lld.drawio.svg)

### Project Structure (Go)

```bash
/cmd/gateway       → entrypoint
/pkg/router        → routing logic
/pkg/logger        → structured logs
/pkg/middleware    → auth, rate limiting (future)
```

## Key Considerations

- Concurrency → goroutines for each request.
- Avoid shared state → channels/sync carefully.
- Error handling → timeouts, retries, fallbacks.
