import React, { useEffect, useState } from "react";
import { Button, Space, Typography, theme } from "antd";
import { CloseOutlined, DownloadOutlined } from "@ant-design/icons";

const { Text } = Typography;

interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed" }>;
}

const InstallPWABanner: React.FC = () => {
  const [installPrompt, setInstallPrompt] =
    useState<BeforeInstallPromptEvent | null>(null);
  const [showBanner, setShowBanner] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const { token } = theme.useToken();

  // Check if device is mobile
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);

    return () => {
      window.removeEventListener("resize", checkMobile);
    };
  }, []);

  useEffect(() => {
    const handler = (e: Event) => {
      e.preventDefault();
      setInstallPrompt(e as BeforeInstallPromptEvent);
      setShowBanner(true);
    };

    window.addEventListener("beforeinstallprompt", handler);

    return () => {
      window.removeEventListener("beforeinstallprompt", handler);
    };
  }, []);

  const handleInstallClick = async () => {
    if (!installPrompt) return;

    installPrompt.prompt();
    const { outcome } = await installPrompt.userChoice;

    if (outcome === "accepted") {
      setShowBanner(false);
    }
  };

  // Don't show banner if not mobile or banner shouldn't be shown
  if (!isMobile || !showBanner) return null;

  return (
    <div
      style={{
        position: "fixed",
        bottom: 0,
        left: 0,
        right: 0,
        backgroundColor: token.colorBgContainer,
        padding: "12px 16px",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        zIndex: 1000,
      }}
    >
      <Space>
        <DownloadOutlined style={{ fontSize: 20, color: "#1db954" }} />
        <div>
          <Text strong>Install G-hopper</Text>
          <br />
          <Text type="secondary">
            Add to home screen for the best experience
          </Text>
        </div>
      </Space>
      <Space>
        <Button type="primary" onClick={handleInstallClick}>
          Install
        </Button>
        <Button
          type="text"
          icon={<CloseOutlined />}
          onClick={() => setShowBanner(false)}
          aria-label="Close"
        />
      </Space>
    </div>
  );
};

export default InstallPWABanner;
