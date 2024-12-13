import { Flex, Space, Spin } from "antd";
import { useAuth } from "../hooks/useAuth";
import Title from "antd/es/typography/Title";
import { TopArtists } from "../components/TopArtists";
import { TopTracks } from "../components/TopTracks";
import { GeneratedPlaylists } from "../components/GeneratedPlaylists";

const UserProfile = () => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <Flex style={{ textAlign: "center", padding: "50px" }}>
        <Spin />
      </Flex>
    );
  }

  return (
    <Space
      direction="vertical"
      size="large"
      style={{ width: "100%", padding: "24px" }}
    >
      <Title
        style={{
          textAlign: "center",
          marginBottom: "24px",
        }}
        level={1}
      >
        Welcome,
        <span style={{ color: "#1db954" }}> {user?.display_name}</span>
      </Title>

      <GeneratedPlaylists />
      <TopArtists />
      <TopTracks />
    </Space>
  );
};

export default UserProfile;
