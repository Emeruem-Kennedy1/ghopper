import { Form, Input, Button, Select, Space, theme } from "antd";
import { MusicGenre } from "../types/enums";
import { getSongs } from "../services/songSearchService";
import { ArtistSearchData } from "../types";
// import { Node, Edge } from "@xyflow/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

interface FormValues {
  songName: string;
  artistName: string;
  genre: MusicGenre;
}

const { Item } = Form;
export const AnalysisForm = () => {
  const [form] = Form.useForm<FormValues>();
  const { token } = theme.useToken();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (data: ArtistSearchData) => getSongs(data),
    onSuccess: (data) => {
      // Invalidate and refetch
      queryClient.setQueryData(["graphData"], data);
    },
  });

  const onFinish = (values: FormValues) => {
    const data: ArtistSearchData = {
      songs: [{ title: values.songName, artist: values.artistName }],
      genre: values.genre,
      maxDepth: 2,
    };

    mutation.mutate(data);
  };

  return (
    <Form
      form={form}
      onFinish={onFinish}
      layout="vertical"
      style={{ width: "100%" }}
    >
      <Space.Compact
        block
        style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}
      >
        <Item
          label="Song Name"
          name="songName"
          rules={[{ required: true, message: "Please input the song name!" }]}
          style={{ flex: 1, minWidth: "200px" }}
        >
          <Input
            placeholder="Enter song name"
            style={{
              background: token.colorBgContainer,
              borderColor: token.colorBorder,
            }}
          />
        </Item>

        <Item
          label="Artist Name"
          name="artistName"
          rules={[{ required: true, message: "Please input the artist name!" }]}
          style={{ flex: 1, minWidth: "200px" }}
        >
          <Input
            placeholder="Enter artist name"
            style={{
              background: token.colorBgContainer,
              borderColor: token.colorBorder,
            }}
          />
        </Item>

        <Item
          label="Genre to Explore"
          name="genre"
          rules={[{ required: true, message: "Please select a genre!" }]}
          style={{ flex: 1, minWidth: "200px" }}
        >
          <Select
            placeholder="Select genre"
            style={{
              background: token.colorBgContainer,
            }}
          >
            {Object.values(MusicGenre).map((genre) => (
              <Select.Option key={genre} value={genre}>
                {genre}
              </Select.Option>
            ))}
          </Select>
        </Item>

        <Item
          label=" "
          style={{
            minWidth: "100px",
            display: "flex",
            alignItems: "flex-end",
          }}
        >
          <Button
            type="primary"
            htmlType="submit"
            style={{ width: "100%" }}
            loading={mutation.isPending}
          >
            Analyze
          </Button>
        </Item>
      </Space.Compact>
    </Form>
  );
};
