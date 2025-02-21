import React from "react";
import { Layout, Avatar, Dropdown, Space, theme, Modal } from "antd";
import {
  LogoutOutlined,
  DashboardOutlined,
  NodeIndexOutlined,
  DeleteOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Link } from "react-router-dom";
import Logo from "../common/Logo";
import { useAuth } from "../../hooks/useAuth";
import { config } from "../../config";
import type { MenuProps } from "antd";

const { Header } = Layout;
const { useToken } = theme;

const CustomHeader: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const { token } = useToken();
  const { user, logout, deleteAccount } = useAuth();

  const handleDeleteAccount = () => {
    Modal.confirm({
      title: "Delete Account",
      content:
        "Are you sure you want to delete your account? This action cannot be undone and will delete all your data including created playlists.",
      okText: "Delete",
      okType: "danger",
      cancelText: "Cancel",
      onOk: deleteAccount,
    });
  };

  const menuItems: MenuProps['items'] = [
    {
      key: "dashboard",
      icon: <DashboardOutlined />,
      label: <Link to="/dashboard">Dashboard</Link>,
    },
    {
      key: "analysis",
      icon: <NodeIndexOutlined />,
      label: <Link to="/analysis">Analysis</Link>,
    },
    {
      type: "divider",
    },
    {
      key: "delete-account",
      icon: <DeleteOutlined />,
      label: "Delete Account",
      danger: true,
      onClick: handleDeleteAccount,
    },
    {
      key: "logout",
      icon: <LogoutOutlined />,
      label: "Logout",
      onClick: () => {
        logout();
      },
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
      <Link to="/">
        <Logo themeColor={isDarkMode ? "#fafafa" : "transparent"} />
      </Link>

      <Space>
        {children}
        {user && (
          <Dropdown menu={{ items: menuItems }} placement="bottomLeft">
            <Avatar
              style={{ cursor: "pointer" }}
              src={user.image}
              icon={!user.image && <UserOutlined />}
            />
          </Dropdown>
        )}
      </Space>
    </Header>
  );
};

export default CustomHeader;
