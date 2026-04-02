import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { getFeeds, getSleep, getGrowth, getDiapers } from "../api";

function SummaryCard({ icon, title, color, linkTo, children }) {
  return (
    <div className="card-hover group">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2.5">
          <span
            className="w-9 h-9 rounded-lg flex items-center justify-center text-lg"
            style={{ background: `color-mix(in srgb, ${color} 15%, transparent)` }}
          >
            {icon}
          </span>
          <h3 className="font-display text-base font-bold text-fg-heading">
            {title}
          </h3>
        </div>
        <Link
          to={linkTo}
          className="w-8 h-8 rounded-lg bg-surface-raised flex items-center justify-center
                     text-fg-muted hover:text-accent hover:bg-surface-hover transition-colors
                     group-hover:bg-accent/10"
          title={`Add ${title}`}
        >
          <span className="text-lg leading-none">+</span>
        </Link>
      </div>
      {children}
    </div>
  );
}

function MiniList({ items, empty }) {
  if (items.length === 0) {
    return <p className="text-xs text-fg-subtle">{empty}</p>;
  }
  return (
    <ul className="space-y-1.5">
      {items.map((item, i) => (
        <li key={i} className="text-xs text-fg-muted flex items-baseline gap-1.5">
          {item}
        </li>
      ))}
    </ul>
  );
}

export default function Dashboard() {
  const [feeds, setFeeds] = useState([]);
  const [sleep, setSleep] = useState([]);
  const [growth, setGrowth] = useState([]);
  const [diapers, setDiapers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const [f, s, g, d] = await Promise.all([
          getFeeds().catch(() => []),
          getSleep().catch(() => []),
          getGrowth().catch(() => []),
          getDiapers().catch(() => []),
        ]);
        setFeeds(f.slice(-5).reverse());
        setSleep(s.slice(-5).reverse());
        setGrowth(g.slice(-3).reverse());
        setDiapers(d.slice(-5).reverse());
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="w-8 h-8 border-3 border-accent/30 border-t-accent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-4 animate-slide-up">
      {/* Greeting */}
      <div className="mb-2">
        <h2 className="font-display text-2xl font-bold text-fg-heading">
          Dashboard
        </h2>
        <p className="text-sm text-fg-muted mt-1">
          Today&apos;s overview at a glance
        </p>
      </div>

      {/* Quick Entry Strip */}
      <div className="flex gap-2 overflow-x-auto pb-1">
        {[
          { to: "/feeds", icon: "🍼", label: "Feed", color: "var(--mod-feeds)" },
          { to: "/sleep", icon: "😴", label: "Sleep", color: "var(--mod-sleep)" },
          { to: "/growth", icon: "📏", label: "Growth", color: "var(--mod-growth)" },
          { to: "/susupoty", icon: "🧷", label: "Diaper", color: "var(--mod-diaper)" },
        ].map(({ to, icon, label, color }) => (
          <Link
            key={to}
            to={to}
            className="btn-quick flex items-center gap-2 shrink-0"
            style={{ borderColor: `color-mix(in srgb, ${color} 30%, transparent)` }}
          >
            <span>{icon}</span>
            <span>+ {label}</span>
          </Link>
        ))}
      </div>

      {/* Summary Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {/* Feeds */}
        <SummaryCard
          icon="🍼"
          title="Feeds"
          color="var(--mod-feeds)"
          linkTo="/feeds"
        >
          <MiniList
            empty="No feeds logged yet"
            items={feeds.map((f) => (
              <span key={f.id}>
                <span className="font-semibold text-fg">{f.type}</span>
                {" "}<span className="text-fg-subtle">{f.date}</span>
                {f.quantity > 0 && <span className="text-mod-feeds"> {f.quantity}ml</span>}
              </span>
            ))}
          />
        </SummaryCard>

        {/* Sleep */}
        <SummaryCard
          icon="😴"
          title="Sleep"
          color="var(--mod-sleep)"
          linkTo="/sleep"
        >
          <MiniList
            empty="No sleep logged yet"
            items={sleep.map((s) => (
              <span key={s.id}>
                <span className="font-semibold text-fg">{s.type}</span>
                {" "}<span className="text-fg-subtle">{s.date}</span>
                {s.quality && <span className="text-mod-sleep"> {s.quality}</span>}
              </span>
            ))}
          />
        </SummaryCard>

        {/* Growth */}
        <SummaryCard
          icon="📏"
          title="Growth"
          color="var(--mod-growth)"
          linkTo="/growth"
        >
          {growth.length === 0 ? (
            <p className="text-xs text-fg-subtle">No measurements yet</p>
          ) : (
            <div className="space-y-2">
              {/* Latest stats */}
              <div className="flex gap-4">
                {growth[0].weight > 0 && (
                  <div>
                    <p className="text-lg font-bold text-fg-heading">{growth[0].weight}</p>
                    <p className="text-[0.65rem] text-fg-subtle uppercase tracking-wider">kg</p>
                  </div>
                )}
                {growth[0].height > 0 && (
                  <div>
                    <p className="text-lg font-bold text-fg-heading">{growth[0].height}</p>
                    <p className="text-[0.65rem] text-fg-subtle uppercase tracking-wider">cm</p>
                  </div>
                )}
                {growth[0].head_circ > 0 && (
                  <div>
                    <p className="text-lg font-bold text-fg-heading">{growth[0].head_circ}</p>
                    <p className="text-[0.65rem] text-fg-subtle uppercase tracking-wider">head cm</p>
                  </div>
                )}
              </div>
              <p className="text-[0.65rem] text-fg-subtle">
                Last measured: {growth[0].date}
              </p>
            </div>
          )}
        </SummaryCard>

        {/* Diapers */}
        <SummaryCard
          icon="🧷"
          title="Diapers"
          color="var(--mod-diaper)"
          linkTo="/susupoty"
        >
          <MiniList
            empty="No changes logged yet"
            items={diapers.map((d) => (
              <span key={d.id}>
                <span className="font-semibold text-fg">{d.type}</span>
                {" "}<span className="text-fg-subtle">{d.date}</span>
                {d.notes && <span className="text-fg-subtle"> — {d.notes}</span>}
              </span>
            ))}
          />
        </SummaryCard>
      </div>
    </div>
  );
}
