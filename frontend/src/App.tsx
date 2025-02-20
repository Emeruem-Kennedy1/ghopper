import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter } from "react-router-dom";
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
import { Content } from "antd/es/layout/layout";
import MainLayout from "./components/layout/MainLayout";
import {App as AntApp} from "antd";
import PrivacyFooter from "./components/layout/PrivacyFooter";

const queryClient = new QueryClient();

const App = () => {
  const [currentTheme, setCurrentTheme] = useState(() => {
    return localStorage.getItem("theme") || "default";
  });

  useEffect(() => {
    // Update data-theme attribute when theme changes
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
  return (
    <AntApp>
      <ConfigProvider theme={getTheme()}>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <AuthProvider>
              <MainLayout>
                <CustomHeader>
                  <Switch
                    checked={currentTheme === "dark"}
                    onChange={toggleTheme}
                    checkedChildren="🌙"
                    unCheckedChildren="☀️"
                  />
                </CustomHeader>
                <Content>
                  <AppRoutes />
                </Content>
                <PrivacyFooter />
              </MainLayout>
            </AuthProvider>
          </BrowserRouter>
        </QueryClientProvider>
      </ConfigProvider>
    </AntApp>
  );
};

export default App;
