import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import App from "./App";
import * as serviceWorkerRegistration from "./serviceWorkerRegistration";

const root = createRoot(document.getElementById("root"));
root.render(
  <StrictMode>
    <App />
  </StrictMode>,
);

// Register service worker for PWA (offline + installable)
serviceWorkerRegistration.register();
