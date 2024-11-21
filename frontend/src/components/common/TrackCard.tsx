import { Avatar, Typography } from "antd";
import { Track } from "../../types";

const { Text } = Typography;

interface TrackCardProps {
  track: Track;
}

export const TrackCard = ({ track }: TrackCardProps) => (
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
      href={`https://open.spotify.com/track/${track.id}`}
      target="_blank"
      rel="noopener noreferrer"
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "0.5rem",
      }}
    >
      <Avatar
        src={track.image}
        shape="square"
        size={140}
        style={{
          border: "2px solid #e8e8e8",
          borderRadius: "8px",
        }}
      />
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
          {track.name}
        </Text>
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
          {track.artists.map((artist) => artist).join(", ")}
        </Text>
      </div>
    </a>
  </div>
);
