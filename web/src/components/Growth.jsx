import { useState, useEffect } from "react";
import { getGrowth, logGrowth } from "../api";
import { useLoadMore } from "../useLoadMore";

const getToday = () => new Date().toISOString().slice(0, 10);

export default function Growth() {
  const [date, setDate] = useState(getToday());
  const [weight, setWeight] = useState("");
  const [height, setHeight] = useState("");
  const [headCirc, setHeadCirc] = useState("");
  const [notes, setNotes] = useState("");
  const [feedback, setFeedback] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { items: recentGrowth, loading: listLoading, error: fetchError, loadMore, hasMore, refresh, sentinelRef } = useLoadMore(getGrowth);

  useEffect(() => {
    refresh();
  }, [refresh]);

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
      refresh();
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
          <span className="w-10 h-10 rounded-xl bg-mod-growth/15 flex items-center justify-center text-xl">
            📏
          </span>
          <h2 className="font-display text-xl font-bold text-fg-heading">
            Log Growth
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
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

          <div className="grid grid-cols-3 gap-3">
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">Weight</label>
              <input
                type="number"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                className="input-field"
                min="0"
                step="0.01"
                placeholder="kg"
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">Height</label>
              <input
                type="number"
                value={height}
                onChange={(e) => setHeight(e.target.value)}
                className="input-field"
                min="0"
                step="0.1"
                placeholder="cm"
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-semibold text-fg-muted">Head</label>
              <input
                type="number"
                value={headCirc}
                onChange={(e) => setHeadCirc(e.target.value)}
                className="input-field"
                min="0"
                step="0.1"
                placeholder="cm"
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="input-field min-h-[60px] resize-y"
              placeholder="Growth observations..."
            />
          </div>

          <div className="pt-2">
            <button type="submit" disabled={loading} className="btn-primary w-full disabled:opacity-50">
              {loading ? "Saving..." : "Log Growth"}
            </button>
          </div>

          {feedback && (
            <p className="text-sm text-mod-growth font-medium animate-fade-in">{feedback}</p>
          )}
          {error && (
            <p className="text-sm text-red-500 font-medium animate-fade-in">{error}</p>
          )}
        </form>
      </div>

      {/* Recent Measurements */}
      <div className="card">
        <h3 className="font-display text-lg font-bold text-fg-heading mb-4">
          Recent Measurements
        </h3>
        {fetchError ? (
          <p className="text-sm text-red-500">{fetchError}</p>
        ) : recentGrowth.length === 0 ? (
          <p className="text-sm text-fg-muted">No growth entries logged yet.</p>
        ) : (
          <>
            <ul className="space-y-2">
              {recentGrowth.map((entry) => (
                <li
                  key={entry.id}
                  className="flex items-baseline gap-2 py-2 border-b border-line-subtle last:border-0"
                >
                  <span className="text-sm font-semibold text-fg">{entry.date}</span>
                  {entry.weight > 0 && (
                    <span className="text-xs text-mod-growth font-medium">{entry.weight} kg</span>
                  )}
                  {entry.height > 0 && (
                    <span className="text-xs text-fg-muted">{entry.height} cm</span>
                  )}
                  {entry.head_circ > 0 && (
                    <span className="text-xs text-fg-subtle">Head: {entry.head_circ} cm</span>
                  )}
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
