import { Avatar, Typography, Button } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import { PlaylistResponse } from "../../types";

const { Text } = Typography;

interface PlaylistCardProps {
  playlist: PlaylistResponse;
  onDelete: (id: string) => void;
}

export const PlaylistCard = ({ playlist, onDelete }: PlaylistCardProps) => (
  <div
    style={{
      padding: "1rem",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      gap: "0.75rem",
    }}
  >
    <a
      href={playlist.url}
      target="_blank"
      rel="noopener noreferrer"
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "0.5rem",
      }}
    >
      <div style={{ position: "relative" }}>
        <Avatar
          src={playlist.image}
          shape="square"
          size={140}
          style={{
            border: "2px solid #e8e8e8",
            borderRadius: "8px",
          }}
        />
        <Button
          type="text"
          icon={<DeleteOutlined />}
          onClick={(e) => {
            e.preventDefault(); // Prevent link navigation when clicking delete
            onDelete(playlist.id);
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
        {playlist.description && (
          <Text
            type="secondary"
            style={{
              fontSize: "12px",
              overflow: "hidden",
              textOverflow: "ellipsis",
              display: "-webkit-box",
              WebkitLineClamp: 1,
              WebkitBoxOrient: "vertical",
            }}
          >
            {playlist.description}
          </Text>
        )}
      </div>
    </a>
  </div>
);
