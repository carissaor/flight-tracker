import { useState } from "react";
import axios from "axios";
import { DESTINATION_LABELS, DESTINATION_EMOJI } from '../constants';

const API = import.meta.env.VITE_API_URL;

export default function SearchPanel() {
  const [destination, setDestination] = useState("LHR");
  const [month, setMonth] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSearch = () => {
    if (!month) return;
    setLoading(true);
    setError(null);
    setResults(null);

    axios
      .get(`${API}/api/search`, {
        params: { origin: "YVR", destination, month },
      })
      .then((res) => {
        setResults(res.data.results || []);
        setLoading(false);
      })
      .catch(() => {
        setError("Failed to fetch prices. Try again.");
        setLoading(false);
      });
  };

  return (
    <div className="search-panel">
      <div className="search-controls">
        <div className="search-field">
          <label className="search-label">From</label>
          <div className="search-static">YVR — Vancouver</div>
        </div>

        <div className="search-field">
          <label className="search-label">To</label>
          <select
            className="search-select"
            value={destination}
            onChange={(e) => setDestination(e.target.value)}
          >
            {Object.entries(DESTINATION_LABELS).map(([code, name]) => (
              <option key={code} value={code}>
                {DESTINATION_EMOJI[code]} {name} ({code})
              </option>
            ))}
          </select>
        </div>

        <div className="search-field">
          <label className="search-label">Departure Month</label>
          <input
            type="month"
            className="search-input"
            value={month}
            onChange={(e) => setMonth(e.target.value)}
            min={new Date().toISOString().slice(0, 7)}
          />
        </div>

        <button
          className="search-btn"
          onClick={handleSearch}
          disabled={!month || loading}
        >
          {loading ? "Searching..." : "Search Flights"}
        </button>
      </div>

      {error && <div className="search-error">{error}</div>}

      {results && results.length === 0 && (
        <div className="search-empty">No flights found for this month.</div>
      )}

      {results && results.length > 0 && (
        <ul className="search-results">
          {results.map((r, i) => (
            <li key={i} className="search-result-item">
              <div className="result-route">
                <span className="result-flag">{DESTINATION_EMOJI[r.destination]}</span>
                <span className="result-dest">YVR → {r.destination}</span>
                <span className="result-airline">{r.airline}</span>
              </div>
              <div className="result-right">
                <div className="result-price">${r.price.toLocaleString()}</div>
                <div className="result-meta">
                  {r.depart_date?.slice(0, 10)} · {r.transfers === 0 ? "Direct" : `${r.transfers} stop${r.transfers > 1 ? "s" : ""}`}
                </div>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}