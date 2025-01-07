import React from "react";
import { Layout, Avatar, Dropdown, Space, theme } from "antd";
import {
  LogoutOutlined,
  DashboardOutlined,
  NodeIndexOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Link } from "react-router-dom";
import Logo from "../common/Logo";
import { useAuth } from "../../hooks/useAuth";
import { config } from "../../config";

const { Header } = Layout;
const { useToken } = theme;

const CustomHeader: React.FC<{ children?: React.ReactNode }> = ({
  children,
}) => {
  const { token } = useToken();
  const { user, logout } = useAuth();

  const menuItems = [
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
