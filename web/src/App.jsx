// BabyTracker/web/src/routes.js
import { BrowserRouter as Router } from "react-router-dom";

import MainLayout from "./Main";
import AppRoutes from "./Routes";

// App
export default function App() {
  return (
    <Router>
      <MainLayout>
        <AppRoutes />
      </MainLayout>
    </Router>
  );
}
