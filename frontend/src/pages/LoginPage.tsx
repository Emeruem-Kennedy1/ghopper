import { SpotifyOutlined } from "@ant-design/icons";
import { Alert, Button, Col, Layout, Row } from "antd";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";

interface LocationState {
  error?: string;
  returnTo?: string;
}

const LoginPage = () => {
  const location = useLocation();
  const { error, returnTo } = (location.state as LocationState) || {};
  const { user } = useAuth();

  const handleSpotifyLogin = () => {
    // Store the return path in localStorage before redirecting
    if (returnTo) {
      localStorage.setItem("returnTo", returnTo);
    }
    window.location.href = `/api/auth/spotify/login`;
  };

  if (user) {
    // Check if there's a stored return path
    const storedReturnTo = localStorage.getItem("returnTo");
    // Clean up the stored path
    localStorage.removeItem("returnTo");
    // Redirect to the stored path or home
    return <Navigate to={storedReturnTo || "/"} replace />;
  }

  return (
    <Layout>
      <Content
        style={{
          minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Row>
          <Col style={{ width: "100%", maxWidth: 400 }}>
            {error && (
              <Alert
                message="Authentication Error"
                description={error}
                type="error"
                showIcon
                style={{ marginBottom: 16 }}
              />
            )}
            <Button
              icon={<SpotifyOutlined />}
              type="primary"
              onClick={handleSpotifyLogin}
              size="large"
              block
            >
              Login with Spotify
            </Button>
          </Col>
        </Row>
      </Content>
    </Layout>
  );
};

export default LoginPage;
