// frontend/src/pages/NonSpotifyDashboardPage.tsx
import React, { useEffect, useState } from "react";
import {
  Button,
  Col,
  Empty,
  Row,
  Space,
  Spin,
  Typography,
  message,
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import { Link, useNavigate } from "react-router-dom";
import { useNonSpotifyAuth } from "../hooks/useNonSpotifyAuth";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";
import { NonSpotifyPlaylist } from "../types/non-spotify";
import { getUserPlaylists } from "../services/nonSpotifyPlaylistService";
import NonSpotifyPlaylistCard from "../components/common/NonSpotifyPlaylistCard";

const { Title } = Typography;

const NonSpotifyDashboardPage: React.FC = () => {
  const { userId } = useNonSpotifyAuth();
  const [playlists, setPlaylists] = useState<NonSpotifyPlaylist[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPlaylists = async () => {
      try {
        const data = await getUserPlaylists();
        setPlaylists(data);
        setError(null);
      } catch (err) {
        console.error(err);
        setError("Failed to fetch playlists");
        message.error("Failed to load playlists");
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylists();
  }, []);

  const handleDeletePlaylist = (playlistId: string) => {
    // Filter out the deleted playlist from state
    setPlaylists(playlists.filter((playlist) => playlist.id !== playlistId));
  };

  const renderContent = () => {
    if (loading) {
      return (
        <div style={{ textAlign: "center", padding: "50px" }}>
          <Spin size="large" />
        </div>
      );
    }

    if (error) {
      return (
        <Empty
          description={<span>Error loading playlists: {error}</span>}
          style={{ margin: "50px 0" }}
        >
          <Button type="primary" onClick={() => window.location.reload()}>
            Try Again
          </Button>
        </Empty>
      );
    }

    if (playlists.length === 0) {
      return (
        <Empty
          description={<span>No playlists found</span>}
          style={{ margin: "50px 0" }}
        >
          <Button type="primary">
            <Link to="/non-spotify/create-playlist">
              Create Your First Playlist
            </Link>
          </Button>
        </Empty>
      );
    }

    return (
      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        {playlists.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()).map((playlist) => (
          <Col xs={24} sm={12} md={8} lg={6} key={playlist.id}>
            <NonSpotifyPlaylistCard
              playlist={playlist}
              onDelete={handleDeletePlaylist}
            />
          </Col>
        ))}
      </Row>
    );
  };

  return (
    <Content
      style={{
        padding: "24px",
        minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
      }}
    >
      <Space direction="vertical" size="large" style={{ width: "100%" }}>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Title level={2}>Your Playlists</Title>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              navigate("/non-spotify/create-playlist");
            }}
          >
            Create New Playlist
          </Button>
        </div>

        {userId && (
          <Typography.Paragraph>
            Welcome back, <strong>{userId}</strong>
          </Typography.Paragraph>
        )}

        {renderContent()}
      </Space>
    </Content>
  );
};

export default NonSpotifyDashboardPage;
