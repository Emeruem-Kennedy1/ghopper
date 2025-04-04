import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./assets/styles/App.css";
import { AuthProvider } from "./context/AuthProvider";
import AppRoutes from "./routes/AppRoutes";
import { ConfigProvider, Switch } from "antd";
import { defaultTheme } from "./theme/defaultTheme";
import lightTheme from "./theme/lightTheme";
import darkTheme from "./theme/darkTheme";
import "antd/dist/reset.css";
import CustomHeader from "./components/layout/CustomHeader";
import { useEffect, useState } from "react";
import MainLayout from "./components/layout/MainLayout";
import NonSpotifyLayout from "./components/layout/NonSpotifyLayout";
import { App as AntApp } from "antd";
import { NonSpotifyAuthProvider } from "./context/NonSpotifyAuthContext";
import NonSpotifyHeader from "./components/layout/NonSpotifyHeader";

const queryClient = new QueryClient();

const App = () => {
  const [currentTheme, setCurrentTheme] = useState(() => {
    return localStorage.getItem("theme") || "default";
  });

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", currentTheme);
  }, [currentTheme]);

  useEffect(() => {
    localStorage.setItem("theme", currentTheme);
  }, [currentTheme]);

  const toggleTheme = (checked: boolean) => {
    setCurrentTheme(checked ? "dark" : "light");
  };

  const getTheme = () => {
    switch (currentTheme) {
      case "light":
        return lightTheme;
      case "dark":
        return darkTheme;
      default:
        return defaultTheme;
    }
  };

  // Theme toggle component to be used in both headers
  const ThemeToggle = (
    <Switch
      checked={currentTheme === "dark"}
      onChange={toggleTheme}
      checkedChildren="🌙"
      unCheckedChildren="☀️"
    />
  );

  return (
    <AntApp>
      <ConfigProvider theme={getTheme()}>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <Routes>
              {/* Non-Spotify routes with NonSpotifyLayout */}
              <Route
                path="/non-spotify/*"
                element={
                  <NonSpotifyAuthProvider>
                    <NonSpotifyLayout>
                      <NonSpotifyHeader>{ThemeToggle}</NonSpotifyHeader>
                      <AppRoutes />
                    </NonSpotifyLayout>
                  </NonSpotifyAuthProvider>
                }
              />

              {/* Regular routes with MainLayout */}
              <Route
                path="*"
                element={
                  <AuthProvider>
                    <MainLayout>
                      <CustomHeader>{ThemeToggle}</CustomHeader>
                      <AppRoutes />
                    </MainLayout>
                  </AuthProvider>
                }
              />
            </Routes>
          </BrowserRouter>
        </QueryClientProvider>
      </ConfigProvider>
    </AntApp>
  );
};

export default App;
