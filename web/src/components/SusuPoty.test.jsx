import { vi, describe, it, expect, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import SusuPoty from "./SusuPoty";

vi.mock("../api", () => ({
  getDiapers: vi.fn(),
  logDiaper: vi.fn(),
}));

const { getDiapers, logDiaper } = await import("../api");

beforeEach(() => {
  getDiapers.mockReset();
  logDiaper.mockReset();
});

describe("SusuPoty", () => {
  it("renders the form", () => {
    getDiapers.mockResolvedValue([]);
    render(<SusuPoty />);
    expect(screen.getByText("The Susu-Poty Chronicles")).toBeInTheDocument();
  });

  it("shows fetch error when API fails", async () => {
    getDiapers.mockRejectedValue(new Error("fail"));
    render(<SusuPoty />);
    await waitFor(() => {
      expect(screen.getByText(/could not load diaper/i)).toBeInTheDocument();
    });
  });

  it("shows recent entries on successful load", async () => {
    getDiapers.mockResolvedValue([
      { id: 1, type: "Wet", date: "2026-01-01", notes: "" },
    ]);
    render(<SusuPoty />);
    await waitFor(() => {
      expect(screen.getByText("2026-01-01", { exact: false })).toBeInTheDocument();
    });
  });

  it("submits diaper entry", async () => {
    const user = userEvent.setup();
    getDiapers.mockResolvedValue([]);
    logDiaper.mockResolvedValue({ id: 1 });

    render(<SusuPoty />);
    await user.selectOptions(screen.getByRole("combobox"), "Wet");
    await user.click(screen.getByRole("button", { name: /log change/i }));

    await waitFor(() => {
      expect(logDiaper).toHaveBeenCalledWith(
        expect.objectContaining({ type: "Wet" })
      );
    });
  });

  it("quick wet sets type", async () => {
    const user = userEvent.setup();
    getDiapers.mockResolvedValue([]);
    render(<SusuPoty />);

    await user.click(screen.getByRole("button", { name: /quick wet/i }));
    expect(screen.getByRole("combobox").value).toBe("Wet");
  });

  it("shows error on submit failure", async () => {
    const user = userEvent.setup();
    getDiapers.mockResolvedValue([]);
    logDiaper.mockRejectedValue(new Error("type required"));

    render(<SusuPoty />);
    await user.selectOptions(screen.getByRole("combobox"), "Dirty");
    await user.click(screen.getByRole("button", { name: /log change/i }));

    await waitFor(() => {
      expect(screen.getByText("type required")).toBeInTheDocument();
    });
  });
});
