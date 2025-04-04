import React from "react";
import { Layout } from "antd";
import NonSpotifyHeader from "./NonSpotifyHeader";
import { config } from "../../config";
import PrivacyFooter from "./PrivacyFooter";

const { Content } = Layout;

const NonSpotifyLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <Layout>
      <NonSpotifyHeader />
      <Content
        style={{
          marginTop: config.headerHeight,
          minHeight: `calc(100vh - ${config.headerHeight}px)`,
          padding: "24px 24px",
        }}
      >
        {children}
      </Content>
      <PrivacyFooter />
    </Layout>
  );
};

export default NonSpotifyLayout;
