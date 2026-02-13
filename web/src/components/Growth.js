import React, { useState, useEffect } from "react";
import { getGrowth, logGrowth } from "../api";

const getToday = () => new Date().toISOString().slice(0, 10);

export default function Growth() {
  const [date, setDate] = useState(getToday());
  const [weight, setWeight] = useState("");
  const [height, setHeight] = useState("");
  const [headCirc, setHeadCirc] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");
  const [recentGrowth, setRecentGrowth] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchGrowth = async () => {
    try {
      const entries = await getGrowth();
      setRecentGrowth(entries.slice(-10).reverse());
    } catch {
      // API may not be running
    }
  };

  useEffect(() => {
    fetchGrowth();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setFeedback("");
    setLoading(true);
    try {
      await logGrowth({
        date,
        weight: weight ? parseFloat(weight) : 0,
        height: height ? parseFloat(height) : 0,
        head_circ: headCirc ? parseFloat(headCirc) : 0,
        notes,
      });
      setFeedback(`Growth logged for ${date}`);
      setDate(getToday());
      setWeight("");
      setHeight("");
      setHeadCirc("");
      setNotes("");
      fetchGrowth();
      setTimeout(() => setFeedback(""), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 500, margin: "0 auto" }}>
      <h2>Log Growth</h2>
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
              <b>Weight (kg)</b>
              <input
                type="number"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                style={{ marginLeft: 8, width: 100 }}
                min="0"
                step="0.01"
                placeholder="kg"
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Height (cm)</b>
              <input
                type="number"
                value={height}
                onChange={(e) => setHeight(e.target.value)}
                style={{ marginLeft: 8, width: 100 }}
                min="0"
                step="0.1"
                placeholder="cm"
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Head Circumference (cm)</b>
              <input
                type="number"
                value={headCirc}
                onChange={(e) => setHeadCirc(e.target.value)}
                style={{ marginLeft: 8, width: 100 }}
                min="0"
                step="0.1"
                placeholder="cm"
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
                placeholder="Growth observations..."
              />
            </label>
          </div>
          <button type="submit" disabled={loading}>
            {loading ? "Saving..." : "Log Growth"}
          </button>
          {feedback && (
            <div style={{ color: "green", marginTop: 8 }}>{feedback}</div>
          )}
          {error && <div style={{ color: "red", marginTop: 8 }}>{error}</div>}
        </form>
      </div>
      <hr />
      <div style={{ marginTop: 24 }}>
        <h3>Recent Measurements</h3>
        <div style={{ background: "#f8f8f8", padding: 16, borderRadius: 8 }}>
          {recentGrowth.length === 0 ? (
            <p style={{ color: "#666" }}>No growth entries logged yet.</p>
          ) : (
            <ul style={{ margin: 0, padding: "0 0 0 20px" }}>
              {recentGrowth.map((entry) => (
                <li key={entry.id} style={{ marginBottom: 4 }}>
                  <b>{entry.date}</b>
                  {entry.weight > 0 && ` - ${entry.weight} kg`}
                  {entry.height > 0 && ` - ${entry.height} cm`}
                  {entry.head_circ > 0 && ` - Head: ${entry.head_circ} cm`}
                  {entry.notes && ` - ${entry.notes}`}
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
