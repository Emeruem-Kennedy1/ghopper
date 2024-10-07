import { SpotifyOutlined } from "@ant-design/icons";
import { Button, Col, Layout, Row } from "antd";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";

const LoginPage = () => {
  const location = useLocation();
  const error = location.state?.error;
  const { user } = useAuth();

  const handleSpotifyLogin = () => {
    window.location.href = `api/auth/spotify/login`;
  };

  if (user) {
    return <Navigate to="/" replace />;
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
          <Col>
            {error && <p style={{ color: "red" }}>{error}</p>}
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
