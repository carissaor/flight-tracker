package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ---------------------------------------------------------------------------
// Travelpayouts types
//
// The /v1/prices/cheap endpoint returns a map keyed by destination IATA code.
// Each value is a map keyed by a stringified index ("0", "1", ...) containing
// the cheapest price found for that departure date.
// ---------------------------------------------------------------------------

type PriceEntry struct {
	Price      float64 `json:"price"`
	DepartDate string  `json:"departure_at"`
	Airline    string  `json:"airline"`
	FlightNum  int     `json:"flight_number"`
	Transfers  int     `json:"transfers"`
}

// PriceResponse is map[destination]map[index]PriceEntry
type PriceResponse struct {
	Success bool                             `json:"success"`
	Data    map[string]map[string]PriceEntry `json:"data"`
}

// ---------------------------------------------------------------------------
// Fetch
// ---------------------------------------------------------------------------

func fetchPrices(token, origin, destination string) ([]PriceEntry, error) {
	url := fmt.Sprintf(
		"https://api.travelpayouts.com/v1/prices/cheap?origin=%s&destination=%s&token=%s&currency=usd",
		origin, destination, token,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Access-Token", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PriceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w\nraw body: %s", err, string(body))
	}
	if !result.Success {
		return nil, fmt.Errorf("travelpayouts returned success=false; body: %s", string(body))
	}

	var entries []PriceEntry
	for _, byIndex := range result.Data {
		for _, entry := range byIndex {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

// ---------------------------------------------------------------------------
// Save
// ---------------------------------------------------------------------------

// ensureRoute inserts the route if it doesn't exist and returns its id.
func ensureRoute(db *sql.DB, origin, destination string) (int, error) {
	var id int
	err := db.QueryRow(`
		INSERT INTO routes (origin, destination)
		VALUES ($1, $2)
		ON CONFLICT (origin, destination) DO UPDATE SET origin = EXCLUDED.origin
		RETURNING id`,
		origin, destination,
	).Scan(&id)
	return id, err
}

func savePrices(db *sql.DB, routeID int, entries []PriceEntry) {
	for _, e := range entries {
		var departDate *time.Time
		if e.DepartDate != "" {
			t, err := time.Parse(time.RFC3339, e.DepartDate)
			if err == nil {
				departDate = &t
			}
		}

		_, err := db.Exec(`
			INSERT INTO prices (route_id, price, currency, depart_date, fetched_at)
			VALUES ($1, $2, 'USD', $3, NOW())`,
			routeID, e.Price, departDate,
		)
		if err != nil {
			log.Println("Error inserting price:", err)
		} else {
			fmt.Printf("  💰 $%.0f | departs %s | %s\n",
				e.Price, e.DepartDate[:10], transferLabel(e.Transfers))
		}
	}
}

func transferLabel(n int) string {
	switch n {
	case 0:
		return "direct"
	case 1:
		return "1 stop"
	default:
		return fmt.Sprintf("%d stops", n)
	}
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Cannot reach database:", err)
	}
	fmt.Println("🐘 Connected to PostgreSQL!")

	origin := os.Getenv("ORIGIN")
	if origin == "" {
		log.Fatal("ORIGIN must be set in .env")
	}
	token := os.Getenv("TRAVELPAYOUTS_TOKEN")
	if token == "" {
		log.Fatal("TRAVELPAYOUTS_TOKEN must be set in .env")
	}

	// Routes to track — TODO: make this user configurable
	destinations := []string{
		"LHR", // London
		"NRT", // Tokyo
		"SYD", // Sydney
		"CDG", // Paris
		"JFK", // New York
		"HKG", // Hong Kong
	}

	for _, dest := range destinations {
		fmt.Printf("\n🔍 %s → %s\n", origin, dest)

		routeID, err := ensureRoute(db, origin, dest)
		if err != nil {
			log.Printf("  Error ensuring route: %v", err)
			continue
		}

		entries, err := fetchPrices(token, origin, dest)
		if err != nil {
			log.Printf("  Error fetching prices: %v", err)
			continue
		}
		if len(entries) == 0 {
			fmt.Println("  ⚠️  No prices found")
			continue
		}

		savePrices(db, routeID, entries)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n✅ Done!")
}