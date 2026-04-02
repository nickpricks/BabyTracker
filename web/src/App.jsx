import { BrowserRouter as Router } from "react-router-dom";

import MainLayout from "./Main";
import AppRoutes from "./Routes";
import ErrorBoundary from "./components/ErrorBoundary";
import { useTheme } from "./themes";

export default function App() {
  const theme = useTheme();

  return (
    <Router>
      <MainLayout theme={theme}>
        <ErrorBoundary>
          <AppRoutes />
        </ErrorBoundary>
      </MainLayout>
    </Router>
  );
}
