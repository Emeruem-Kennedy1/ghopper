import React, { useEffect, useState } from "react";
import { Typography, Button, Popconfirm, Modal, Avatar } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import { Link } from "react-router-dom";
import { NonSpotifyPlaylist } from "../../types/non-spotify";
import { deletePlaylist } from "../../services/nonSpotifyPlaylistService";

// Import all cover images dynamically
// For Webpack
const coverImages = import.meta.glob("../../assets/covers/*.{jpg,jpeg,png}", {
  eager: true,
});

const { Text } = Typography;

interface NonSpotifyPlaylistCardProps {
  playlist: NonSpotifyPlaylist;
  onDelete: (playlistId: string) => void;
}

const NonSpotifyPlaylistCard: React.FC<NonSpotifyPlaylistCardProps> = ({
  playlist,
  onDelete,
}) => {
  const [coverImage, setCoverImage] = useState<string>("/default-playlist.jpg");

  useEffect(() => {
    // Determine image path based on playlist.image_url
    if (playlist.image_url) {
      // Look for the image in our dynamic imports
      const imagePath = Object.keys(coverImages).find((path) =>
        path.includes(`/${playlist.image_url}`)
      );

      if (imagePath) {
        const image = coverImages[imagePath] as { default?: string } | string;
        setCoverImage(typeof image === "string" ? image : image.default || "/default-playlist.jpg");
      } else {
        setCoverImage("/default-playlist.jpg");
      }
    }
  }, [playlist.image_url]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleDelete = async (e: React.MouseEvent) => {
    // Stop event propagation to prevent Link navigation
    e.stopPropagation();

    try {
      await deletePlaylist(playlist.id);
      onDelete(playlist.id);
    } catch (err) {
      console.error(err);
      Modal.error({
        title: "Delete Failed",
        content: "Failed to delete playlist. Please try again later.",
      });
    }
  };

  return (
    <Link
      to={`/non-spotify/playlists/${playlist.id}`}
      style={{
        padding: "1rem",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "0.75rem",
        textDecoration: "none",
        color: "inherit",
      }}
    >
      <div style={{ position: "relative" }}>
        <Avatar
          src={coverImage}
          shape="square"
          size={140}
          style={{
            border: "2px solid #e8e8e8",
            borderRadius: "8px",
          }}
        />
        <Popconfirm
          title="Delete this playlist?"
          description="This action cannot be undone."
          okText="Yes"
          cancelText="No"
          onConfirm={(e) => handleDelete(e as React.MouseEvent)}
          onCancel={(e) => e?.stopPropagation()}
        >
          <Button
            type="text"
            icon={<DeleteOutlined />}
            onClick={(e) => {
              e.preventDefault();
            }}
            style={{
              position: "absolute",
              top: "4px",
              right: "4px",
              color: "#ff4d4f",
              background: "rgba(255, 255, 255, 0.8)",
              padding: "4px 8px",
              height: "auto",
              borderRadius: "4px",
            }}
            aria-label="Delete playlist"
          />
        </Popconfirm>
      </div>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          textAlign: "center",
          maxWidth: "200px",
        }}
      >
        <Text
          strong
          style={{
            fontSize: "14px",
            marginBottom: "4px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            display: "-webkit-box",
            WebkitLineClamp: 2,
            WebkitBoxOrient: "vertical",
          }}
        >
          {playlist.name}
        </Text>
        <Text
          type="secondary"
          style={{
            fontSize: "12px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            display: "-webkit-box",
            WebkitLineClamp: 2,
            WebkitBoxOrient: "vertical",
            marginBottom: "4px",
          }}
        >
          {playlist.description || `Playlist for genre: ${playlist.genre}`}
        </Text>
        <Text
          type="secondary"
          style={{
            fontSize: "11px",
          }}
        >
          Created: {formatDate(playlist.created_at)}
        </Text>
      </div>
    </Link>
  );
};

export default NonSpotifyPlaylistCard;
