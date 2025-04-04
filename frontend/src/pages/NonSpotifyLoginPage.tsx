import React, { useState } from "react";
import {
  Alert,
  Button,
  Card,
  Divider,
  Form,
  Input,
  Typography,
  App,
} from "antd";
import { useNavigate, useLocation, Link } from "react-router-dom";
import { useNonSpotifyAuth } from "../hooks/useNonSpotifyAuth";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";

const { Title, Text, Paragraph } = Typography;

const NonSpotifyLoginPage: React.FC = () => {
  const [form] = Form.useForm();
  const { login } = useNonSpotifyAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { message: messageApi } = App.useApp();

  // Get the path to redirect to after login (if any)
  interface LocationState {
    from?: string;
  }

  const from =
    (location.state as LocationState)?.from || "/non-spotify/dashboard";

  const onFinish = async (values: { userId: string; passphrase: string }) => {
    setLoading(true);
    setError(null);

    try {
      const success = await login(values.userId, values.passphrase);

      if (success) {
        messageApi.success("Logged in successfully");
        navigate(from, { replace: true });
      } else {
        setError("Invalid user ID or passphrase");
      }
    } catch (err) {
      console.error(err);
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Content
      style={{
        minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        padding: "24px",
      }}
    >
      <Card style={{ width: 400, maxWidth: "100%" }}>
        <Title level={2} style={{ textAlign: "center" }}>
          Welcome to G-hopper
        </Title>
        <Paragraph style={{ textAlign: "center" }}>
          Log in to discover new music and create playlists
        </Paragraph>

        {error && (
          <Alert
            message="Login Failed"
            description={error}
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
          />
        )}

        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="userId"
            label="User ID"
            rules={[{ required: true, message: "Please enter your user ID" }]}
          >
            <Input placeholder="Enter your user ID" />
          </Form.Item>

          <Form.Item
            name="passphrase"
            label="Passphrase"
            rules={[
              { required: true, message: "Please enter your passphrase" },
            ]}
          >
            <Input.Password placeholder="Enter your passphrase" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              Log In
            </Button>
          </Form.Item>
        </Form>

        <Divider>Or</Divider>

        <div style={{ textAlign: "center" }}>
          <Text>Don't have an account?</Text>
          <br />
          <Button type="link" style={{ padding: 0 }}>
            <Link to="/non-spotify/register">Register now</Link>
          </Button>
        </div>

        <div style={{ marginTop: 16, textAlign: "center" }}>
          <Text type="secondary">Have a Spotify account?</Text>
          <br />
          <Button type="link" style={{ padding: 0 }}>
            <Link to="/login">Login with Spotify</Link>
          </Button>
        </div>
      </Card>
    </Content>
  );
};

export default NonSpotifyLoginPage;
