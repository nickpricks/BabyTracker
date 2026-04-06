import { API_BASE } from "./config";

const API_KEY = import.meta.env.VITE_API_KEY || "";

function authHeaders() {
  const h = {};
  if (API_KEY) h["Authorization"] = `Bearer ${API_KEY}`;
  return h;
}

async function apiGet(path) {
  const res = await fetch(`${API_BASE}${path}`, { headers: authHeaders() });
  if (!res.ok) {
    throw new Error(`GET ${path} failed: ${res.status}`);
  }
  return res.json();
}

async function apiPost(path, body) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `POST ${path} failed: ${res.status}`);
  }
  return res.json();
}

async function apiPut(path, body) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `PUT ${path} failed: ${res.status}`);
  }
  return res.json();
}

async function apiDelete(path) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "DELETE",
    headers: authHeaders(),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `DELETE ${path} failed: ${res.status}`);
  }
  return res.json();
}

// Feeds
export const getFeeds = (limit, offset) => apiGet(`/feeds?limit=${limit ?? 10}&offset=${offset ?? 0}`);
export const logFeed = (feed) => apiPost("/feeds", feed);
export const updateFeed = (id, feed) => apiPut(`/feeds/${id}`, feed);
export const deleteFeed = (id) => apiDelete(`/feeds/${id}`);

// Sleep
export const getSleep = (limit, offset) => apiGet(`/sleep?limit=${limit ?? 10}&offset=${offset ?? 0}`);
export const logSleep = (entry) => apiPost("/sleep", entry);
export const updateSleep = (id, entry) => apiPut(`/sleep/${id}`, entry);
export const deleteSleep = (id) => apiDelete(`/sleep/${id}`);

// Growth
export const getGrowth = (limit, offset) => apiGet(`/growth?limit=${limit ?? 10}&offset=${offset ?? 0}`);
export const logGrowth = (entry) => apiPost("/growth", entry);
export const updateGrowth = (id, entry) => apiPut(`/growth/${id}`, entry);
export const deleteGrowth = (id) => apiDelete(`/growth/${id}`);

// Diapers
export const getDiapers = (limit, offset) => apiGet(`/diapers?limit=${limit ?? 10}&offset=${offset ?? 0}`);
export const logDiaper = (entry) => apiPost("/diapers", entry);
export const updateDiaper = (id, entry) => apiPut(`/diapers/${id}`, entry);
export const deleteDiaper = (id) => apiDelete(`/diapers/${id}`);
