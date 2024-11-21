import { Flex, Space, Spin } from "antd";
import { useAuth } from "../hooks/useAuth";
import Title from "antd/es/typography/Title";
import { TopArtists } from "../components/TopArtists";
import { TopTracks } from "../components/TopTracks";

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
      <Title level={2}>Welcome, {user?.display_name}</Title>
      <TopArtists />
      <TopTracks />
    </Space>
  );
};

export default UserProfile;
