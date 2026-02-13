// BabyTracker/web/src/routes.js
import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";

import { Feeds, Growth, Sleep, SusuPoty } from "./components/";

// App routes
export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/feeds" replace />} />
      <Route path="/feeds" element={<Feeds />} />
      <Route path="/sleep" element={<Sleep />} />
      <Route path="/growth" element={<Growth />} />
      <Route path="/susupoty" element={<SusuPoty />} />
      <Route path="*" element={<div>404 Not Found</div>} />
    </Routes>
  );
}
