// frontend/src/components/layout/NonSpotifyHeader.tsx
import React from "react";
import { Layout, Avatar, Dropdown, Space, theme, Modal } from "antd";
import {
  LogoutOutlined,
  DashboardOutlined,
  PlusOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Link } from "react-router-dom";
import Logo from "../common/Logo";
import { useNonSpotifyAuth } from "../../hooks/useNonSpotifyAuth";
import { config } from "../../config";
import type { MenuProps } from "antd";

const { Header } = Layout;
const { useToken } = theme;

const NonSpotifyHeader: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const { token } = useToken();
  const { userId, logout } = useNonSpotifyAuth();

  const handleLogout = () => {
    Modal.confirm({
      title: "Logout",
      content: "Are you sure you want to logout?",
      okText: "Logout",
      cancelText: "Cancel",
      onOk: logout,
    });
  };

  const menuItems: MenuProps["items"] = [
    {
      key: "dashboard",
      icon: <DashboardOutlined />,
      label: <Link to="/non-spotify/dashboard">Dashboard</Link>,
    },
    {
      key: "create-playlist",
      icon: <PlusOutlined />,
      label: <Link to="/non-spotify/create-playlist">Create Playlist</Link>,
    },
    {
      type: "divider",
    },
    {
      key: "logout",
      icon: <LogoutOutlined />,
      label: "Logout",
      danger: true,
      onClick: handleLogout,
    },
  ];

  const isDarkMode = token.colorTextBase === token.colorWhite;

  return (
    <Header
      style={{
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        padding: "0 24px",
        background: token.colorBgContainer,
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        zIndex: 1000,
        height: config.headerHeight,
      }}
    >
      <Link to="/non-spotify/dashboard">
        <Logo themeColor={isDarkMode ? "#fafafa" : "transparent"} />
      </Link>

      <Space>
        {userId && (
          <Space>
            {children}
            <Dropdown menu={{ items: menuItems }} placement="bottomRight">
              <Space style={{ cursor: "pointer" }}>
                <Avatar
                  style={{ backgroundColor: token.colorPrimary }}
                  icon={<UserOutlined />}
                />
                <span>{userId}</span>
              </Space>
            </Dropdown>
          </Space>
        )}
        <span
          style={{
            padding: "4px 8px",
            backgroundColor: token.colorPrimaryBg,
            color: token.colorPrimary,
            borderRadius: "10px",
            fontSize: "12px",
          }}
        >
          General Mode
        </span>
      </Space>
    </Header>
  );
};

export default NonSpotifyHeader;
