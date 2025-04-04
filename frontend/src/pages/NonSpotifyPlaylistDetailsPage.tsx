import React, { useEffect, useState } from "react";
import {
  Alert,
  Avatar,
  Button,
  Card,
  Checkbox,
  Col,
  Divider,
  List,
  Row,
  Spin,
  Tag,
  Typography,
  message,
} from "antd";
import { ArrowLeftOutlined, SoundOutlined } from "@ant-design/icons";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";
import { Link, useNavigate, useParams } from "react-router-dom";
import { NonSpotifyPlaylistWithTracks } from "../types/non-spotify";
import {
  getPlaylistDetails,
  updateTrackStatus,
} from "../services/nonSpotifyPlaylistService";

// Import all cover images dynamically
const coverImages = import.meta.glob("../assets/covers/*.{jpg,jpeg,png}", {
  eager: true,
});

const { Title, Text, Paragraph } = Typography;

const NonSpotifyPlaylistDetailsPage: React.FC = () => {
  const { playlistId } = useParams<{ playlistId: string }>();
  const [playlist, setPlaylist] = useState<NonSpotifyPlaylistWithTracks | null>(
    null
  );
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [coverImage, setCoverImage] = useState<string>("/default-playlist.jpg");
  const navigate = useNavigate();

  useEffect(() => {
    const fetchPlaylistDetails = async () => {
      if (!playlistId) return;

      try {
        const data = await getPlaylistDetails(playlistId);
        setPlaylist(data);
        setError(null);
      } catch (err) {
        console.error(err);
        setError("Failed to fetch playlist details");
        message.error("Failed to load playlist");
      } finally {
        setLoading(false);
      }
    };

    fetchPlaylistDetails();
  }, [playlistId]);

  // Load cover image when playlist changes
  useEffect(() => {
    if (playlist?.image_url) {
      // Look for the image in our dynamic imports
      const imagePath = Object.keys(coverImages).find((path) =>
        path.includes(`/${playlist.image_url}`)
      );

      if (imagePath) {
        const image = coverImages[imagePath] as { default?: string } | string;
        setCoverImage(
          typeof image === "string"
            ? image
            : image.default || "/default-playlist.jpg"
        );
      } else {
        setCoverImage("/default-playlist.jpg");
      }
    }
  }, [playlist]);

  const handleCheckboxChange = async (trackId: string, checked: boolean) => {
    try {
      await updateTrackStatus(trackId, checked);

      // Update local state
      if (playlist) {
        setPlaylist({
          ...playlist,
          tracks: playlist.tracks.map((track) =>
            track.id === trackId
              ? { ...track, added_to_playlist: checked }
              : track
          ),
        });
      }
    } catch (err) {
      console.error(err);
      message.error("Failed to update track status");
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // Calculate progress
  const getTotalTracks = () => playlist?.tracks.length || 0;
  const getAddedTracks = () =>
    playlist?.tracks.filter((t) => t.added_to_playlist).length || 0;
  const getProgressPercentage = () => {
    const total = getTotalTracks();
    return total > 0 ? Math.round((getAddedTracks() / total) * 100) : 0;
  };

  if (loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Spin size="large" />
      </div>
    );
  }

  if (error || !playlist) {
    return (
      <Content
        style={{
          padding: "24px",
          minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
        }}
      >
        <Alert
          message="Error"
          description={error || "Playlist not found"}
          type="error"
          showIcon
          action={
            <Button type="primary">
              <Link to="/non-spotify/dashboard">Return to Dashboard</Link>
            </Button>
          }
        />
      </Content>
    );
  }

  return (
    <Content
      style={{
        padding: "24px",
        minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
      }}
    >
      <Button
        type="link"
        icon={<ArrowLeftOutlined />}
        style={{ marginBottom: 16, padding: 0 }}
        onClick={() => navigate("/non-spotify/dashboard")}
      >
        Back to Playlists
      </Button>

      <Card>
        <Row gutter={[24, 24]}>
          <Col xs={24} md={8}>
            <div style={{ marginBottom: 16 }}>
              <Avatar
                src={coverImage}
                icon={!coverImage && <SoundOutlined />}
                size={400}
                shape="square"
                style={{ width: "100%", height: "auto", maxWidth: "400px" }}
              />
            </div>

            <Title level={3}>{playlist.name}</Title>
            <Paragraph>{playlist.description}</Paragraph>

            <div style={{ marginBottom: 16 }}>
              <Tag
                style={{
                  backgroundColor: "#f0f0f0",
                  color: "#000",
                }}
              >
                {playlist.genre}
              </Tag>
              <Text type="secondary" style={{ display: "block", marginTop: 8 }}>
                Created: {formatDate(playlist.created_at)}
              </Text>
            </div>

            <div style={{ marginBottom: 16 }}>
              <Text strong>Progress:</Text>
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 8,
                  marginTop: 4,
                }}
              >
                <div
                  style={{
                    width: "100%",
                    height: 8,
                    backgroundColor: "#e0e0e0",
                    borderRadius: 4,
                    overflow: "hidden",
                  }}
                >
                  <div
                    style={{
                      width: `${getProgressPercentage()}%`,
                      height: "100%",
                      backgroundColor: "#1db954",
                      borderRadius: 4,
                      transition: "width 0.3s ease",
                    }}
                  />
                </div>
                <Text style={{ marginLeft: 8, fontSize: 12, width: 40 }}>
                  {getAddedTracks()}/{getTotalTracks()}
                </Text>
              </div>
            </div>

            <div>
              <Title level={4}>Seed Songs</Title>
              <Paragraph>
                These are the songs you used to generate this playlist.
              </Paragraph>

              <List
                dataSource={playlist.seed_tracks}
                renderItem={(track) => (
                  <List.Item>
                    <Text strong>{track.title}</Text> -{" "}
                    <Text>{track.artist}</Text>
                  </List.Item>
                )}
              />
            </div>
          </Col>

          <Col xs={24} md={16}>
            <div>
              <Title level={4}>Tracks</Title>
              <Paragraph>
                Check the box once you've added a song to your personal
                playlist.
              </Paragraph>

              <List
                dataSource={playlist.tracks}
                renderItem={(track) => (
                  <List.Item>
                    <div
                      style={{
                        display: "flex",
                        width: "100%",
                        alignItems: "center",
                      }}
                    >
                      <Checkbox
                        checked={track.added_to_playlist}
                        onChange={(e) =>
                          handleCheckboxChange(track.id, e.target.checked)
                        }
                        style={{ marginRight: 16 }}
                      />
                      <div style={{ flex: 1 }}>
                        <Text strong>{track.title}</Text>
                        <br />
                        <Text type="secondary">{track.artist}</Text>
                      </div>
                    </div>
                  </List.Item>
                )}
              />
            </div>

            <Divider />
          </Col>
        </Row>
      </Card>
    </Content>
  );
};

export default NonSpotifyPlaylistDetailsPage;
