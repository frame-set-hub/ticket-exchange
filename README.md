# TicketX 🎫 - Secure Secondary Ticket Marketplace

TicketX is a Proof-of-Concept (POC) for a secondary market ticket exchange platform. It focuses heavily on security and an **Admin-centric Escrow system** to prevent fraud (e.g., selling fake tickets or not paying).

## Core Mechanism (Escrow)
When a buyer decides to purchase a ticket, a transaction is created:
1. **Seller <-> Admin**: Seller uploads the digital ticket to the Admin. Admin "Holds" it.
2. **Buyer <-> Admin**: Buyer uploads proof of payment to the Admin. Admin verifies it.
3. **Completion**: Admin releases the ticket to the buyer and the funds to the seller. 
*(Direct chatting between Buyer and Seller is restricted to prevent off-platform deals)*

## Tech Stack
*   **Frontend**: React (Vite), Tailwind CSS, Zustand, Axios, Lucide React.
*   **Backend**: Golang (Gin), GORM (PostgreSQL).
*   **Real-time**: WebSockets (Gorilla) for secure Admin chats.

## Setup & Run Instructions

### 1. Database Setup (Docker)
We use PostgreSQL. You can start it via Docker:
```bash
docker run --name ticketx-postgres -e POSTGRES_USER=ticketx_user -e POSTGRES_PASSWORD=secret_password -e POSTGRES_DB=ticketx_db -p 5432:5432 -d postgres:alpine
```

### 2. Backend Setup Go
1. Navigate to the backend directory: `cd backend`
2. Create your `.env` file (never commit this!): `cp .env.example .env`
3. Download dependencies: `go mod download`
4. Run the server: `go run cmd/api/main.go` (Server starts on port 8080)

### 3. Frontend Setup
1. Navigate to the frontend directory: `cd frontend`
2. Install dependencies (we use Bun, but npm works too): `bun install`
3. Start the dev server: `bun run dev`

## Security Considerations 🔒
*   **Credentials**: Do not commit actual `DB_PASSWORD` or `JWT_SECRET` to source control. They are ignored via `.gitignore` and mapped using `.env`.
*   **Authentication**: All sensitive endpoints and WebSockets are guarded by JWT tokens.
*   **Passwords**: User passwords are encrypted using `bcrypt` before database entry.