# Culinary Digital Twin | Blinkit-Style Mobile Redesign

This repository contains the complete implementation of the **Culinary Digital Twin** under the **PromptWars Mumbai Domination Framework**. 

The application has been transformed from a cyberpunk HUD into an **ultra-clean, mobile-first Blinkit-style consumer application** designed for rapid single-handed tracking, featuring an interactive **3D Anime Assistant Girl**, **Gemini Live AI Voice Chat**, **Google OAuth**, and **Google Calendar synchronizers**.

---

## 🏗️ System Architecture

The project is structured into three highly isolated, production-grade components:

### 1. Frontend (`/frontend-nextjs`)
Built with **Next.js 14 (TypeScript)** and styled with **TailwindCSS**:
* **Blinkit Interface**: Clean top header with location tracking, high-contrast 8-minute ETA badge, a search bar, and grid-based grocery items with increment/decrement counters.
* **Procedural 3D Anime Assistant**: Rendered using standard Three.js (React Three Fiber) primitives to guarantee zero CORS or gltf loading crashes. Includes breathing, cursor-tracking, blinking, and speech-amplitude mouth oscillations.
* **Gemini Voice API**: Converts user speech to text (Speech-to-Text), queries Google's Gemini-1.5-flash endpoint (injecting real-time metabolic context), and speaks the answer out loud (Text-to-Speech) while animating the mouth in perfect sync.
* **Google OAuth & Calendar Sync**: Implements interactive profile sign-in popup flows and direct calendar slot schedulers.
* **Dynamic API URL Resolution**: Solves the "localhost pinning" production issue. Fetch endpoints dynamically prioritize `process.env.NEXT_PUBLIC_API_URL` and fall back to the live GCP Cloud Run URL in production.

### 2. Backend (`/backend-go`)
A high-performance **Go Gin-Gonic** microservice:
* **Auto-Migrations**: Automatically verifies connection pools and runs startup `CREATE TABLE IF NOT EXISTS` commands for `metabolic_logs`, `calendar_slots`, and `checkout_transactions` tables.
* **Tuned Connection Pooling**: Configured with `db.SetMaxOpenConns(25)`, `db.SetMaxIdleConns(10)`, and `db.SetConnMaxLifetime(15 * time.Minute)` to prevent descriptor leakages on high serverless container scale-ups.
* **Hardened Security**: Custom middlewares inject explicit CORS access controls and security headers (CSP, HSTS, X-Content-Type, X-Frame DENY).

### 3. Infrastructure (`/infrastructure`)
Infrastructure-as-Code utilizing **Terraform**:
* **Secured PostgreSQL Cloud SQL**: Provisions PostgreSQL 15 instances with enforced SSL requirements.
* **Cloud KMS Keyring & CMEK**: Provisions cryptographic key rings and Biometric keys to enforce Customer-Managed Encryption Keys (CMEK) at rest for database storage.

---

## 🧪 Comprehensive Automated Testing

To ensure total code correctness and secure compilation pathways, full testing blocks have been integrated directly within the folders:

### 1. Go Unit Test Suite (`/backend-go/main_test.go`)
Uses `net/http/httptest` and `github.com/DATA-DOG/go-sqlmock` to test 100% of the target Gin endpoints under both standalone fallback modes (`db == nil`) and active PostgreSQL connection states:
* `GET /api/v1/health/metabolic`
* `POST /api/v1/health/metabolic`
* `POST /api/v1/checkout`
* `POST /api/v1/calendar`
* `GET /api/v1/calendar`
* Custom security and CORS header injections.

### 2. Frontend Bindings Verifier (`/frontend-nextjs/test-bindings.js`)
An automated integration checker in Node.js that runs parsing audits, validates contract payload structures, and simulates API calls under active or offline database states to verify that the frontend has 100% failover resilience (zero UI crashes).

---

## 🚀 Local Quickstart Guide

Follow these steps to run the entire stack locally:

### 1. Provision the Database
Ensure a local PostgreSQL instance is running on port `5432` with a database named `digital_twin`, or start one via Docker:
```bash
docker run --name twin-postgres -e POSTGRES_DB=digital_twin -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

### 2. Run Go Backend
Navigate to `/backend-go` and execute:
```bash
go run main.go
```
*Note: The backend will auto-create all database tables and log: `[INFO] Startup Migration Steps succeeded.`*

### 3. Run Next.js Frontend
Navigate to `/frontend-nextjs` and execute:
```bash
npm install
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) to interact with the Culinary Digital Twin!

---

## 🔒 Security & Cryptographic Compliance
* **Storage CMEK**: Enforced at rest utilizing global Cloud KMS keys.
* **Transport Encryption**: SSL required for SQL handshakes and standard HTTPS routing.
* **Authentication Boundary**: OAuth processed via browser-popup interfaces, isolating private token states from database boundaries.
* **Application KMS Sealing**: Checkout triggers RSA simulated KMS biometric signing key handshakes, saving cryptographically sealed signature envelopes directly to the database ledger.
