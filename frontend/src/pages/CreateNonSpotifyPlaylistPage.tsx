/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useState } from "react";
import {
  Button,
  Card,
  Form,
  Input,
  Select,
  Space,
  Typography,
  Divider,
  Alert,
  Steps,
  Result,
  App,
} from "antd";
import { MinusCircleOutlined, PlusOutlined } from "@ant-design/icons";
import { Content } from "antd/es/layout/layout";
import { config } from "../config";
import { Link } from "react-router-dom";
import { SeedTrack } from "../types/non-spotify";
import { generatePlaylist } from "../services/nonSpotifyPlaylistService";

const { Title, Paragraph } = Typography;
const { Option } = Select;

const genres = [
  { value: "hip-hop", label: "Hip-Hop / Rap / R&B" },
  { value: "electronic", label: "Electronic / Dance" },
  { value: "rock", label: "Rock / Pop" },
  { value: "soul", label: "Soul / Funk / Disco" },
  { value: "jazz", label: "Jazz / Blues" },
  { value: "reggae", label: "Reggae / Dub" },
  { value: "country", label: "Country / Folk" },
  { value: "world", label: "World / Latin" },
  { value: "soundtrack", label: "Soundtrack" },
  { value: "classical", label: "Classical" },
];

const CreateNonSpotifyPlaylistPage: React.FC = () => {
  const [form] = Form.useForm();
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [playlistId, setPlaylistId] = useState<string | null>(null);
  const [seedTracks, setSeedTracks] = useState<SeedTrack[]>([
    { title: "", artist: "" },
  ]);
  const { message: messageApi } = App.useApp();

  // Handle form submission
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const genre = values.genre;

      if (!seedTracks || seedTracks.length === 0) {
        setError("Please add at least one song");
        return;
      }

      // Filter out empty tracks
      const validSeeds = seedTracks.filter((seed) => seed.title && seed.artist);

      if (validSeeds.length === 0) {
        setError(
          "Please add at least one valid song with both title and artist"
        );
        return;
      }

      setLoading(true);

      const playlist = await generatePlaylist(validSeeds, genre);
      setPlaylistId(playlist.id);
      setSuccess(true);
      messageApi.success("Playlist created successfully");
    } catch (err: any) {
      console.error("Error creating playlist:", err);
      setError(err.response?.data?.error || "Failed to create playlist");
      messageApi.error("Failed to create playlist");
    } finally {
      setLoading(false);
    }
  };

  // Handle adding a new seed track
  const addSeedTrack = () => {
    if (seedTracks.length >= 20) {
      messageApi.warning("You can add at most 20 songs");
      return;
    }
    setSeedTracks([...seedTracks, { title: "", artist: "" }]);
  };

  // Handle removing a seed track
  const removeSeedTrack = (index: number) => {
    const newSeeds = [...seedTracks];
    newSeeds.splice(index, 1);
    setSeedTracks(newSeeds);
  };

  // Handle updating a seed track
  const updateSeedTrack = (
    index: number,
    field: "title" | "artist",
    value: string
  ) => {
    const newSeeds = [...seedTracks];
    newSeeds[index] = { ...newSeeds[index], [field]: value };
    setSeedTracks(newSeeds);
  };

  // Handle next step button
  const goToNextStep = () => {
    const validSeeds = seedTracks.filter((seed) => seed.title && seed.artist);

    if (validSeeds.length === 0) {
      messageApi.error(
        "Please add at least one song with both title and artist"
      );
      return;
    }

    setCurrentStep(1);
  };

  // Render steps content
  const stepsContent = [
    // Step 1: Add songs
    <>
      <Paragraph>
        Add up to 20 songs you like. We'll use these to find similar songs in
        your chosen genre.{" "}
        <span style={{ color: "#1db954" }}>
          The more songs you add, the better the results!
        </span>
      </Paragraph>

      {seedTracks.map((track, index) => (
        <Space
          key={index}
          style={{ display: "flex", marginBottom: 8, width: "100%" }}
          align="baseline"
        >
          <Input
            placeholder="Song Title"
            value={track.title}
            onChange={(e) => updateSeedTrack(index, "title", e.target.value)}
            style={{ flex: 1 }}
          />
          <Input
            placeholder="Artist Name"
            value={track.artist}
            onChange={(e) => updateSeedTrack(index, "artist", e.target.value)}
            style={{ flex: 1 }}
          />
          <Button
            type="text"
            icon={<MinusCircleOutlined />}
            onClick={() => removeSeedTrack(index)}
            disabled={seedTracks.length <= 1}
          />
        </Space>
      ))}

      <Button
        type="dashed"
        onClick={addSeedTrack}
        style={{ width: "100%", marginBottom: 16 }}
        icon={<PlusOutlined />}
        disabled={seedTracks.length >= 20}
      >
        Add Song
      </Button>

      {seedTracks.length >= 20 && (
        <Paragraph type="secondary" style={{ marginTop: 8 }}>
          You've reached the maximum of 20 songs.
        </Paragraph>
      )}

      <Button type="primary" onClick={goToNextStep}>
        Next
      </Button>
    </>,

    // Step 2: Select genre
    <>
      <Paragraph>
        Choose a genre you want to explore. We'll find songs that connect your
        favorites to this genre.
      </Paragraph>

      <Form.Item
        name="genre"
        rules={[{ required: true, message: "Please select a genre" }]}
      >
        <Select placeholder="Select a genre" style={{ width: "100%" }}>
          {genres.map((genre) => (
            <Option key={genre.value} value={genre.value}>
              {genre.label}
            </Option>
          ))}
        </Select>
      </Form.Item>

      <Space>
        <Button onClick={() => setCurrentStep(0)}>Previous</Button>
        <Button type="primary" onClick={handleSubmit} loading={loading}>
          Generate Playlist
        </Button>
      </Space>
    </>,
  ];

  if (success && playlistId) {
    return (
      <Content
        style={{
          padding: "24px",
          minHeight: `calc(100vh - ${config.headerHeight * 2}px)`,
        }}
      >
        <Card>
          <Result
            status="success"
            title="Playlist Created Successfully!"
            subTitle="Your playlist has been created and is ready to use."
            extra={[
              <Button type="primary" key="view">
                <Link to={`/non-spotify/playlists/${playlistId}`}>
                  View Playlist
                </Link>
              </Button>,
              <Button key="dashboard">
                <Link to="/non-spotify/dashboard">Return to Dashboard</Link>
              </Button>,
            ]}
          />
        </Card>
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
      <Card>
        <Title level={2}>Create New Playlist</Title>
        <Paragraph>
          Create a playlist by adding songs you like and selecting a genre to
          explore.
        </Paragraph>

        {error && (
          <Alert
            message="Error"
            description={error}
            type="error"
            showIcon
            style={{ marginBottom: 16 }}
          />
        )}

        <Steps
          current={currentStep}
          items={[{ title: "Add Songs" }, { title: "Select Genre" }]}
          style={{ marginBottom: 24 }}
        />

        <Divider />

        <Form form={form} layout="vertical">
          {stepsContent[currentStep]}
        </Form>
      </Card>
    </Content>
  );
};

export default CreateNonSpotifyPlaylistPage;
