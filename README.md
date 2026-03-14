<p align="center">
  <h1 align="center">🎫 TicketX — Secure Secondary Ticket Marketplace</h1>
  <p align="center">
    A Proof-of-Concept escrow-based ticket exchange with <strong>Clean Architecture</strong> on both Backend (Go) and Frontend (React).
  </p>
  <p align="center">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
    <img src="https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React" />
    <img src="https://img.shields.io/badge/PostgreSQL-15-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white" alt="Tailwind" />
    <img src="https://img.shields.io/badge/WebSocket-Real--time-010101?style=for-the-badge&logo=socketdotio&logoColor=white" alt="WebSocket" />
  </p>
</p>

---

## Screenshots

> Add your own screenshots here! Run the app and capture the UI.

| Marketplace | Escrow Transaction | Admin Dashboard |
|:-----------:|:------------------:|:---------------:|
| ![Marketplace](docs/screenshots/marketplace.png) | ![Transaction](docs/screenshots/transaction.png) | ![Admin](docs/screenshots/admin.png) |

<details>
<summary>📸 How to add screenshots</summary>

1. Create a `docs/screenshots/` folder
2. Take screenshots of each page
3. Save them as `marketplace.png`, `transaction.png`, `admin.png`
4. They will automatically appear in this table

</details>

---

## Core Mechanism — Escrow Flow

```mermaid
sequenceDiagram
    actor Buyer
    actor Seller
    participant Platform as 🛡️ TicketX Platform
    participant Admin as 👮 Admin (Escrow)

    Seller->>Platform: List ticket for sale
    Buyer->>Platform: Click "Buy Now"
    Platform->>Admin: Create escrow transaction

    rect rgb(30, 41, 59)
        Note over Buyer,Admin: 🔒 Escrow Phase — Funds & Tickets held by Admin
        Seller->>Admin: Upload digital ticket
        Admin-->>Seller: ✅ Ticket received & held
        Buyer->>Admin: Upload payment proof
        Admin-->>Buyer: ✅ Payment verified
    end

    Admin->>Buyer: 🎫 Release ticket
    Admin->>Seller: 💰 Release funds
    Note over Buyer,Seller: ❌ Direct Buyer↔Seller chat blocked (prevents off-platform deals)
```

---

## Real-Time Chat — WebSocket Architecture

Each escrow transaction has a private chat room. Communication flows through WebSocket for real-time delivery, with REST fallback for message history.

```mermaid
sequenceDiagram
    actor User as Buyer / Seller
    participant FE as React Frontend
    participant WS as WebSocket Server
    participant DB as PostgreSQL

    Note over FE,DB: 📡 Connection Phase
    FE->>WS: GET /api/chat/ws/:tx_id?token=JWT
    WS->>WS: Validate JWT + check buyer/seller/admin
    WS-->>FE: 101 Switching Protocols
    FE->>FE: GET /api/chat/transactions/:tx_id/messages (REST)
    FE->>FE: Render chat history

    Note over FE,DB: 💬 Real-Time Messaging
    User->>FE: Type & send message
    FE->>WS: ws.send({content: "..."})
    WS->>DB: INSERT INTO message_models
    WS->>WS: Broadcast to all clients in room
    WS-->>FE: onmessage → render instantly

    Note over FE,DB: 🔐 Security
    Note right of WS: Only buyer, seller, and admin<br/>can join the room
```

**Key Design:**
- **Hub pattern** — one Hub manages all rooms (`tx:<id>`), rooms are created/destroyed lazily
- **History via REST** — `GET /api/chat/transactions/:tx_id/messages` loads past messages on page load
- **New messages via WS** — sent through WebSocket, persisted to DB, then broadcast to all room participants
- **Auth** — JWT token passed as query param on WS upgrade (browsers can't set headers on WebSocket)

---

## Tech Stack

```mermaid
graph LR
    subgraph Frontend
        React[React 19] --> Zustand
        React --> Axios
        React --> TW[Tailwind CSS 4]
        React --> Lucide[Lucide Icons]
    end

    subgraph Backend
        Go[Go / Gin] --> GORM
        GORM --> PG[(PostgreSQL)]
        Go --> JWT[JWT Auth]
        Go --> WS[WebSocket / Gorilla]
    end

    Axios -->|REST API| Go
    React -->|WebSocket| WS

    style Frontend fill:#0f172a,stroke:#3b82f6,color:#e2e8f0
    style Backend fill:#0f172a,stroke:#10b981,color:#e2e8f0
```

---

## Quick Start

### 1. Database (Docker)
```bash
# Option A: Quick start (hardcoded values)
docker run --name ticketx-postgres \
  -e POSTGRES_USER=ticketx_user \
  -e POSTGRES_PASSWORD=secret_password \
  -e POSTGRES_DB=ticketx_db \
  -p 5432:5432 -d postgres:alpine

# Option B: Load from .env (recommended — single source of truth)
cd backend && cp .env.example .env   # edit .env if needed
source .env
docker run --name ticketx-postgres \
  -e POSTGRES_USER=$DB_USER \
  -e POSTGRES_PASSWORD=$DB_PASSWORD \
  -e POSTGRES_DB=$DB_NAME \
  -p $DB_PORT:5432 -d postgres:alpine
```

### 2. Backend
```bash
cd backend
cp .env.example .env    # Never commit this!
go mod download
go run cmd/api/main.go  # → http://localhost:8080
```

### 3. Frontend
```bash
cd frontend
bun install             # or npm install
bun run dev             # → http://localhost:5173
```

---

## Architecture Overview

Both backend and frontend follow **Clean Architecture** — inner layers never depend on outer layers.

```mermaid
graph TB
    subgraph "Clean Architecture — Dependency Rule"
        direction TB
        E["🟢 Entity / Domain<br/><i>Pure business objects</i>"]
        U["🔵 Use Case / Features<br/><i>Application logic</i>"]
        R["🟠 Repository / Infrastructure<br/><i>Data access</i>"]
        I["🔴 Interface / Pages<br/><i>HTTP, UI, Framework</i>"]

        I -->|depends on| U
        I -->|depends on| R
        U -->|depends on| E
        R -->|depends on| E
    end

    style E fill:#065f46,stroke:#10b981,color:#d1fae5
    style U fill:#1e3a5f,stroke:#3b82f6,color:#dbeafe
    style R fill:#78350f,stroke:#f59e0b,color:#fef3c7
    style I fill:#7f1d1d,stroke:#ef4444,color:#fee2e2
```

---

## Backend Architecture (Go)

### Why `cmd/api/main.go`?

Follows the **Go community standard** project layout. Entry points live under `cmd/` to support multiple binaries (`api`, `worker`, `migrate`). `main.go` is purely a wiring point — DI only.

### 4-Layer Structure

```mermaid
graph LR
    Main["main.go<br/><i>DI Wiring</i>"] --> Handler
    Handler["Interface<br/><i>gin_server/</i>"] --> UseCase["Use Case<br/><i>use_case/</i>"]
    UseCase --> Entity["Entity<br/><i>entity/</i>"]
    Repo["Repository<br/><i>repository/</i>"] --> Entity

    Main --> Repo
    Main --> UseCase
    Handler -.->|injects| UseCase
    UseCase -.->|calls| Repo

    style Main fill:#1e293b,stroke:#64748b,color:#e2e8f0
    style Handler fill:#7f1d1d,stroke:#ef4444,color:#fee2e2
    style UseCase fill:#1e3a5f,stroke:#3b82f6,color:#dbeafe
    style Entity fill:#065f46,stroke:#10b981,color:#d1fae5
    style Repo fill:#78350f,stroke:#f59e0b,color:#fef3c7
```

```
backend/
├── cmd/api/main.go                  # Entry point — DI wiring only
├── internal/
│   ├── entity/                      # Layer 1: Domain models & interfaces
│   │   ├── user/                    #   Pure business objects, no framework imports
│   │   ├── ticket/
│   │   ├── transaction/
│   │   └── message/
│   ├── use_case/                    # Layer 2: Business logic
│   │   ├── auth.go                  #   Orchestrates entities & repositories
│   │   ├── ticket.go                #   Knows WHAT to do, not HOW
│   │   └── transaction.go
│   ├── repository/                  # Layer 3: Data access
│   │   ├── user_repository/         #   Implements interfaces from inner layers
│   │   ├── ticket_repository/       #   Contains GORM/DB-specific code
│   │   ├── transaction_repository/
│   │   └── message_repository/
│   └── interface/                   # Layer 4: HTTP handlers & routing
│       └── gin_server/              #   Maps HTTP ↔ Use Case
│           ├── handler_auth.go
│           ├── handler_ticket.go
│           └── middleware/
└── pkg/utils/                       # Shared utilities (JWT, bcrypt)
```

### Password Security (Defense in Depth)

```mermaid
graph LR
    PW["Plain Password"] -->|bcrypt cost 14| Hash["Hashed in DB"]
    Hash -->|json:'-' tag| API["API Response<br/><i>password excluded</i>"]
    Hash -->|FindByID omits column| Query["DB Query<br/><i>never loaded</i>"]

    style PW fill:#7f1d1d,stroke:#ef4444,color:#fee2e2
    style Hash fill:#78350f,stroke:#f59e0b,color:#fef3c7
    style API fill:#065f46,stroke:#10b981,color:#d1fae5
    style Query fill:#065f46,stroke:#10b981,color:#d1fae5
```

| Layer | Protection |
|-------|-----------|
| `bcrypt` cost 14 | ~1s per hash, brute force resistant |
| `json:"-"` tag | Password hash never serialized to JSON |
| `FindByID()` omits password | DB query excludes the password column |

---

## Frontend Architecture (React)

### 4-Layer Structure

```mermaid
graph LR
    Pages["Pages<br/><i>Login, Home, Admin...</i>"] --> Hooks["Features / Hooks<br/><i>useLogin, useTickets...</i>"]
    Hooks --> Infra["Infrastructure<br/><i>API repositories</i>"]
    Infra --> Domain["Domains<br/><i>Entities & Interfaces</i>"]

    DI["ServiceContainer<br/><i>DI</i>"] --> Infra
    Hooks --> DI

    style Pages fill:#7f1d1d,stroke:#ef4444,color:#fee2e2
    style Hooks fill:#1e3a5f,stroke:#3b82f6,color:#dbeafe
    style Infra fill:#78350f,stroke:#f59e0b,color:#fef3c7
    style Domain fill:#065f46,stroke:#10b981,color:#d1fae5
    style DI fill:#1e293b,stroke:#64748b,color:#e2e8f0
```

```
frontend/src/
├── domains/                         # Layer 1: Domain (innermost)
│   ├── auth/
│   │   ├── entities/User.ts         #   Pure TypeScript interfaces
│   │   └── repositories/AuthRepository.ts
│   ├── ticket/
│   │   ├── entities/Ticket.ts
│   │   └── repositories/TicketRepository.ts
│   ├── transaction/
│   │   ├── entities/Transaction.ts
│   │   └── repositories/TransactionRepository.ts
│   └── chat/
│       ├── entities/Message.ts
│       └── repositories/ChatRepository.ts
├── infrastructure/                  # Layer 2: Infrastructure
│   ├── api/
│   │   ├── apiClient.ts             #   Axios instance + JWT interceptor
│   │   ├── authRepository.ts        #   Implements IAuthRepository
│   │   ├── ticketRepository.ts
│   │   ├── transactionRepository.ts
│   │   └── chatRepository.ts        #   WebSocket implementation
│   └── services/
│       └── ServiceContainer.ts      #   DI container
├── features/                        # Layer 3: Features (hooks)
│   ├── auth/hooks/useLogin.ts
│   ├── ticket/hooks/useTickets.ts
│   ├── transaction/hooks/useTransaction.ts
│   └── chat/hooks/useChat.ts
├── pages/                           # Layer 4: Presentation (outermost)
├── components/                      # Shared UI components
├── store/                           # Zustand (auth state only)
└── App.tsx                          # Router + Layout
```

### Key Design Decisions
- **Domains** = pure TypeScript interfaces only. No React, no Axios, no framework code
- **Infrastructure** implements domain interfaces with real HTTP/WebSocket
- **ServiceContainer** = DI container. Swap implementations without touching features
- **Features** expose custom hooks. Pages stay thin (UI only)
- **Zustand** handles only auth state. Domain state lives in hooks

---

## API Endpoints

```mermaid
graph LR
    subgraph Public
        A1["POST /auth/register"]
        A2["POST /auth/login"]
    end

    subgraph Protected
        T1["GET /tickets"]
        T2["GET /tickets/my"]
        T3["POST /tickets"]
        T4["DELETE /tickets/:id"]
        TX1["POST /transactions"]
        TX2["GET /transactions/my"]
        TX3["POST /transactions/:id/status"]
        C1["GET /chat/transactions/:id/messages"]
        C2["POST /chat/transactions/:id/messages"]
        C3["🔌 WS /chat/ws/:id?token=JWT"]
    end

    subgraph Admin Only
        AD1["GET /admin/transactions"]
        AD2["POST /admin/transactions/:id/status"]
    end

    style Public fill:#065f46,stroke:#10b981,color:#d1fae5
    style Protected fill:#1e3a5f,stroke:#3b82f6,color:#dbeafe
    style Admin Only fill:#78350f,stroke:#f59e0b,color:#fef3c7
```

---

## Security Considerations

*   **Credentials**: Do not commit `DB_PASSWORD` or `JWT_SECRET`. They are managed via `.env` and ignored by `.gitignore`.
*   **Authentication**: All sensitive endpoints and WebSockets are guarded by JWT tokens.
*   **Passwords**: Hashed with `bcrypt` (cost 14) before storage. Never exposed in API responses (see Defense in Depth above).
*   **HTTPS**: Plaintext password transmission is safe over TLS — this is the industry standard used by all major platforms.

---

<p align="center">
  Built with ❤️ as a learning project for <strong>Clean Architecture</strong> + <strong>Escrow Systems</strong>
</p>
