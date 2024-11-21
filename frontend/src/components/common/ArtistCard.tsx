import { Avatar, Tag, Typography } from "antd";
import { Artist } from "../../types";

const { Text } = Typography;

interface ArtistCardProps {
  artist: Artist;
}

export const ArtistCard = ({ artist }: ArtistCardProps) => (
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
      href={artist.uri}
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
        src={artist.image}
        size={140}
        style={{ border: "2px solid #e8e8e8" }}
      />
      <Text strong>{artist.name}</Text>
    </a> 
    <div
      style={{
        display: "flex",
        flexWrap: "wrap",
        justifyContent: "center",
        gap: "4px",
        maxWidth: "200px",
      }}
    >
      {artist.genres.slice(0, 1).map((genre) => (
        <Tag
          key={genre}
          style={{
            margin: 0,
            fontSize: "12px",
            borderRadius: "1rem",
          }}
        >
          {genre}
        </Tag>
      ))}
    </div>
  </div>
);
