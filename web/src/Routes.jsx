import { Routes, Route } from "react-router-dom";

import Dashboard from "./components/Dashboard";
import { Feeds, Growth, Sleep, SusuPoty } from "./components/";

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />
      <Route path="/feeds" element={<Feeds />} />
      <Route path="/sleep" element={<Sleep />} />
      <Route path="/growth" element={<Growth />} />
      <Route path="/susupoty" element={<SusuPoty />} />
      <Route
        path="*"
        element={
          <div className="card text-center py-12">
            <p className="text-4xl mb-4">🍼</p>
            <h2 className="font-display text-xl font-bold text-fg-heading mb-2">
              404 — Not Found
            </h2>
            <p className="text-fg-muted">
              This page wandered off during naptime.
            </p>
          </div>
        }
      />
    </Routes>
  );
}
