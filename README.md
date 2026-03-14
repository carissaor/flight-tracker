# ✈️ flight-tracker

## Overview

I love travelling!! But between global conflicts, pandemics, economic shifts, and geopolitical tensions, flight prices have never been more unpredictable. **flight-tracker** pulls real-time pricing data across multiple routes and world-event signals to build up a dataset for forecasting where prices are headed.

🌐 **[Live Demo](https://flight-tracker-pink.vercel.app)**

---

## What It Does

- Fetches the cheapest available fares for 6 routes out of YVR via the Travelpayouts API
- Pulls world-event signals from Polymarket — real money prediction markets for geopolitical events (conflicts, pandemics, oil prices, travel bans)
- Calculates a **Global Chaos Score** (0–100) from weighted Polymarket probabilities to signal when prices are likely to spike
- Saves each price snapshot and event probability with a timestamp so history accumulates over time
- REST API serving price history, world-event data, and chaos score to a React frontend
- Collector runs every 6 hours on Railway, building up price history automatically

---

## Routes Tracked

| Origin | Destination |
|--------|-------------|
| YVR | LHR — London |
| YVR | NRT — Tokyo |
| YVR | SYD — Sydney |
| YVR | CDG — Paris |
| YVR | JFK — New York |
| YVR | HKG — Hong Kong |

---

## Global Chaos Score

A single 0–100 score computed from a weighted average of all Polymarket event probabilities. Higher-volume markets carry more weight since they represent more reliable crowd signals.

| Score | Level | Meaning |
|-------|-------|---------|
| 60+ | 😭 We are so cooked | Book ASAP and get a refundable ticket! |
| 40+ | 🌪️ It's giving chaos | Things are getting spicy...don't wait! |
| 20+ | 👀 Sus but manageable | Could be nothing. Could be everything. Check back soon! |
| 0+ | ✌️ Calm skies | Weirdly calm, book before that changes! |

---

## World Event Signals

Uses the [Polymarket](https://polymarket.com) Gamma API (no API key required) to fetch prediction market probabilities for events that historically impact flight prices:

- Wars and invasions
- Pandemic declarations
- Travel bans and airspace closures
- Ceasefires and peace deals
- Crude oil price movements
- Financial crises

Each market returns a 0–1 probability representing what traders think is the likelihood of that event occurring. These get stored alongside price snapshots for correlation analysis.

---

## Roadmap

- [x] Route price fetching (Travelpayouts)
- [x] PostgreSQL storage with timestamped snapshots
- [x] World-event signals (Polymarket)
- [x] Global Chaos Score
- [x] REST API (Go)
- [x] React frontend dashboard
- [x] Scheduled data collection (Railway cron)
- [x] Deploy API to Railway, frontend to Vercel
- [ ] Price prediction model (Python)
- [ ] Price alert notifications

---

## Getting Started

### Prerequisites

- Go 1.26+
- Node.js 18+ (for frontend)
- [Travelpayouts API token](https://travelpayouts.com) (free)

### Installation

```bash
git clone https://github.com/carissaor/flight-tracker.git
cd flight-tracker
go mod tidy
```

### Database Setup

```bash
psql postgres -c "CREATE DATABASE flight_tracker;"
psql "postgres://YOUR_USER@localhost:5432/flight_tracker" -f schema.sql
```

### Configuration

```bash
cp .env.example .env
```

```env
DATABASE_URL=postgres://YOUR_USER@localhost:5432/flight_tracker?sslmode=disable
TRAVELPAYOUTS_TOKEN=your_token_here
ORIGIN=YVR
```

### Run the Collector

```bash
go run ./cmd/collector
```

### Run the API Server

```bash
go run ./cmd/api
```

API runs on `http://localhost:8080`

### Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/routes` | All routes with latest and lowest price |
| GET | `/api/prices?route=YVR-LHR` | Price history for a specific route |
| GET | `/api/events` | Latest Polymarket world-event signals |
| GET | `/api/chaos` | Global chaos score and level |

### Example Responses

**GET /api/routes**
```json
[
  {
    "id": 1,
    "origin": "YVR",
    "destination": "LHR",
    "lowest_price": 787,
    "latest_price": 787,
    "depart_date": "2026-04-27"
  }
]
```

**GET /api/chaos**
```json
{
  "score": 28.7,
  "level": "MODERATE",
  "label": "sus but manageable 👀",
  "insight": "Could be nothing. Could be everything. Check back soon!",
  "market_count": 9
}
```

**GET /api/events**
```json
[
  {
    "question": "US x Iran ceasefire by March 31?",
    "probability": 0.18,
    "volume": 42381,
    "fetched_at": "2026-03-13T20:42:51Z"
  }
]
```

---

## Project Structure

```
flight-tracker/
├── cmd/
│   ├── api/
│   │   └── main.go          # REST API server
│   └── collector/
│       └── main.go          # Price + event collector (runs every 6h on Railway)
├── frontend/                # React dashboard (Vite)
├── schema.sql               # Database table definitions
├── Dockerfile               # API server Docker build
├── Dockerfile.collector     # Collector Docker build
├── go.mod
├── go.sum
└── README.md
```

---

## Deployment

| Service | Platform | Notes |
|---------|----------|-------|
| Frontend | Vercel | Auto-deploys on push |
| API server | Railway | Always on |
| Collector | Railway | Cron job every 6 hours |
| PostgreSQL | Railway | Persistent |

---

## Database Schema

```sql
routes   -- city pairs being tracked (e.g. YVR → LHR)
prices   -- price snapshots per route with timestamps
events   -- Polymarket world-event probabilities with timestamps
```

---

## Data Sources

| Source | Purpose | Auth |
|--------|---------|------|
| [Travelpayouts](https://travelpayouts.com) | Live flight prices by route | API token |
| [Polymarket](https://polymarket.com) | World-event prediction markets | None |

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Disclaimer

Flight price predictions are based on historical data and world-event signals. They are not financial advice. Always verify prices directly with airlines or booking platforms before making travel decisions.