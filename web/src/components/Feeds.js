import React, { useState, useEffect } from "react";
import { getFeeds, logFeed } from "../api";

const FEED_TYPES = [
  "Bottle",
  "Breast (Left)",
  "Breast (Right)",
  "Breast (Both)",
  "Solid Food",
];

const getToday = () => new Date().toISOString().slice(0, 10);
const getNow = () =>
  new Date().toLocaleTimeString("en-GB", { hour12: false }).slice(0, 8);

export default function Feeds() {
  const [feedType, setFeedType] = useState("");
  const [date, setDate] = useState(getToday());
  const [time, setTime] = useState(getNow());
  const [quantity, setQuantity] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");
  const [recentFeeds, setRecentFeeds] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchFeeds = async () => {
    try {
      const feeds = await getFeeds();
      setRecentFeeds(feeds.slice(-10).reverse());
    } catch {
      // API may not be running - that's okay
    }
  };

  useEffect(() => {
    fetchFeeds();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setFeedback("");
    setLoading(true);
    try {
      await logFeed({
        type: feedType,
        date: date,
        time: `${date}T${time}`,
        quantity: quantity ? parseFloat(quantity) : 0,
        notes: notes,
      });
      setFeedback(`Feed logged: ${feedType} on ${date} at ${time}`);
      setFeedType("");
      setDate(getToday());
      setTime(getNow());
      setQuantity("");
      setNotes("");
      fetchFeeds();
      setTimeout(() => setFeedback(""), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const quickBottle = () => {
    setFeedType("Bottle");
    setDate(getToday());
    setTime(getNow());
  };

  const quickBreast = () => {
    setFeedType("Breast (Both)");
    setDate(getToday());
    setTime(getNow());
  };

  return (
    <div style={{ maxWidth: 500, margin: "0 auto" }}>
      <h2>Log New Feed</h2>
      <div
        style={{
          background: "#f8f8f8",
          padding: 20,
          borderRadius: 8,
          marginBottom: 24,
        }}
      >
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Feed Type</b>
              <select
                value={feedType}
                onChange={(e) => setFeedType(e.target.value)}
                style={{ marginLeft: 8, width: "60%" }}
                required
              >
                <option value="">Select feed type...</option>
                {FEED_TYPES.map((t) => (
                  <option key={t} value={t}>
                    {t}
                  </option>
                ))}
              </select>
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Date</b>
              <input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                style={{ marginLeft: 8 }}
                required
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Time</b>
              <input
                type="time"
                step="1"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                style={{ marginLeft: 8 }}
                required
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Quantity (optional)</b>
              <input
                type="number"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                style={{ marginLeft: 8, width: 100 }}
                min="0"
                step="any"
                placeholder="ml or oz"
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Notes</b>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                style={{ marginLeft: 8, width: "80%", minHeight: 48 }}
                placeholder="How did baby respond? Any concerns?"
              />
            </label>
          </div>
          <div style={{ display: "flex", gap: 12, marginBottom: 12 }}>
            <button type="button" onClick={quickBottle}>
              Quick Bottle
            </button>
            <button type="button" onClick={quickBreast}>
              Quick Breast
            </button>
            <button
              type="submit"
              style={{ marginLeft: "auto" }}
              disabled={loading}
            >
              {loading ? "Saving..." : "Log Feed"}
            </button>
          </div>
          {feedback && (
            <div style={{ color: "green", marginTop: 8 }}>{feedback}</div>
          )}
          {error && <div style={{ color: "red", marginTop: 8 }}>{error}</div>}
        </form>
      </div>
      <hr />
      <div style={{ marginTop: 24 }}>
        <h3>Recent Feeds</h3>
        <div style={{ background: "#f8f8f8", padding: 16, borderRadius: 8 }}>
          {recentFeeds.length === 0 ? (
            <p style={{ color: "#666" }}>No feeds logged yet.</p>
          ) : (
            <ul style={{ margin: 0, padding: "0 0 0 20px" }}>
              {recentFeeds.map((feed) => (
                <li key={feed.id} style={{ marginBottom: 4 }}>
                  <b>{feed.type}</b> on {feed.date}
                  {feed.quantity > 0 && ` - ${feed.quantity} ml/oz`}
                  {feed.notes && ` - ${feed.notes}`}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
