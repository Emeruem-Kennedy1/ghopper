import React, { useState } from "react";
import {
  Alert,
  Button,
  Card,
  Divider,
  Form,
  Input,
  Typography,
  message,
  Result,
  Space,
} from "antd";
import { Link } from "react-router-dom";
import { useNonSpotifyAuth } from "../hooks/useNonSpotifyAuth";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";
import { CopyOutlined } from "@ant-design/icons";

const { Title, Text, Paragraph } = Typography;

const NonSpotifyRegisterPage: React.FC = () => {
  const [form] = Form.useForm();
  const { register } = useNonSpotifyAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [registered, setRegistered] = useState(false);
  const [passphrase, setPassphrase] = useState<string>("");
  const [userId, setUserId] = useState<string>("");

  const onFinish = async (values: { userId: string }) => {
    setLoading(true);
    setError(null);

    try {
      const response = await register(values.userId);
      setUserId(response.userId);
      setPassphrase(response.passphrase);
      setRegistered(true);
      message.success("Registration successful");
    } catch (err: unknown) {
      if (
        typeof err === "object" &&
        err !== null &&
        "response" in err &&
        (err as { response?: { status?: number } }).response?.status === 409
      ) {
        setError("This user ID already exists. Please try another one.");
      } else {
        setError("Registration failed. Please try again.");
      }
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    message.success("Copied to clipboard");
  };

  if (registered) {
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
        <Card style={{ width: 500, maxWidth: "100%" }}>
          <Result
            status="success"
            title="Registration Successful!"
            subTitle="Please save your passphrase in a secure location. You will need it to log in."
          />

          <div style={{ padding: "0 24px", marginBottom: 24 }}>
            <div style={{ marginBottom: 16 }}>
              <Text strong>User ID:</Text>
              <div
                style={{ display: "flex", alignItems: "center", marginTop: 8 }}
              >
                <Input value={userId} readOnly />
                <Button
                  icon={<CopyOutlined />}
                  onClick={() => copyToClipboard(userId)}
                  style={{ marginLeft: 8 }}
                />
              </div>
            </div>

            <div style={{ marginBottom: 24 }}>
              <Text strong>Passphrase:</Text>
              <div
                style={{ display: "flex", alignItems: "center", marginTop: 8 }}
              >
                <Input value={passphrase} readOnly />
                <Button
                  icon={<CopyOutlined />}
                  onClick={() => copyToClipboard(passphrase)}
                  style={{ marginLeft: 8 }}
                />
              </div>
            </div>

            <Alert
              message="Important"
              description="Please write down your passphrase or save it somewhere secure. You won't be able to recover it if lost."
              type="warning"
              showIcon
              style={{ marginBottom: 24 }}
            />

            <Space direction="vertical" style={{ width: "100%" }}>
              <Button type="primary" block>
                <Link to="/non-spotify/dashboard">Continue to Dashboard</Link>
              </Button>
              <Button block>
                <Link to="/non-spotify/login">Go to Login</Link>
              </Button>
            </Space>
          </div>
        </Card>
      </Content>
    );
  }

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
          Register New Account
        </Title>
        <Paragraph style={{ textAlign: "center" }}>
          Create a new account to discover music and build playlists
        </Paragraph>

        {error && (
          <Alert
            message="Registration Failed"
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
            tooltip="Choose a unique identifier for your account"
            rules={[
              { required: true, message: "Please enter a user ID" },
              { min: 3, message: "User ID must be at least 3 characters" },
              { max: 30, message: "User ID must be less than 30 characters" },
            ]}
          >
            <Input placeholder="Enter a unique user ID" />
          </Form.Item>

          <Paragraph type="secondary" style={{ marginBottom: 16 }}>
            A secure passphrase will be generated for you after registration.
          </Paragraph>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              Register
            </Button>
          </Form.Item>
        </Form>

        <Divider>Or</Divider>

        <div style={{ textAlign: "center" }}>
          <Text>Already have an account?</Text>
          <br />
          <Button type="link" style={{ padding: 0 }}>
            <Link to="/non-spotify/login">Log in now</Link>
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

export default NonSpotifyRegisterPage;
