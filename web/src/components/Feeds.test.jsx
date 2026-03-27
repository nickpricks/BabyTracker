import { vi, describe, it, expect, beforeEach } from "vitest";
import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Feeds from "./Feeds";

vi.mock("../api", () => ({
  getFeeds: vi.fn(),
  logFeed: vi.fn(),
}));

const { getFeeds, logFeed } = await import("../api");

beforeEach(() => {
  getFeeds.mockReset();
  logFeed.mockReset();
});

describe("Feeds", () => {
  it("renders the form", () => {
    getFeeds.mockResolvedValue([]);
    render(<Feeds />);
    expect(screen.getByText("Log New Feed")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /log feed/i })).toBeInTheDocument();
  });

  it("shows fetch error when API fails", async () => {
    getFeeds.mockRejectedValue(new Error("network error"));
    render(<Feeds />);
    await waitFor(() => {
      expect(screen.getByText(/could not load feeds/i)).toBeInTheDocument();
    });
  });

  it("shows recent feeds on successful load", async () => {
    getFeeds.mockResolvedValue([
      { id: 1, type: "Bottle", date: "2026-01-01", quantity: 120, notes: "" },
    ]);
    render(<Feeds />);
    await waitFor(() => {
      expect(screen.getByText("2026-01-01", { exact: false })).toBeInTheDocument();
    });
  });

  it("shows empty state when no feeds", async () => {
    getFeeds.mockResolvedValue([]);
    render(<Feeds />);
    await waitFor(() => {
      expect(screen.getByText(/no feeds logged yet/i)).toBeInTheDocument();
    });
  });

  it("submits a feed and shows feedback", async () => {
    const user = userEvent.setup();
    getFeeds.mockResolvedValue([]);
    logFeed.mockResolvedValue({ id: 1 });

    render(<Feeds />);

    await user.selectOptions(screen.getByRole("combobox"), "Bottle");
    await user.click(screen.getByRole("button", { name: /log feed/i }));

    await waitFor(() => {
      expect(logFeed).toHaveBeenCalledWith(
        expect.objectContaining({ type: "Bottle" })
      );
    });
  });

  it("shows error on submit failure", async () => {
    const user = userEvent.setup();
    getFeeds.mockResolvedValue([]);
    logFeed.mockRejectedValue(new Error("type is required"));

    render(<Feeds />);

    await user.selectOptions(screen.getByRole("combobox"), "Bottle");
    await user.click(screen.getByRole("button", { name: /log feed/i }));

    await waitFor(() => {
      expect(screen.getByText("type is required")).toBeInTheDocument();
    });
  });

  it("quick bottle sets type", async () => {
    const user = userEvent.setup();
    getFeeds.mockResolvedValue([]);
    render(<Feeds />);

    await user.click(screen.getByRole("button", { name: /quick bottle/i }));
    expect(screen.getByRole("combobox").value).toBe("Bottle");
  });
});
