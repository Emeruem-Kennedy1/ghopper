// components/layout/PrivacyFooter.tsx
import React, { useState } from "react";
import { Modal, Typography, Space, Layout, theme } from "antd";

const { Footer } = Layout;
const { Text, Title, Paragraph } = Typography;
const { useToken } = theme;

interface PrivacyFooterProps {
  spotify: boolean;
}

const PrivacyFooter: React.FC<PrivacyFooterProps> = ({ spotify }) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  const { token } = useToken();

  return (
    <Footer
      style={{
        textAlign: "center",
        padding: "1rem 10px",
        position: "fixed",
        bottom: 0,
        left: 0,
        right: 0,
        width: "100%",
        zIndex: 999,
        margin: 0,
        boxSizing: "border-box",
      }}
    >
      <Space>
        {spotify && (
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
        )}
        <Text>© {new Date().getFullYear()} G-hopper</Text>
        <Text type="secondary">|</Text>
        <Text
          className="privacy-link"
          onClick={showModal}
          style={{ cursor: "pointer", textDecoration: "underline" }}
        >
          Privacy Policy
        </Text>
      </Space>

      <Modal
        title="Privacy Policy"
        open={isModalOpen}
        onCancel={handleCancel}
        footer={null}
        width={700}
      >
        <Typography>
          <Title level={4}>G-hopper Privacy Policy</Title>
          <Paragraph>Last updated: February 20, 2025</Paragraph>

          <Title level={5}>Introduction</Title>
          <Paragraph>
            G-hopper ("we", "our", or "us") is committed to protecting your
            privacy. This Privacy Policy explains how we collect, use, and
            safeguard your information when you use our web application that
            helps you discover new music by creating playlists based on
            connections between songs.
          </Paragraph>

          <Title level={5}>Information We Collect</Title>
          <Paragraph>When you use G-hopper, we collect:</Paragraph>
          <ul>
            <li>
              <strong>Spotify Account Information:</strong> When you
              authenticate with Spotify, we receive your Spotify user ID,
              display name, email, country, profile image, and Spotify URI.
            </li>
            <li>
              <strong>Music Listening Data:</strong> We access your top 50
              tracks and artists from Spotify to analyze your listening
              preferences.
            </li>
            <li>
              <strong>Created Playlists:</strong> We store information about
              playlists created through our service.
            </li>
            <li>
              <strong>Usage Data:</strong> We collect information about how you
              interact with G-hopper, including genres explored and features
              used.
            </li>
          </ul>

          <Title level={5}>How We Use Your Information</Title>
          <Paragraph>We use your information to:</Paragraph>
          <ul>
            <li>Create and manage your G-hopper account</li>
            <li>
              Analyze your music preferences to create personalized playlists
            </li>
            <li>Create and manage playlists in your Spotify account</li>
            <li>Enhance the user experience and develop new features</li>
            <li>Ensure the security and proper functioning of our service</li>
          </ul>

          <Title level={5}>Data Storage and Retention</Title>
          <Paragraph>
            Your data is stored securely in our database using industry-standard
            security practices. We retain your account information and created
            playlists for as long as you maintain an active G-hopper account.
            Your Spotify authentication tokens are securely stored and refreshed
            automatically when needed.
          </Paragraph>

          <Title level={5}>Account Deletion</Title>
          <Paragraph>
            You can delete your G-hopper account at any time through the account
            settings in the application. When you delete your account:
          </Paragraph>
          <ul>
            <li>
              All personal information associated with your account will be
              permanently removed from our database
            </li>
            <li>Your authentication tokens will be revoked and deleted</li>
            <li>Records of your top tracks and artists will be deleted</li>
            <li>
              Information about playlists created through G-hopper will be
              deleted from our database
            </li>
          </ul>
          <Paragraph>
            Please note: Playlists that were created in your Spotify account
            through G-hopper will remain in your Spotify account unless you
            manually delete them through Spotify. We cannot delete content
            directly from your Spotify account after it has been created.
          </Paragraph>

          <Title level={5}>Data Sharing and Third-Party Services</Title>
          <Paragraph>
            G-hopper integrates with Spotify to provide our service. When you
            use G-hopper:
          </Paragraph>
          <ul>
            <li>
              We request specific permissions from Spotify through their
              authorization system
            </li>
            <li>
              We create playlists in your Spotify account based on your
              selections
            </li>
            <li>
              We do not sell or share your personal information with third
              parties
            </li>
          </ul>
          <Paragraph>
            We use secure, industry-standard hosting and connection methods to
            protect your data. Our service providers do not have access to your
            personal data or Spotify information.
          </Paragraph>

          <Title level={5}>Your Rights</Title>
          <Paragraph>You have the right to:</Paragraph>
          <ul>
            <li>Access the personal information we hold about you</li>
            <li>Request correction of inaccurate personal information</li>
            <li>Request deletion of your account and associated data</li>
            <li>
              Revoke Spotify permissions through your Spotify account settings
            </li>
            <li>Obtain information about how we use your data</li>
          </ul>

          <Title level={5}>Security</Title>
          <Paragraph>
            We implement appropriate security measures to protect your personal
            information, including:
          </Paragraph>
          <ul>
            <li>Secure HTTPS connections</li>
            <li>Secure authentication methods</li>
            <li>Proper session management</li>
            <li>Regular security updates</li>
          </ul>

          <Title level={5}>Changes to This Privacy Policy</Title>
          <Paragraph>
            We may update our Privacy Policy from time to time. We will notify
            you of any changes by posting the new Privacy Policy on this page
            and updating the "Last updated" date.
          </Paragraph>

          <Title level={5}>Contact Us</Title>
          <Paragraph>
            If you have questions or concerns about this Privacy Policy or our
            data practices, please contact us at: kennedyemeruem@uni.minerva.edu
          </Paragraph>
        </Typography>
      </Modal>
    </Footer>
  );
};

export default PrivacyFooter;
