import React, { useState } from "react";
import { Button, Typography, App } from "antd";
import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { useAuth } from "../hooks/useAuth";
import { getToken } from "../utils/auth";
import axios from "axios";
import loadingGif from "../assets/loading-animation.gif";
import { SearchResultModal } from "../components/common/SearchResultModal";

const { Title, Paragraph } = Typography;

const containerStyle = {
  display: "flex",
  flexDirection: "column" as const,
  alignItems: "center",
  justifyContent: "center",
  minHeight: "calc(100vh - 124px)",
  padding: "2rem",
};

const loadingContainerStyle = {
  position: "fixed" as const,
  top: 0,
  left: 0,
  width: "100vw",
  height: "100vh",
  backgroundColor: "rgba(0, 0, 0, 0.7)",
  display: "flex",
  flexDirection: "column" as const,
  justifyContent: "center",
  alignItems: "center",
  zIndex: 1000,
};

const loadingContentStyle = {
  textAlign: "center" as const,
  color: "white",
};

const contentStyle = {
  maxWidth: "1024px",
  textAlign: "center" as const,
};

const titleStyle = {
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  gap: "4px",
};

const highlightStyle = {
  color: "#1db954",
};

const buttonContainerStyle = {
  display: "flex",
  flexWrap: "wrap" as const,
  justifyContent: "center",
  gap: "1rem",
  marginBottom: "2rem",
};

const genres = [
  { id: "hip-hop", name: "Hip-Hop / Rap / R&B" },
  { id: "electronic", name: "Electronic / Dance" },
  { id: "rock", name: "Rock / Pop" },
  { id: "soul", name: "Soul / Funk / Disco" },
  { id: "jazz", name: "Jazz / Blues" },
  { id: "reggae", name: "Reggae / Dub" },
  { id: "country", name: "Country / Folk" },
  { id: "world", name: "World / Latin" },
  { id: "soundtrack", name: "Soundtrack" },
  { id: "classical", name: "Classical" },
];

interface PlaylistResponse {
  playlist: string;
  songs: string[];
  message: string;
}

interface GeneratePlaylistParams {
  genre: string;
  userId: string;
}

const Homepage: React.FC = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const token = getToken();
  const [isModalVisible, setIsModalVisible] = useState(false);
  const { message: messageApi } = App.useApp();

  const generatePlaylistMutation = useMutation<
    PlaylistResponse,
    Error,
    GeneratePlaylistParams
  >({
    mutationFn: async ({ genre, userId }) => {
      try {
        const response = await axios.post(
          "/api/api/toptracks-analysis",
          { genre, userId },
          {
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
          }
        );
        return response.data;
      } catch (error) {
        if (axios.isAxiosError(error)) {
          throw new Error(error.response?.data?.error || error.message);
        }
        throw error;
      }
    },
    onSuccess: () => {
      setIsModalVisible(true);
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : "";

      switch (errorMessage) {
        case "Unauthorized":
          messageApi.error("Please log in again to continue");
          navigate("/login");
          break;
        case "no songs found for user":
          messageApi.info(
            "We couldn't find any recent listening history. Try playing some songs first!"
          );
          break;
        case "failed to get user's top tracks":
          messageApi.error(
            "Couldn't access your Spotify history. Please try again later."
          );
          break;
        case "failed to create playlist":
          messageApi.error(
            "Failed to create playlist. Please try again later."
          );
          break;
        case "genre is required":
          messageApi.error("Please select a genre to continue");
          break;
        default:
          messageApi.error(
            errorMessage || "An unexpected error occurred. Please try again."
          );
      }
    },
  });

  const handleGenreClick = (genreId: string) => {
    if (!user) {
      messageApi.info("Please sign in to generate playlists");
      navigate("/login");
      return;
    }
    generatePlaylistMutation.mutate({ genre: genreId, userId: user.id });
  };
  return (
    <div style={containerStyle}>
      {generatePlaylistMutation.isPending && (
        <div style={loadingContainerStyle}>
          <div style={loadingContentStyle}>
            <div>
              <img
                src={loadingGif}
                alt="Loading..."
                style={{
                  width: "150px",
                  height: "150px",
                }}
              />
            </div>
            <h2 style={{ marginTop: "1rem", color: "white" }}>
              Generating Your Playlist...
            </h2>
            <p>We're finding the perfect songs to match your taste</p>
          </div>
        </div>
      )}
      <div style={contentStyle}>
        <div style={titleStyle}>
          <Title level={1}>
            Welcome to G<span style={highlightStyle}>hopper</span> (Beta)
          </Title>
        </div>

        <Paragraph style={{ fontSize: "1.125rem", marginBottom: "2rem" }}>
          Break free from your music bubble with Ghopper! Streaming algorithms
          often trap us in comfortable but limiting musical loops. Ghopper uses
          music sampling relationships - where elements of existing songs appear
          in new works - to build bridges between the music you love and
          unexplored genres. By tracking how artists have sampled and reimagined
          each other's work, we create personalized playlists that connect your
          current favorites to exciting new territories. Choose a genre below,
          and we'll craft a playlist that uses familiar elements from your
          listening history to ease your journey into fresh musical landscapes.
        </Paragraph>

        <Title level={3} style={{ marginBottom: "1.5rem" }}>
          Choose a Genre to Explore
        </Title>

        <div style={buttonContainerStyle}>
          {genres.map((genre) => (
            <Button
              key={genre.id}
              type="primary"
              size="large"
              onClick={() => handleGenreClick(genre.id)}
              disabled={generatePlaylistMutation.isPending}
            >
              {genre.name}
            </Button>
          ))}
        </div>

        <SearchResultModal
          isModalVisible={isModalVisible}
          setIsModalVisible={setIsModalVisible}
          generatePlaylistMutation={generatePlaylistMutation}
        />
      </div>
    </div>
  );
};

export default Homepage;
