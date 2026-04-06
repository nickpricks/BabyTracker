import { vi, describe, it, expect, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Sleep from "./Sleep";

vi.mock("../api", () => ({
  getSleep: vi.fn(),
  logSleep: vi.fn(),
}));

const { getSleep, logSleep } = await import("../api");

beforeEach(() => {
  getSleep.mockReset();
  logSleep.mockReset();
});

describe("Sleep", () => {
  it("renders the form", () => {
    getSleep.mockResolvedValue({ items: [], total: 0 });
    render(<Sleep />);
    expect(screen.getByRole("heading", { name: /log sleep/i })).toBeInTheDocument();
  });

  it("shows fetch error when API fails", async () => {
    getSleep.mockRejectedValue(new Error("fail"));
    render(<Sleep />);
    await waitFor(() => {
      expect(screen.getByText("fail")).toBeInTheDocument();
    });
  });

  it("shows recent entries on successful load", async () => {
    getSleep.mockResolvedValue({
      items: [{ id: 1, type: "Nap", date: "2026-01-01", quality: "Good", notes: "" }],
      total: 1,
    });
    render(<Sleep />);
    await waitFor(() => {
      expect(screen.getByText("2026-01-01", { exact: false })).toBeInTheDocument();
    });
  });

  it("submits sleep entry", async () => {
    const user = userEvent.setup();
    getSleep.mockResolvedValue({ items: [], total: 0 });
    logSleep.mockResolvedValue({ id: 1 });

    render(<Sleep />);
    await user.selectOptions(screen.getAllByRole("combobox")[0], "Nap");
    await user.click(screen.getByRole("button", { name: /log sleep/i }));

    await waitFor(() => {
      expect(logSleep).toHaveBeenCalledWith(
        expect.objectContaining({ type: "Nap" })
      );
    });
  });

  it("shows error on submit failure", async () => {
    const user = userEvent.setup();
    getSleep.mockResolvedValue({ items: [], total: 0 });
    logSleep.mockRejectedValue(new Error("date required"));

    render(<Sleep />);
    await user.selectOptions(screen.getAllByRole("combobox")[0], "Night");
    await user.click(screen.getByRole("button", { name: /log sleep/i }));

    await waitFor(() => {
      expect(screen.getByText("date required")).toBeInTheDocument();
    });
  });
});
