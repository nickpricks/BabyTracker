// BabyTracker/web/src/routes.js
import { BrowserRouter as Router } from "react-router-dom";

import MainLayout from "./Main";
import AppRoutes from "./Routes";
import ErrorBoundary from "./components/ErrorBoundary";

// App
export default function App() {
  return (
    <Router>
      <MainLayout>
        <ErrorBoundary>
          <AppRoutes />
        </ErrorBoundary>
      </MainLayout>
    </Router>
  );
}
