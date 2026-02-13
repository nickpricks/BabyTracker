import { API_BASE } from "./config";

async function apiGet(path) {
  const res = await fetch(`${API_BASE}${path}`);
  if (!res.ok) {
    throw new Error(`GET ${path} failed: ${res.status}`);
  }
  return res.json();
}

async function apiPost(path, body) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `POST ${path} failed: ${res.status}`);
  }
  return res.json();
}

// Feeds
export const getFeeds = () => apiGet("/feeds");
export const logFeed = (feed) => apiPost("/feeds", feed);

// Sleep
export const getSleep = () => apiGet("/sleep");
export const logSleep = (entry) => apiPost("/sleep", entry);

// Growth
export const getGrowth = () => apiGet("/growth");
export const logGrowth = (entry) => apiPost("/growth", entry);

// Diapers
export const getDiapers = () => apiGet("/diapers");
export const logDiaper = (entry) => apiPost("/diapers", entry);
