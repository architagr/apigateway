# apigateway
Scalable API Gateway in Go

## 🎯 Epic: Build a Scalable API Gateway in Go
The goal of this epic is to design and implement an API Gateway in Golang that can scale under load, enforce different rate limiting strategies, perform service discovery, validate APIs, and route traffic efficiently.

## Story 1: System Design Blueprint (Weekend 1)
**Objective:** Create a scalable system design for the API Gateway.
- Define the **architecture** capable of handling 100–1000 RPS.
- Document **components:** load balancer, gateway, backend services.
- Decide how requests flow: Client → Gateway → API Server.
- Define high-level **scalability strategies** (horizontal scaling, stateless gateway, caching).
- Output: A system design diagram + reasoning.

## Story 2: Basic API Gateway Skeleton (Weekend 2)
**Objective:** Implement a barebones API gateway in Go.
- Accept incoming HTTP requests.
- Forward requests to a configured backend service.
- Return the response back to the client.
- Add logging for each request (to prepare for rate limiting).
- Output: Working Go code for a gateway → backend flow.

## Story 3: Configurable Rate Limiter (Weekend 3)
**Objective:** Extend the gateway with a rate limiter.
- Support at least one rate limiting strategy (e.g., Token Bucket).
- Rate limit should be **configurable at startup** (limit & burst size).
- Apply limits at the **gateway level** (all requests).
- Output: Gateway that drops/throttles requests beyond allowed RPS.

## Story 4: Advanced Rate Limiting (Weekend 4)
**Objective:** Implement fine-grained & flexible rate limiting.
- Support **all 4 algorithms:** Token Bucket, Leaky Bucket, Fixed Window, Sliding Window.
- Rate limit by **IP** and by **User ID** (present in the JWT token).
- Configurable at startup which algorithm to use.
- Output: Gateway with pluggable rate limiter system.
  
## Story 5: Service Discovery & API Validation (Weekend 5–6)
**Objective:** Make the gateway “smart.”
- Implement **service registration:** Gateway should know which services are available.
- Add **health checks** for services (periodic ping).
- Reject requests for **invalid/unregistered APIs** before hitting backend.
- Route requests only to healthy instances.
- Output: Gateway with service discovery + validation layer.
<hr/>

👉 By the end of this series, the community will have built a mini Kong/NGINX-style API Gateway in Go 💪
