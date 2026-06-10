# Fullstack Coding Test — 7Solutions
 
This repository contains all submissions for the Full-Stack Developer coding test, covering a React frontend, a Go REST/gRPC API, a TypeScript data transformation service, and a lottery system design proposal.
 
---
 
## Repository Structure
 
```
fullstack-coding-test/
├── frontend/           # React + TypeScript todo list (deployed on Vercel)
├── backend-go/         # Go REST API + gRPC with MongoDB and JWT
├── backend-ts/         # TypeScript service — group users by department (optional)
└── lottery-design/     # Lottery search system design proposal (no code)
```
 
---
 
## 1. Frontend — Auto Delete Todo List
 
A React + TypeScript todo list where items move into Fruit/Vegetable columns and auto-return after 5 seconds.
 
🔗 **Live Demo:** [https://fullstack-coding-test-tawny.vercel.app](https://fullstack-coding-test-tawny.vercel.app)
 
**Tech:** React 18, TypeScript, Vite
 
→ See [frontend/README.md](./frontend/README.md) for setup instructions.
 
---
 
## 2. Backend — User Management API (Go)
 
A RESTful API with MongoDB persistence, JWT authentication, input validation, and a gRPC server.
 
**Tech:** Go 1.26, MongoDB, JWT (HS256), Docker, gRPC
 
**Highlights:**
- Full CRUD for users — register, login, list, get, update, delete
- JWT middleware protecting all `/users` endpoints
- Background goroutine logging user count every 10 seconds
- Graceful shutdown on SIGINT/SIGTERM
- Unit tests with mocked repository (no MongoDB required)
- Docker + docker-compose — single command to run API + gRPC + MongoDB
- gRPC server on port `50051` with `CreateUser` and `GetUser`
**Quick start:**
```bash
cd backend-go
cp .env.example .env
docker-compose up --build
```
 
→ See [backend-go/README.md](./backend-go/README.md) for full setup and JWT guide.
 
---
 
## 3. Backend — User Department Grouping (TypeScript) *(Optional)*
 
A Node.js + TypeScript service that fetches users from [dummyjson.com](https://dummyjson.com/users) and groups them by department.
 
**Tech:** TypeScript, Node.js, Express, Jest
 
**Quick start:**
```bash
cd backend-ts
npm install
npm start
```
 
Open `http://localhost:3000/users/grouped`
 
→ See [backend-ts/README.md](./backend-ts/README.md) for details.
 
---
 
## 4. Lottery Search System Design *(Design proposal only — no code)*
 
A design document for searching 1 million lottery tickets with wildcard pattern matching (`****23`, `1****5`, `123***`) while ensuring no duplicate simultaneous results.
 
**Covers:**
- PostgreSQL schema with per-digit indexing for fast wildcard lookup
- `SELECT FOR UPDATE SKIP LOCKED` for atomic ticket reservation
- Reservation expiry to prevent tickets being held indefinitely
- Performance analysis and tradeoffs
→ See [lottery-design/DESIGN.md](./lottery-design/DESIGN.md)