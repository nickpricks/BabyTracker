import React, { useState, useEffect } from "react";
import { getSleep, logSleep } from "../api";

const SLEEP_TYPES = ["Nap", "Night"];
const QUALITY_OPTIONS = ["Good", "Fair", "Poor"];

const getToday = () => new Date().toISOString().slice(0, 10);
const getNow = () =>
  new Date().toLocaleTimeString("en-GB", { hour12: false }).slice(0, 5);

export default function Sleep() {
  const [date, setDate] = useState(getToday());
  const [startTime, setStartTime] = useState(getNow());
  const [endTime, setEndTime] = useState("");
  const [sleepType, setSleepType] = useState("");
  const [quality, setQuality] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");
  const [recentSleep, setRecentSleep] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchSleep = async () => {
    try {
      const entries = await getSleep();
      setRecentSleep(entries.slice(-10).reverse());
    } catch {
      // API may not be running
    }
  };

  useEffect(() => {
    fetchSleep();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setFeedback("");
    setLoading(true);
    try {
      await logSleep({
        date,
        start_time: `${date}T${startTime}:00`,
        end_time: endTime ? `${date}T${endTime}:00` : undefined,
        type: sleepType,
        quality,
        notes,
      });
      setFeedback(`Sleep logged: ${sleepType} on ${date}`);
      setDate(getToday());
      setStartTime(getNow());
      setEndTime("");
      setSleepType("");
      setQuality("");
      setNotes("");
      fetchSleep();
      setTimeout(() => setFeedback(""), 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 500, margin: "0 auto" }}>
      <h2>Log Sleep</h2>
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
              <b>Sleep Type</b>
              <select
                value={sleepType}
                onChange={(e) => setSleepType(e.target.value)}
                style={{ marginLeft: 8 }}
                required
              >
                <option value="">Select type...</option>
                {SLEEP_TYPES.map((t) => (
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
              <b>Start Time</b>
              <input
                type="time"
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                style={{ marginLeft: 8 }}
                required
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>End Time (optional)</b>
              <input
                type="time"
                value={endTime}
                onChange={(e) => setEndTime(e.target.value)}
                style={{ marginLeft: 8 }}
              />
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Quality</b>
              <select
                value={quality}
                onChange={(e) => setQuality(e.target.value)}
                style={{ marginLeft: 8 }}
              >
                <option value="">Select quality...</option>
                {QUALITY_OPTIONS.map((q) => (
                  <option key={q} value={q}>
                    {q}
                  </option>
                ))}
              </select>
            </label>
          </div>
          <div style={{ marginBottom: 12 }}>
            <label>
              <b>Notes</b>
              <textarea
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                style={{ marginLeft: 8, width: "80%", minHeight: 48 }}
                placeholder="Sleep observations..."
              />
            </label>
          </div>
          <button type="submit" disabled={loading}>
            {loading ? "Saving..." : "Log Sleep"}
          </button>
          {feedback && (
            <div style={{ color: "green", marginTop: 8 }}>{feedback}</div>
          )}
          {error && <div style={{ color: "red", marginTop: 8 }}>{error}</div>}
        </form>
      </div>
      <hr />
      <div style={{ marginTop: 24 }}>
        <h3>Recent Sleep</h3>
        <div style={{ background: "#f8f8f8", padding: 16, borderRadius: 8 }}>
          {recentSleep.length === 0 ? (
            <p style={{ color: "#666" }}>No sleep entries logged yet.</p>
          ) : (
            <ul style={{ margin: 0, padding: "0 0 0 20px" }}>
              {recentSleep.map((entry) => (
                <li key={entry.id} style={{ marginBottom: 4 }}>
                  <b>{entry.type}</b> on {entry.date}
                  {entry.quality && ` - Quality: ${entry.quality}`}
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
