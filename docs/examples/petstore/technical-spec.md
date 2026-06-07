# Technical Specification: PetStore API

This walkthrough demonstrates creating technical specifications from the Working Backwards artifacts.

## The Technical Flow

```
PRD + UXD (Requirements)
        ↓
TRD (Technical Architecture)
        ↓
TPD (Test Plan)
        ↓
IRD (Infrastructure)
        ↓
spec.md (Reconciled Specification)
```

## Step 1: User Experience Design (UXD)

The UXD defines how users interact with the system.

### Create the UXD

```bash
visionspec create uxd -p petstore-api
```

### Author the UXD

Edit `docs/specs/petstore-api/source/uxd.md`:

```markdown
# PetStore API - User Experience Design

## API Design Principles

### 1. RESTful Conventions
- Resource-oriented URLs (`/pets`, `/orders`)
- HTTP methods for actions (GET, POST, PUT, DELETE)
- Consistent response formats

### 2. Developer Experience
- Intuitive endpoint naming
- Comprehensive error messages
- Pagination for list endpoints
- Filtering and sorting support

### 3. Versioning
- URL versioning (`/v1/pets`)
- Backward compatibility within major versions
- Clear deprecation policy

## API Contracts

### Pet Resource

```yaml
Pet:
  type: object
  required:
    - name
    - status
  properties:
    id:
      type: integer
      format: int64
      readOnly: true
    name:
      type: string
      example: "Buddy"
    category:
      $ref: '#/components/schemas/Category'
    photoUrls:
      type: array
      items:
        type: string
    tags:
      type: array
      items:
        $ref: '#/components/schemas/Tag'
    status:
      type: string
      enum: [available, pending, sold]
```

### User Journeys

#### Journey 1: Store Owner Adds Pet

```
1. Owner logs into dashboard
2. Clicks "Add Pet"
3. Fills pet details form
4. Uploads photos
5. Clicks "Save"
6. Pet appears in inventory
```

**API Flow:**
```
POST /user/login → token
POST /pets → pet created
POST /pets/{id}/photos → photos uploaded
```

#### Journey 2: Customer Searches and Orders

```
1. Customer browses pets
2. Filters by criteria
3. Views pet details
4. Adds to cart
5. Completes checkout
6. Receives confirmation
```

**API Flow:**
```
GET /pets?status=available&category=dog → pet list
GET /pets/{id} → pet details
POST /store/order → order created
GET /store/order/{id} → order confirmation
```

#### Journey 3: Developer Integrates API

```
1. Developer signs up
2. Gets API key
3. Reads documentation
4. Tests in sandbox
5. Integrates into app
6. Goes live
```

**API Flow:**
```
POST /user → account created
GET /user/{id}/apikey → API key
GET /pets (with API key) → test integration
```

## Error Handling

### Standard Error Response

```json
{
  "code": 400,
  "type": "validation_error",
  "message": "Invalid pet status",
  "details": {
    "field": "status",
    "value": "invalid",
    "allowed": ["available", "pending", "sold"]
  }
}
```

### Error Codes

| Code | Type | Description |
|------|------|-------------|
| 400 | validation_error | Invalid request data |
| 401 | authentication_error | Invalid or missing auth |
| 403 | authorization_error | Insufficient permissions |
| 404 | not_found | Resource doesn't exist |
| 429 | rate_limit_exceeded | Too many requests |
| 500 | internal_error | Server error |

## Pagination

### Request

```
GET /pets?page=2&limit=20
```

### Response

```json
{
  "data": [...],
  "pagination": {
    "page": 2,
    "limit": 20,
    "total": 150,
    "pages": 8
  }
}
```

## Authentication

### API Key (Server-to-Server)

```
GET /pets
Authorization: ApiKey sk_live_xxxxx
```

### JWT (User Sessions)

```
POST /user/login
→ {"token": "eyJ..."}

GET /pets
Authorization: Bearer eyJ...
```
```

## Step 2: Technical Requirements (TRD)

The TRD synthesizes requirements into technical architecture.

### Synthesize the TRD

```bash
visionspec synthesize trd -p petstore-api
```

### Review Generated TRD

`docs/specs/petstore-api/technical/trd.md`:

```markdown
# PetStore API - Technical Requirements Document

## Architecture Overview

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Clients   │────→│  API Layer  │────→│  Database   │
│ (Web/Mobile)│     │    (Go)     │     │ (PostgreSQL)│
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │    Cache    │
                    │   (Redis)   │
                    └─────────────┘
```

## Technology Stack

### Backend
- **Language**: Go 1.22+
- **Framework**: Chi router
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Auth**: JWT + OAuth 2.0

### Infrastructure
- **Hosting**: AWS (ECS Fargate)
- **Database**: RDS PostgreSQL
- **Cache**: ElastiCache Redis
- **CDN**: CloudFront
- **DNS**: Route 53

### Development
- **API Spec**: OpenAPI 3.0
- **Testing**: Go testing + testify
- **CI/CD**: GitHub Actions
- **Monitoring**: Datadog

## Component Design

### API Layer

```go
// cmd/server/main.go
func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)
    cache := redis.Connect(cfg.RedisURL)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RateLimit(cfg.RateLimit))

    // Mount routes
    r.Mount("/v1/pets", pet.NewRouter(db, cache))
    r.Mount("/v1/store", store.NewRouter(db, cache))
    r.Mount("/v1/user", user.NewRouter(db))

    http.ListenAndServe(":8080", r)
}
```

### Pet Service

```go
// internal/pet/service.go
type Service struct {
    repo  Repository
    cache Cache
}

func (s *Service) Create(ctx context.Context, pet *Pet) error {
    if err := pet.Validate(); err != nil {
        return err
    }
    if err := s.repo.Create(ctx, pet); err != nil {
        return err
    }
    s.cache.Invalidate(ctx, "pets:list")
    return nil
}

func (s *Service) List(ctx context.Context, filter Filter) ([]Pet, error) {
    key := filter.CacheKey()
    if cached, ok := s.cache.Get(ctx, key); ok {
        return cached, nil
    }
    pets, err := s.repo.List(ctx, filter)
    if err != nil {
        return nil, err
    }
    s.cache.Set(ctx, key, pets, 5*time.Minute)
    return pets, nil
}
```

### Data Models

```go
// internal/pet/model.go
type Pet struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Category  *Category `json:"category,omitempty"`
    PhotoURLs []string  `json:"photoUrls" db:"photo_urls"`
    Tags      []Tag     `json:"tags,omitempty"`
    Status    Status    `json:"status" db:"status"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Status string

const (
    StatusAvailable Status = "available"
    StatusPending   Status = "pending"
    StatusSold      Status = "sold"
)
```

### Database Schema

```sql
-- migrations/001_initial.sql
CREATE TABLE pets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id BIGINT REFERENCES categories(id),
    photo_urls TEXT[],
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pets_status ON pets(status);
CREATE INDEX idx_pets_category ON pets(category_id);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE pet_tags (
    pet_id BIGINT REFERENCES pets(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (pet_id, tag_id)
);
```

## API Specification

### OpenAPI 3.0

```yaml
openapi: 3.0.3
info:
  title: PetStore API
  version: 1.0.0

paths:
  /v1/pets:
    get:
      summary: List pets
      parameters:
        - name: status
          in: query
          schema:
            type: string
            enum: [available, pending, sold]
        - name: page
          in: query
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Pet'
                  pagination:
                    $ref: '#/components/schemas/Pagination'
```

## Non-Functional Requirements

### Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| Response time (p50) | < 50ms | Datadog APM |
| Response time (p95) | < 200ms | Datadog APM |
| Response time (p99) | < 500ms | Datadog APM |
| Throughput | 1000 req/s | Load testing |

### Scalability

- Horizontal scaling via ECS auto-scaling
- Database read replicas for read-heavy workloads
- Redis cluster for cache scaling
- CDN for static assets

### Security

- TLS 1.3 for all connections
- JWT tokens with 1-hour expiry
- API key rotation support
- Rate limiting per client
- Input validation on all endpoints
- SQL injection prevention via parameterized queries

### Reliability

- Multi-AZ deployment
- Database automated backups (7-day retention)
- Health check endpoints
- Graceful shutdown handling
- Circuit breaker for external dependencies
```

## Step 3: Test Plan (TPD)

The TPD defines the testing strategy.

### Synthesize the TPD

```bash
visionspec synthesize tpd -p petstore-api
```

### Review Generated TPD

`docs/specs/petstore-api/technical/tpd.md`:

```markdown
# PetStore API - Test Plan Document

## Test Strategy

### Test Pyramid

```
        ┌─────────────┐
        │    E2E      │  10%
        ├─────────────┤
        │ Integration │  30%
        ├─────────────┤
        │    Unit     │  60%
        └─────────────┘
```

## Unit Tests

### Pet Service Tests

```go
// internal/pet/service_test.go
func TestService_Create(t *testing.T) {
    tests := []struct {
        name    string
        pet     *Pet
        wantErr bool
    }{
        {
            name: "valid pet",
            pet: &Pet{
                Name:   "Buddy",
                Status: StatusAvailable,
            },
            wantErr: false,
        },
        {
            name: "missing name",
            pet: &Pet{
                Status: StatusAvailable,
            },
            wantErr: true,
        },
        {
            name: "invalid status",
            pet: &Pet{
                Name:   "Buddy",
                Status: "invalid",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

### Coverage Requirements

| Package | Minimum Coverage |
|---------|-----------------|
| internal/pet | 80% |
| internal/store | 80% |
| internal/user | 80% |
| internal/auth | 90% |

## Integration Tests

### Database Tests

```go
// internal/pet/repository_test.go
func TestRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db := testutil.SetupTestDB(t)
    repo := NewRepository(db)

    t.Run("Create and Retrieve", func(t *testing.T) {
        pet := &Pet{Name: "Buddy", Status: StatusAvailable}
        err := repo.Create(context.Background(), pet)
        require.NoError(t, err)
        require.NotZero(t, pet.ID)

        retrieved, err := repo.Get(context.Background(), pet.ID)
        require.NoError(t, err)
        assert.Equal(t, pet.Name, retrieved.Name)
    })
}
```

### API Tests

```go
// api/pets_test.go
func TestPetsAPI(t *testing.T) {
    srv := testutil.NewTestServer(t)
    defer srv.Close()

    t.Run("POST /v1/pets", func(t *testing.T) {
        body := `{"name": "Buddy", "status": "available"}`
        resp := srv.Post("/v1/pets", body)
        assert.Equal(t, 201, resp.StatusCode)

        var pet Pet
        json.Unmarshal(resp.Body, &pet)
        assert.NotZero(t, pet.ID)
        assert.Equal(t, "Buddy", pet.Name)
    })
}
```

## End-to-End Tests

### User Journey: Add Pet

```gherkin
Feature: Add Pet
  As a store owner
  I want to add pets to inventory
  So that customers can purchase them

  Scenario: Successfully add a pet
    Given I am logged in as a store owner
    When I create a pet with name "Buddy" and status "available"
    Then the pet should be created
    And the pet should appear in the inventory list

  Scenario: Add pet with invalid data
    Given I am logged in as a store owner
    When I create a pet without a name
    Then I should receive a validation error
```

## Performance Tests

### Load Test Scenarios

```yaml
# k6/load-test.js
scenarios:
  list_pets:
    executor: constant-vus
    vus: 100
    duration: 5m
    exec: listPets

  create_pets:
    executor: ramping-vus
    startVUs: 0
    stages:
      - duration: 2m, target: 50
      - duration: 5m, target: 50
      - duration: 2m, target: 0
    exec: createPet

thresholds:
  http_req_duration:
    - p(95) < 200
    - p(99) < 500
  http_req_failed:
    - rate < 0.01
```
```

## Step 4: Infrastructure Requirements (IRD)

### Synthesize the IRD

```bash
visionspec synthesize ird -p petstore-api
```

## Step 5: Reconcile

Generate the unified specification:

```bash
# Evaluate technical specs
visionspec eval trd -p petstore-api
visionspec eval tpd -p petstore-api

# Approve
visionspec approve trd -p petstore-api --approver "tech@example.com"
visionspec approve tpd -p petstore-api --approver "qa@example.com"

# Reconcile all specs into spec.md
visionspec reconcile -p petstore-api
```

The reconciliation process:

1. Validates all required specs are approved
2. Detects conflicts between specs
3. Generates unified `spec.md`
4. Creates decision log for tradeoffs

## Next Steps

- [Execution](execution.md) - Export and implement with AI agents
