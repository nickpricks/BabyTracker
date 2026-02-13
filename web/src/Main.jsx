import { Link } from "react-router-dom";

export default function MainLayout({ children }) {
  return (
    <div style={{ maxWidth: 800, margin: "0 auto", padding: 24 }}>
      <header>
        <h1 style={{ fontSize: 20, marginBottom: 12 }}>Baby Tracker</h1>
        <nav style={{ marginBottom: 24, display: "flex", gap: 16 }}>
          <Link to="/feeds">Feeds</Link>
          <Link to="/sleep">Sleep</Link>
          <Link to="/growth">Growth</Link>
          <Link to="/susupoty">Susu-Poty</Link>
        </nav>
      </header>
      <main>{children}</main>
    </div>
  );
}
