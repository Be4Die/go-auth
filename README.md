# go-auth

`go-auth` is a production-grade authentication and authorization service written in Go.
It is designed as a standalone **Identity Provider (IdP)** for microservice-based systems.

The project focuses on **security, scalability, and clean architecture**, rather than being a simple CRUD example.

---

## ✨ Features

### Authentication
- User registration & login
- Secure password hashing (Argon2id / bcrypt)
- Email verification
- Access & refresh tokens (JWT)
- Refresh token rotation & revocation
- Logout (session invalidation)

### Authorization
- Role-Based Access Control (RBAC)
- Fine-grained permissions (`resource.action`)
- Multi-tenant support (organizations / workspaces)

### OAuth2
- Authorization Code Flow
- Client Credentials Flow
- Token introspection & revocation
- Public & confidential clients

### Security
- Short-lived stateless JWT access tokens
- Stateful refresh tokens stored as hashes
- Brute-force protection & rate limiting
- Audit log for security events
- Constant-time comparisons
- CSRF protection (for browser flows)

### Multi-Factor Authentication (MFA)
- TOTP (RFC 6238)
- QR code provisioning
- Backup recovery codes

### Observability
- Structured logging
- Prometheus metrics
- OpenTelemetry tracing

---

## 🏗 Architecture

The service follows **Clean Architecture / Hexagonal Architecture** principles.

```

cmd/
└── auth-service/
└── main.go

internal/
├── domain/          // core entities & interfaces
├── app/             // use cases (business logic)
├── infrastructure/ // database, cache, email, external services
├── transport/
│    ├── http/       // REST API
│    └── grpc/       // internal communication
├── security/        // JWT, password hashing, MFA
└── config/

````

### Why this approach?
- Clear separation of concerns
- Business logic independent from frameworks
- Easy to test and extend
- Suitable for real-world production systems

---

## 🔐 Token Strategy

### Access Token
- JWT
- Short-lived (5–15 minutes)
- Stateless (not stored in DB)

### Refresh Token
- Long-lived
- Stored **hashed** in the database
- Rotated on every refresh
- Can be revoked at any time

This approach balances **performance, security, and scalability**.

---

## 👥 Domain Model (simplified)

- **User**
- **Tenant (Organization)**
- **Role**
- **Permission**
- **Session**
- **OAuth Client**
- **Audit Event**

The model is designed to support **multi-tenant SaaS applications**.

---

## 🌐 API

### REST API
- JSON-based
- OpenAPI 3.0 specification

### gRPC API
- Intended for internal service-to-service communication
- Authentication via interceptors

---

## 🛠 Tech Stack

- Go
- PostgreSQL
- Redis
- JWT
- OAuth2
- Docker
- Prometheus
- OpenTelemetry

---

## 🚀 Running Locally

```bash
docker-compose up
````

Environment-based configuration is used.
Database migrations are applied automatically on startup.

---

## 📚 Project Goals

This project was created to:

* Demonstrate **real-world backend engineering skills**
* Show understanding of **authentication & authorization internals**
* Serve as a reusable authentication service for other projects
* Act as a strong portfolio project for a **Middle Go Backend Developer**

---

## ⚠️ Disclaimer

This project is for educational and portfolio purposes.
It is designed with production principles in mind but should be security-reviewed
before being used in real production environments.
