import React, { useState, useEffect } from "react";
import { getDiapers, logDiaper } from "../api";

const DIAPER_TYPES = ["Wet", "Dirty", "Mixed"];

const getToday = () => new Date().toISOString().slice(0, 10);
const getNow = () =>
  new Date().toLocaleTimeString("en-GB", { hour12: false }).slice(0, 8);

export default function SusuPoty() {
  const [date, setDate] = useState(getToday());
  const [time, setTime] = useState(getNow());
  const [diaperType, setDiaperType] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");
  const [recentDiapers, setRecentDiapers] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchDiapers = async () => {
    try {
      const entries = await getDiapers();
      setRecentDiapers(entries.slice(-10).reverse());
    } catch {
      // API may not be running
    }
  };

  useEffect(() => {
    fetchDiapers();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setFeedback("");
    setLoading(true);
    try {
      await logDiaper({
        date,
        time: `${date}T${time}`,
        type: diaperType,
        notes,
      });
      setFeedback(`Diaper change logged: ${diaperType} on ${date}`);
      setDate(getToday());
      setTime(getNow());
      setDiaperType("");
      setNotes("");
      fetchDiapers();
      setTimeout(() => setFeedback(""), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const quickWet = () => {
    setDiaperType("Wet");
    setDate(getToday());
    setTime(getNow());
  };

  const quickDirty = () => {
    setDiaperType("Dirty");
    setDate(getToday());
    setTime(getNow());
  };

  return (
    <div style={{ maxWidth: 500, margin: "0 auto" }}>
      <h2>The Susu-Poty Chronicles</h2>
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
              <b>Diaper Type</b>
              <select
                value={diaperType}
                onChange={(e) => setDiaperType(e.target.value)}
                style={{ marginLeft: 8 }}
                required
              >
                <option value="">Select type...</option>
                {DIAPER_TYPES.map((t) => (
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
              <b>Notes</b>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                style={{ marginLeft: 8, width: "80%", minHeight: 48 }}
                placeholder="Any observations..."
              />
            </label>
          </div>
          <div style={{ display: "flex", gap: 12, marginBottom: 12 }}>
            <button type="button" onClick={quickWet}>
              Quick Wet
            </button>
            <button type="button" onClick={quickDirty}>
              Quick Dirty
            </button>
            <button
              type="submit"
              style={{ marginLeft: "auto" }}
              disabled={loading}
            >
              {loading ? "Saving..." : "Log Change"}
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
        <h3>Recent Changes</h3>
        <div style={{ background: "#f8f8f8", padding: 16, borderRadius: 8 }}>
          {recentDiapers.length === 0 ? (
            <p style={{ color: "#666" }}>No diaper changes logged yet.</p>
          ) : (
            <ul style={{ margin: 0, padding: "0 0 0 20px" }}>
              {recentDiapers.map((entry) => (
                <li key={entry.id} style={{ marginBottom: 4 }}>
                  <b>{entry.type}</b> on {entry.date}
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
