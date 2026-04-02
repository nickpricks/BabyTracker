/* ═══════════════════════════════════════════════════════════════
   Theme Definitions + Utilities
   Same architecture as Floor-Tracker (localStorage instead of Firebase)
   ═══════════════════════════════════════════════════════════════ */

import { useState, useEffect } from "react";

export const THEMES = {
  lullaby: {
    id: "lullaby",
    name: "Lullaby",
    family: "Nursery",
    darkOnly: false,
    cssClass: "theme-lullaby",
    previewColors: { bg: "#faf6ef", accent: "#e8a44a", text: "#3d3529" },
    previewColorsDark: { bg: "#1a2332", accent: "#f0b35e", text: "#d4dde8" },
  },
  "nursery-os": {
    id: "nursery-os",
    name: "Nursery_OS",
    family: "Cyberpunk",
    darkOnly: true,
    cssClass: "theme-nursery-os",
    previewColors: { bg: "#080810", accent: "#00e4ff", text: "#a8a8c0" },
  },
  "midnight-feed": {
    id: "midnight-feed",
    name: "Midnight Feed",
    family: "Nursery",
    darkOnly: true,
    cssClass: "theme-midnight-feed",
    previewColors: { bg: "#0d0b08", accent: "#cc8833", text: "#998877" },
  },
};

const STORAGE_KEY = "babytracker-theme";
const COLOR_MODE_KEY = "babytracker-color-mode";

export function getStoredTheme() {
  return localStorage.getItem(STORAGE_KEY) || "lullaby";
}

export function getStoredColorMode() {
  return localStorage.getItem(COLOR_MODE_KEY) || "system";
}

export function applyTheme(themeId, colorMode) {
  const theme = THEMES[themeId];
  if (!theme) return;

  const root = document.documentElement;

  // Remove all theme classes
  Object.values(THEMES).forEach((t) => root.classList.remove(t.cssClass));

  // Apply new theme class
  root.classList.add(theme.cssClass);

  // Handle dark mode
  if (theme.darkOnly) {
    root.classList.add("dark");
  } else if (colorMode === "dark") {
    root.classList.add("dark");
  } else if (colorMode === "light") {
    root.classList.remove("dark");
  } else {
    const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    root.classList.toggle("dark", isDark);
  }

  // Persist
  localStorage.setItem(STORAGE_KEY, themeId);
  localStorage.setItem(COLOR_MODE_KEY, colorMode);
}

export function useTheme() {
  const [themeId, setThemeId] = useState(getStoredTheme);
  const [colorMode, setColorMode] = useState(getStoredColorMode);

  useEffect(() => {
    applyTheme(themeId, colorMode);
  }, [themeId, colorMode]);

  // Listen for system color scheme changes
  useEffect(() => {
    if (colorMode !== "system") return;
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    const handler = () => applyTheme(themeId, "system");
    mq.addEventListener("change", handler);
    return () => mq.removeEventListener("change", handler);
  }, [themeId, colorMode]);

  // Apply on mount
  useEffect(() => {
    applyTheme(getStoredTheme(), getStoredColorMode());
  }, []);

  return { themeId, setThemeId, colorMode, setColorMode };
}
