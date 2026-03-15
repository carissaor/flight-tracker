import { DESTINATION_EMOJI, DESTINATION_LABELS } from "../constants";

export default function RouteCard({ route, selected, onClick }) {
  return (
    <button
      className={`route-card ${selected ? "selected" : ""}`}
      onClick={onClick}
    >
      <div className="route-card-top">
        <span className="flag">
          {DESTINATION_EMOJI[route.destination] || "🌍"}
        </span>
        <span className="destination">
          {DESTINATION_LABELS[route.destination] || route.destination}
        </span>
      </div>
      <div className="route-card-price">
        ${route.latest_price.toLocaleString()}
      </div>
      <div className="route-card-meta">
        <span className="depart-date">Cheapest on {route.depart_date}</span>
      </div>
    </button>
  );
}