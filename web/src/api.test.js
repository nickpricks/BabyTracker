import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";

// Must mock before importing the module under test
vi.stubGlobal("fetch", vi.fn());

// Dynamic import after mock setup
const { getFeeds, logFeed, getSleep, logSleep, getGrowth, logGrowth, getDiapers, logDiaper } = await import("./api");

beforeEach(() => {
  fetch.mockReset();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("apiGet", () => {
  it("returns parsed JSON on success", async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([{ id: 1, type: "Bottle" }]),
    });
    const result = await getFeeds();
    expect(result).toEqual([{ id: 1, type: "Bottle" }]);
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/feeds"),
      expect.objectContaining({ headers: expect.any(Object) })
    );
  });

  it("throws on non-ok response", async () => {
    fetch.mockResolvedValueOnce({ ok: false, status: 500 });
    await expect(getFeeds()).rejects.toThrow(/GET \/feeds.*failed: 500/);
  });
});

describe("apiPost", () => {
  it("sends JSON body and returns parsed response", async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ id: 1 }),
    });
    const result = await logFeed({ type: "Bottle", date: "2026-01-01" });
    expect(result).toEqual({ id: 1 });
    const [, opts] = fetch.mock.calls[0];
    expect(opts.method).toBe("POST");
    expect(JSON.parse(opts.body)).toEqual({ type: "Bottle", date: "2026-01-01" });
  });

  it("throws with server error message", async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 400,
      json: () => Promise.resolve({ error: "type is required" }),
    });
    await expect(logFeed({})).rejects.toThrow("type is required");
  });

  it("throws with status on unparseable error", async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: () => Promise.reject(new Error("not json")),
    });
    await expect(logFeed({})).rejects.toThrow("POST /feeds failed: 500");
  });
});

describe("auth header", () => {
  it("does not send Authorization when no key set", async () => {
    const origKey = import.meta.env.VITE_API_KEY;
    import.meta.env.VITE_API_KEY = "";

    // Re-import api.js to pick up the cleared key
    vi.resetModules();
    const { getFeeds: getFeedsNoAuth } = await import("./api.js");

    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });
    await getFeedsNoAuth();
    const [, opts] = fetch.mock.calls[0];
    expect(opts.headers.Authorization).toBeUndefined();

    import.meta.env.VITE_API_KEY = origKey;
  });
});

describe("all resource endpoints", () => {
  it.each([
    ["getFeeds", getFeeds, "/feeds"],
    ["getSleep", getSleep, "/sleep"],
    ["getGrowth", getGrowth, "/growth"],
    ["getDiapers", getDiapers, "/diapers"],
  ])("%s calls correct endpoint", async (name, fn, path) => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });
    await fn();
    expect(fetch.mock.calls[0][0]).toContain(path);
  });

  it.each([
    ["logFeed", logFeed, "/feeds"],
    ["logSleep", logSleep, "/sleep"],
    ["logGrowth", logGrowth, "/growth"],
    ["logDiaper", logDiaper, "/diapers"],
  ])("%s POSTs to correct endpoint", async (name, fn, path) => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ id: 1 }),
    });
    await fn({ test: true });
    const [url, opts] = fetch.mock.calls[0];
    expect(url).toContain(path);
    expect(opts.method).toBe("POST");
  });
});
