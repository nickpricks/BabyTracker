import { useState, useEffect } from "react";
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
  const [fetchError, setFetchError] = useState("");

  const fetchSleep = async () => {
    try {
      setFetchError("");
      const entries = await getSleep();
      setRecentSleep(entries.slice(-10).reverse());
    } catch {
      setFetchError("Could not load sleep entries. Is the API server running?");
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
        start_time: `${date}T${startTime}:00Z`,
        end_time: endTime ? `${date}T${endTime}:00Z` : undefined,
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
    <div className="space-y-6 animate-slide-up">
      {/* Form Card */}
      <div className="card">
        <div className="flex items-center gap-3 mb-5">
          <span className="w-10 h-10 rounded-xl bg-mod-sleep/15 flex items-center justify-center text-xl">
            😴
          </span>
          <h2 className="font-display text-xl font-bold text-fg-heading">
            Log Sleep
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Sleep Type</label>
            <select
              value={sleepType}
              onChange={(e) => setSleepType(e.target.value)}
              className="input-field"
              required
            >
              <option value="">Select type...</option>
              {SLEEP_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Date</label>
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="input-field"
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">Start Time</label>
              <input
                type="time"
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                className="input-field"
                required
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">End Time (optional)</label>
              <input
                type="time"
                value={endTime}
                onChange={(e) => setEndTime(e.target.value)}
                className="input-field"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Quality</label>
            <select
              value={quality}
              onChange={(e) => setQuality(e.target.value)}
              className="input-field"
            >
              <option value="">Select quality...</option>
              {QUALITY_OPTIONS.map((q) => (
                <option key={q} value={q}>{q}</option>
              ))}
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="input-field min-h-[60px] resize-y"
              placeholder="Sleep observations..."
            />
          </div>

          <div className="pt-2">
            <button type="submit" disabled={loading} className="btn-primary w-full disabled:opacity-50">
              {loading ? "Saving..." : "Log Sleep"}
            </button>
          </div>

          {feedback && (
            <p className="text-sm text-mod-sleep font-medium animate-fade-in">{feedback}</p>
          )}
          {error && (
            <p className="text-sm text-red-500 font-medium animate-fade-in">{error}</p>
          )}
        </form>
      </div>

      {/* Recent Sleep */}
      <div className="card">
        <h3 className="font-display text-lg font-bold text-fg-heading mb-4">
          Recent Sleep
        </h3>
        {fetchError ? (
          <p className="text-sm text-red-500">{fetchError}</p>
        ) : recentSleep.length === 0 ? (
          <p className="text-sm text-fg-muted">No sleep entries logged yet.</p>
        ) : (
          <ul className="space-y-2">
            {recentSleep.map((entry) => (
              <li
                key={entry.id}
                className="flex items-baseline gap-2 py-2 border-b border-line-subtle last:border-0"
              >
                <span className="text-sm font-semibold text-fg">{entry.type}</span>
                <span className="text-xs text-fg-muted">{entry.date}</span>
                {entry.quality && (
                  <span className="text-xs text-mod-sleep font-medium">{entry.quality}</span>
                )}
                {entry.notes && (
                  <span className="text-xs text-fg-subtle truncate">{entry.notes}</span>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
