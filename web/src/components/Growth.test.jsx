import { vi, describe, it, expect, beforeEach } from "vitest";
import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Growth from "./Growth";

vi.mock("../api", () => ({
  getGrowth: vi.fn(),
  logGrowth: vi.fn(),
}));

const { getGrowth, logGrowth } = await import("../api");

beforeEach(() => {
  getGrowth.mockReset();
  logGrowth.mockReset();
});

describe("Growth", () => {
  it("renders the form", () => {
    getGrowth.mockResolvedValue([]);
    render(<Growth />);
    expect(screen.getByRole("heading", { name: /log growth/i })).toBeInTheDocument();
  });

  it("shows fetch error when API fails", async () => {
    getGrowth.mockRejectedValue(new Error("fail"));
    render(<Growth />);
    await waitFor(() => {
      expect(screen.getByText(/could not load growth/i)).toBeInTheDocument();
    });
  });

  it("shows recent entries on successful load", async () => {
    getGrowth.mockResolvedValue([
      { id: 1, date: "2026-01-01", weight: 4.5, height: 55, head_circ: 36, notes: "" },
    ]);
    render(<Growth />);
    await waitFor(() => {
      expect(screen.getByText(/4.5 kg/)).toBeInTheDocument();
    });
  });

  it("submits growth entry", async () => {
    const user = userEvent.setup();
    getGrowth.mockResolvedValue([]);
    logGrowth.mockResolvedValue({ id: 1 });

    render(<Growth />);
    await user.click(screen.getByRole("button", { name: /log growth/i }));

    await waitFor(() => {
      expect(logGrowth).toHaveBeenCalled();
    });
  });

  it("shows error on submit failure", async () => {
    const user = userEvent.setup();
    getGrowth.mockResolvedValue([]);
    logGrowth.mockRejectedValue(new Error("date required"));

    render(<Growth />);
    await user.click(screen.getByRole("button", { name: /log growth/i }));

    await waitFor(() => {
      expect(screen.getByText("date required")).toBeInTheDocument();
    });
  });
});
