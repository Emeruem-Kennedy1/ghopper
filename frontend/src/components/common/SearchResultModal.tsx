import { Modal, Button } from "antd";
import { useNavigate } from "react-router-dom";
import brokenLink from "../../assets/broken-link.svg";

export const SearchResultModal = ({
  isModalVisible,
  setIsModalVisible,
  generatePlaylistMutation,
}: {
  isModalVisible: boolean;
  setIsModalVisible: (value: boolean) => void;
  generatePlaylistMutation: {
    data?: {
      songs: string[];
      playlist?: string;
    };
  };
}) => {
  const navigate = useNavigate();
  return (
    <Modal
      title={
        generatePlaylistMutation.data?.songs.length
          ? "Playlist Generated!"
          : "No Matches Found"
      }
      open={isModalVisible}
      onOk={() => setIsModalVisible(false)}
      onCancel={() => setIsModalVisible(false)}
      footer={
        <div style={{ display: "flex", justifyContent: "center", gap: "8px" }}>
          <Button
            key="dashboard"
            type="primary"
            onClick={() => navigate("/dashboard")}
          >
            View on Dashboard
          </Button>
          <Button key="close" onClick={() => setIsModalVisible(false)}>
            Close
          </Button>
        </div>
      }
      style={{ textAlign: "center" }}
    >
      <div style={{ textAlign: "center" }}>
        {generatePlaylistMutation.data?.songs.length ? (
          <>
            <p>Your personalized playlist has been generated successfully!</p>
            {generatePlaylistMutation.data?.playlist && (
              <>
                <p>You can access your playlist here:</p>
                  <a
                    href={generatePlaylistMutation.data.playlist}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{
                      display: "inline-block",
                      padding: "8px 16px",
                      backgroundColor: "#1db954",
                      color: "white",
                      borderRadius: "20px",
                      textDecoration: "none",
                      marginBottom: "1rem",
                    }}
                  >
                    Open Playlist in Spotify
                  </a>
              </>
            )}
            <p style={{ marginTop: "1rem" }}>
              Visit your dashboard to see this and other playlists, along with
              insights about your musical preferences and the songs that
              inspired these recommendations.
            </p>
          </>
        ) : (
          <div>
            <div style={{ marginBottom: "1.5rem" }}>
              <img
                src={brokenLink}
                alt="Broken Link"
                style={{ width: "100px", height: "100px" }}
              />
            </div>
            <p style={{ fontSize: "1.1rem", marginBottom: "1rem" }}>
              Sorry! We couldn't find enough connections in your listening
              history to create a personalized playlist for this genre.
            </p>
            {generatePlaylistMutation.data?.playlist && (
              <>
                <p style={{ marginBottom: "1rem" }}>
                  But don't worry! Here's a curated playlist you might enjoy!
                </p>
                <a
                  href={generatePlaylistMutation.data.playlist}
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{
                    display: "inline-block",
                    padding: "8px 16px",
                    backgroundColor: "#1db954",
                    color: "white",
                    borderRadius: "20px",
                    textDecoration: "none",
                    marginBottom: "1rem",
                  }}
                >
                  Open Playlist in Spotify
                </a>
              </>
            )}
            <p style={{ color: "#666", fontSize: "0.9rem" }}>
              Try exploring a different genre or check out your dashboard for
              more recommendations!
            </p>
          </div>
        )}
      </div>
    </Modal>
  );
};
