import React from "react";
import { Layout } from "antd";
import CustomHeader from "./CustomHeader";
import { config } from "../../config";

const { Content } = Layout;

const MainLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <Layout>
      <CustomHeader />
      <Content
        style={{
          marginTop: config.headerHeight,
          minHeight: `calc(100vh - ${config.headerHeight}px)`,
          padding: "24px 24px",
        }}
      >
        {children}
      </Content>
    </Layout>
  );
};

export default MainLayout;