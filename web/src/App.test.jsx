import { vi, describe, it, expect } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import App from "./App";

vi.mock("./api", () => ({
  getFeeds: vi.fn().mockResolvedValue([]),
  logFeed: vi.fn(),
  getSleep: vi.fn().mockResolvedValue([]),
  logSleep: vi.fn(),
  getGrowth: vi.fn().mockResolvedValue([]),
  logGrowth: vi.fn(),
  getDiapers: vi.fn().mockResolvedValue([]),
  logDiaper: vi.fn(),
}));

describe("App", () => {
  it("renders navigation links", () => {
    render(<App />);
    expect(screen.getByText("Baby Tracker")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /feeds/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /sleep/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /growth/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /diapers/i })).toBeInTheDocument();
  });

  it("renders dashboard by default", async () => {
    render(<App />);
    await waitFor(() => {
      expect(screen.getByText("Dashboard")).toBeInTheDocument();
    });
  });
});
