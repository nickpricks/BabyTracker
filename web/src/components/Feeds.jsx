import { useState, useEffect } from "react";
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
  const [fetchError, setFetchError] = useState("");

  const fetchFeeds = async () => {
    try {
      setFetchError("");
      const feeds = await getFeeds();
      setRecentFeeds(feeds.slice(-10).reverse());
    } catch {
      setFetchError("Could not load feeds. Is the API server running?");
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
        time: `${date}T${time}Z`,
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
    <div className="space-y-6 animate-slide-up">
      {/* Form Card */}
      <div className="card">
        <div className="flex items-center gap-3 mb-5">
          <span className="w-10 h-10 rounded-xl bg-mod-feeds/15 flex items-center justify-center text-xl">
            🍼
          </span>
          <h2 className="font-display text-xl font-bold text-fg-heading">
            Log Feed
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Feed Type</label>
            <select
              value={feedType}
              onChange={(e) => setFeedType(e.target.value)}
              className="input-field"
              required
            >
              <option value="">Select feed type...</option>
              {FEED_TYPES.map((t) => (
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
            <label className="text-sm font-semibold text-fg-muted">Quantity (optional)</label>
            <input
              type="number"
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              className="input-field"
              min="0"
              step="any"
              placeholder="ml or oz"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-semibold text-fg-muted">Notes</label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="input-field min-h-[60px] resize-y"
              placeholder="How did baby respond? Any concerns?"
            />
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button type="button" onClick={quickBottle} className="btn-quick border-mod-feeds/30 hover:border-mod-feeds hover:text-mod-feeds">
              Quick Bottle
            </button>
            <button type="button" onClick={quickBreast} className="btn-quick border-mod-feeds/30 hover:border-mod-feeds hover:text-mod-feeds">
              Quick Breast
            </button>
            <button type="submit" disabled={loading} className="btn-primary ml-auto disabled:opacity-50">
              {loading ? "Saving..." : "Log Feed"}
            </button>
          </div>

          {feedback && (
            <p className="text-sm text-mod-feeds font-medium animate-fade-in">{feedback}</p>
          )}
          {error && (
            <p className="text-sm text-red-500 font-medium animate-fade-in">{error}</p>
          )}
        </form>
      </div>

      {/* Recent Feeds */}
      <div className="card">
        <h3 className="font-display text-lg font-bold text-fg-heading mb-4">
          Recent Feeds
        </h3>
        {fetchError ? (
          <p className="text-sm text-red-500">{fetchError}</p>
        ) : recentFeeds.length === 0 ? (
          <p className="text-sm text-fg-muted">No feeds logged yet.</p>
        ) : (
          <ul className="space-y-2">
            {recentFeeds.map((feed) => (
              <li
                key={feed.id}
                className="flex items-baseline gap-2 py-2 border-b border-line-subtle last:border-0"
              >
                <span className="text-sm font-semibold text-fg">{feed.type}</span>
                <span className="text-xs text-fg-muted">{feed.date}</span>
                {feed.quantity > 0 && (
                  <span className="text-xs text-mod-feeds font-medium">{feed.quantity} ml/oz</span>
                )}
                {feed.notes && (
                  <span className="text-xs text-fg-subtle truncate">{feed.notes}</span>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
