import { useState, useEffect } from "react";
import { getDiapers, logDiaper } from "../api";
import { useLoadMore } from "../useLoadMore";

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
  const [loading, setLoading] = useState(false);

  const { items: recentDiapers, loading: listLoading, error: fetchError, loadMore, hasMore, refresh, sentinelRef } = useLoadMore(getDiapers);

  useEffect(() => {
    refresh();
  }, [refresh]);

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
      refresh();
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
    <div className="space-y-6 animate-slide-up">
      {/* Form Card */}
      <div className="card">
        <div className="flex items-center gap-3 mb-5">
          <span className="w-10 h-10 rounded-xl bg-mod-diaper/15 flex items-center justify-center text-xl">
            🧷
          </span>
          <h2 className="font-display text-xl font-bold text-fg-heading">
            The Susu-Poty Chronicles
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Diaper Type</label>
            <select
              value={diaperType}
              onChange={(e) => setDiaperType(e.target.value)}
              className="input-field"
              required
            >
              <option value="">Select type...</option>
              {DIAPER_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-2 gap-4">
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
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">Time</label>
              <input
                type="time"
                step="1"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                className="input-field"
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="input-field min-h-[60px] resize-y"
              placeholder="Any observations..."
            />
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button type="button" onClick={quickWet} className="btn-quick border-mod-diaper/30 hover:border-mod-diaper hover:text-mod-diaper">
              Quick Wet
            </button>
            <button type="button" onClick={quickDirty} className="btn-quick border-mod-diaper/30 hover:border-mod-diaper hover:text-mod-diaper">
              Quick Dirty
            </button>
            <button type="submit" disabled={loading} className="btn-primary ml-auto disabled:opacity-50">
              {loading ? "Saving..." : "Log Change"}
            </button>
          </div>

          {feedback && (
            <p className="text-sm text-mod-diaper font-medium animate-fade-in">{feedback}</p>
          )}
          {error && (
            <p className="text-sm text-red-500 font-medium animate-fade-in">{error}</p>
          )}
        </form>
      </div>

      {/* Recent Changes */}
      <div className="card">
        <h3 className="font-display text-lg font-bold text-fg-heading mb-4">
          Recent Changes
        </h3>
        {fetchError ? (
          <p className="text-sm text-red-500">{fetchError}</p>
        ) : recentDiapers.length === 0 ? (
          <p className="text-sm text-fg-muted">No diaper changes logged yet.</p>
        ) : (
          <>
            <ul className="space-y-2">
              {recentDiapers.map((entry) => (
                <li
                  key={entry.id}
                  className="flex items-baseline gap-2 py-2 border-b border-line-subtle last:border-0"
                >
                  <span className="text-sm font-semibold text-fg">{entry.type}</span>
                  <span className="text-xs text-fg-muted">{entry.date}</span>
                  {entry.notes && (
                    <span className="text-xs text-fg-subtle truncate">{entry.notes}</span>
                  )}
                </li>
              ))}
            </ul>
            {hasMore && <div ref={sentinelRef} className="py-2 text-center text-xs text-fg-muted">{listLoading ? "Loading..." : ""}</div>}
          </>
        )}
      </div>
    </div>
  );
}
