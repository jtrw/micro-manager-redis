import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const PORT = process.env.PORT || 8080;
const DEVSERVER_BASE_PATH =
  process.env.DEVSERVER_BASE_PATH || `http://127.0.0.1:${PORT}`;

export default defineConfig(({ mode }) => {
  const isDev = mode === "development";
  const RKEYS_URL = isDev ? DEVSERVER_BASE_PATH : "{% MANAGE_RKEYS_URL %}";
  process.env.RKEYS_URL = RKEYS_URL;
  return {
    // build specific config
    plugins: [react()],
    define: {
      "process.env": process.env,
    },
    server: {
      host: true,
    },
    base: "./",
  };
});
