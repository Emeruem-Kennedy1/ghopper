import {
  Alert,
  App,
  Button,
  Col,
  Empty,
  Grid,
  message,
  Row,
  Skeleton,
} from "antd";
import { PlaylistCard } from "./common/PlaylistCard";
import { deletePlaylist, getPlaylists } from "../services/playlistService";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Typography } from "antd";
import { useEffect, useState } from "react";

export const GeneratedPlaylists = () => {
  const queryClient = useQueryClient();
  const { message: messageApi } = App.useApp();
  const [showAll, setShowAll] = useState(false);
  const [itemsToShow, setItemsToShow] = useState(12);
  const { useBreakpoint } = Grid;
  const screens = useBreakpoint();

  useEffect(() => {
    if (screens.xl) setItemsToShow(6);
    else if (screens.lg) setItemsToShow(4);
    else if (screens.md) setItemsToShow(3);
    else if (screens.sm) setItemsToShow(2);
    else setItemsToShow(1);
  }, [screens]);

  const mutation = useMutation({
    mutationFn: deletePlaylist,
    onSuccess: (_, playlistId) => {
      // Optimistically update the UI
      queryClient.setQueryData(["playlists"], (old: { id: string }[]) =>
        old?.filter((playlist) => playlist.id !== playlistId)
      );
      messageApi.success("Playlist deleted successfully");
    },
    onError: (error) => {
      message.error(
        error instanceof Error
          ? error.message
          : "An error occurred. Please try again."
      );
    },
  });

  const {
    data: playlists,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["playlists"],
    queryFn: getPlaylists,
  });

  const handleDelete = (playlistId: string) => {
    mutation.mutate(playlistId);
  };

  if (isLoading)
    return (
      <div style={{ textAlign: "center", padding: "20px" }}>
        <Skeleton active />
      </div>
    );
  if (isError)
    return (
      <Alert
        type="error"
        message="Error"
        description="Failed to load playlists. Please try again later."
        showIcon
      />
    );

  return (
    <>
      <Typography.Title style={{ textAlign: "center" }} level={2}>
        Generated Playlists
      </Typography.Title>

      {!playlists || playlists.length === 0 ? (
        <Empty
          description={<span>No playlists yet</span>}
          style={{ margin: "20px 0" }}
        >
          <Button type="primary" href="/">
            Go to Tracks Analysis
          </Button>
        </Empty>
      ) : (
        <>
          <Row
            gutter={[16, 16]}
            style={{ display: "flex", justifyContent: "center" }}
          >
            {playlists
              .slice(0, showAll ? playlists.length : itemsToShow)
              .map((playlist) => (
                <Col xs={24} sm={12} md={8} lg={6} xl={4} key={playlist.id}>
                  <PlaylistCard playlist={playlist} onDelete={handleDelete} />
                </Col>
              ))}
          </Row>

          {playlists.length > itemsToShow && (
            <div style={{ textAlign: "center", marginTop: "20px" }}>
              <Button type="primary" onClick={() => setShowAll(!showAll)}>
                {showAll
                  ? "Show Less"
                  : `Show More (${playlists.length - itemsToShow} more)`}
              </Button>
            </div>
          )}
        </>
      )}
    </>
  );
};
