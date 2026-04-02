import { NavLink, useLocation } from "react-router-dom";
import { THEMES } from "./themes";

const NAV_ITEMS = [
  { to: "/", icon: "🏠", label: "Home" },
  { to: "/feeds", icon: "🍼", label: "Feeds" },
  { to: "/sleep", icon: "😴", label: "Sleep" },
  { to: "/growth", icon: "📏", label: "Growth" },
  { to: "/susupoty", icon: "🧷", label: "Diapers" },
];

function ThemeSwitcher({ theme }) {
  const { themeId, setThemeId, colorMode, setColorMode } = theme;
  const current = THEMES[themeId];

  return (
    <div className="flex items-center gap-2">
      <select
        value={themeId}
        onChange={(e) => setThemeId(e.target.value)}
        className="input-field !w-auto !py-1.5 !px-3 text-sm"
      >
        {Object.values(THEMES).map((t) => (
          <option key={t.id} value={t.id}>
            {t.name}
          </option>
        ))}
      </select>

      {!current?.darkOnly && (
        <div className="flex rounded-lg border border-line overflow-hidden">
          {["light", "dark", "system"].map((mode) => (
            <button
              key={mode}
              onClick={() => setColorMode(mode)}
              className={`px-2 py-1 text-xs capitalize transition-colors ${
                colorMode === mode
                  ? "bg-accent text-white"
                  : "bg-surface text-fg-muted hover:text-fg"
              }`}
            >
              {mode === "light" ? "☀" : mode === "dark" ? "🌙" : "⚙"}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

export default function MainLayout({ children, theme }) {
  const location = useLocation();

  return (
    <div className="min-h-screen bg-surface bg-topo relative">
      {/* Ambient effects layer */}
      <div className="fx-ambient" />

      {/* Header */}
      <header className="sticky top-0 z-40 bg-surface-card/80 backdrop-blur-md border-b border-line-subtle">
        <div className="max-w-3xl mx-auto px-4 py-3 flex items-center justify-between">
          <h1 className="font-display text-lg font-bold text-fg-heading tracking-tight">
            Baby Tracker
          </h1>
          <ThemeSwitcher theme={theme} />
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-3xl mx-auto px-4 pt-6 pb-24 animate-fade-in">
        {children}
      </main>

      {/* Bottom navigation */}
      <nav className="fixed bottom-0 inset-x-0 z-40 bg-surface-card/90 backdrop-blur-md border-t border-line-subtle">
        <div className="max-w-3xl mx-auto flex justify-around py-2">
          {NAV_ITEMS.map(({ to, icon, label }) => {
            const isActive =
              to === "/"
                ? location.pathname === "/"
                : location.pathname.startsWith(to);

            return (
              <NavLink
                key={to}
                to={to}
                className={`flex flex-col items-center gap-0.5 px-3 py-1 rounded-xl transition-colors ${
                  isActive
                    ? "text-accent"
                    : "text-fg-muted hover:text-fg"
                }`}
              >
                <span className="text-xl">{icon}</span>
                <span className="text-[0.65rem] font-semibold">{label}</span>
              </NavLink>
            );
          })}
        </div>
      </nav>
    </div>
  );
}
