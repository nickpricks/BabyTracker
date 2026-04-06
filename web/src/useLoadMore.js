import { useState, useCallback, useRef, useEffect } from "react";

// Configurable pagination defaults
const INITIAL_LIMIT = 5;
const BULK_LIMIT = 50;
const BULK_THRESHOLD = 3; // switch to bulk after this many clicks

/**
 * Hook for "Load More" pagination with escalation and infinite scroll.
 * After BULK_THRESHOLD loads, fetches BULK_LIMIT items per load instead of INITIAL_LIMIT.
 * Returns a sentinelRef — attach to a div at the bottom of your list for auto-load on scroll.
 *
 * @param {Function} fetchFn - API function that accepts (limit, offset) and returns { items, total }
 * @returns {{ items, total, loading, error, loadMore, hasMore, refresh, sentinelRef }}
 */
export function useLoadMore(fetchFn) {
  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [loadCount, setLoadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const sentinelRef = useRef(null);
  const loadingRef = useRef(false);

  const hasMore = offset < total;

  const refresh = useCallback(async () => {
    setLoading(true);
    loadingRef.current = true;
    setError("");
    try {
      const res = await fetchFn(INITIAL_LIMIT, 0);
      setItems(res.items || []);
      setTotal(res.total);
      setOffset(res.items?.length || 0);
      setLoadCount(0);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
  }, [fetchFn]);

  const loadMore = useCallback(async () => {
    if (loadingRef.current) return;
    setLoading(true);
    loadingRef.current = true;
    setError("");
    try {
      const nextCount = loadCount + 1;
      const limit = nextCount >= BULK_THRESHOLD ? BULK_LIMIT : INITIAL_LIMIT;
      const res = await fetchFn(limit, offset);
      setItems((prev) => [...prev, ...(res.items || [])]);
      setTotal(res.total);
      setOffset((prev) => prev + (res.items?.length || 0));
      setLoadCount(nextCount);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
  }, [fetchFn, offset, loadCount]);

  // IntersectionObserver for infinite scroll
  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !loadingRef.current) {
          loadMore();
        }
      },
      { rootMargin: "100px" },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [loadMore]);

  return { items, total, loading, error, loadMore, hasMore, refresh, sentinelRef };
}
