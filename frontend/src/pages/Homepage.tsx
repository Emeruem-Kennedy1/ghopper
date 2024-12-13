import React, { useState } from "react";
import { Button, Typography, App } from "antd";
import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { useAuth } from "../hooks/useAuth";
import { getToken } from "../utils/auth";
import axios from "axios";
import spotifyGif from "../assets/spotify.gif";
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
  { id: "hip-hop", name: "Hip-Hop" },
  { id: "rap", name: "Rap" },
  { id: "r&b", name: "R&B" },
  { id: "electronic", name: "Electronic" },
  { id: "dance", name: "Dance" },
  { id: "rock", name: "Rock" },
  { id: "pop", name: "Pop" },
  { id: "soul", name: "Soul" },
  { id: "funk", name: "Funk" },
  { id: "disco", name: "Disco" },
  { id: "jazz", name: "Jazz" },
  { id: "blues", name: "Blues" },
  { id: "reggae", name: "Reggae" },
  { id: "dub", name: "Dub" },
  { id: "country", name: "Country" },
  { id: "folk", name: "Folk" },
  { id: "world", name: "World" },
  { id: "latin", name: "Latin" },
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
      const response = await axios.post(
        "/api/api/toptracks-analysis",
        {
          genre,
          userId,
        },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status !== 200) throw new Error(response.data.message);

      return response.data;
    },
    onSuccess: (data) => {
      console.log(data);
      setIsModalVisible(true);
    },
    onError: (error) =>
      messageApi.error(
        error instanceof Error
          ? error.message
          : "An error occurred. Please try again."
      ),
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
                src={spotifyGif}
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
