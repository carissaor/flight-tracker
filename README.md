# ✈️ flight-tracker

> Flight price monitoring and prediction — because the world situation shouldn't catch your wallet off guard.

---

## Overview

We all love travelling. But between global conflicts, pandemics, economic shifts, and geopolitical tensions, flight prices have never been more unpredictable. **flight-tracker** pulls real-time pricing data across multiple routes and builds up a historical record to eventually forecast where prices are headed.

---

## What It Does Right Now

- Connects to a local PostgreSQL database
- Fetches the cheapest available fares for 5 routes out of YVR via the Travelpayouts API
- Saves each price snapshot with a timestamp so price history accumulates over time
- Tracks routes: YVR → LHR, NRT, SYD, CDG, JFK

Each run adds new rows to the database. Run it daily and you build up the price history needed for prediction.

---

## Roadmap

- [x] Route price fetching (Travelpayouts)
- [x] PostgreSQL storage with timestamped snapshots
- [ ] Scheduled data collection (cron / cloud deploy)
- [ ] World-event signal integration (conflicts, pandemics, travel bans)
- [ ] Price prediction model
- [ ] Route comparison and trend visualization
- [ ] Price alert notifications

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL running locally
- [Travelpayouts API token](https://travelpayouts.com) (free)

### Installation

```bash
git clone https://github.com/your-username/flight-tracker.git
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

### Run

```bash
go run main.go
```

### Example Output

```
🐘 Connected to PostgreSQL!

🔍 YVR → LHR
  💰 $787 | departs 2026-04-27 | 1 stop

🔍 YVR → NRT
  💰 $478 | departs 2026-09-28 | direct

🔍 YVR → SYD
  💰 $1152 | departs 2026-04-20 | 1 stop

🔍 YVR → CDG
  💰 $624 | departs 2026-04-12 | 1 stop

🔍 YVR → JFK
  💰 $327 | departs 2026-04-30 | direct

✅ Done!
```

---

## Project Structure

```
flight-tracker/
├── main.go         # Fetch prices and save to DB
├── schema.sql      # Database table definitions
├── .env.example    # Environment variable template
├── go.mod
└── go.sum
```

---

## Database Schema

```sql
routes   -- city pairs being tracked (e.g. YVR → LHR)
prices   -- price snapshots per route with timestamps
```

---

## Data Sources

| Source | Purpose |
|---|---|
| [Travelpayouts](https://travelpayouts.com) | Live flight prices by route |
| NewsAPI / GDELT | World event signals *(planned)* |

---

## Disclaimer

Flight price predictions are based on historical data and world-event signals. They are not financial advice. Always verify prices directly with airlines or booking platforms before making travel decisions.