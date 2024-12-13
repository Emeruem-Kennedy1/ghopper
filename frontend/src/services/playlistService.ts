import axios from "axios";
import { PlaylistResponse } from "../types";
import { getToken } from "../utils/auth";

async function getPlaylists(): Promise<PlaylistResponse[]> {
  const token = getToken();

  const response = await axios.get("/api/api/user/playlists", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  const playlists = response.data.playlists as PlaylistResponse[];

  if (!playlists) {
    return [];
  }

  return playlists.map((playlist: PlaylistResponse) => ({
    id: playlist.id,
    name: playlist.name,
    description: playlist.description,
    url: playlist.url,
    image: playlist.image,
  }));
}

const deletePlaylist = async (playlistId: string) => {
    const token = getToken();
    
    await axios.delete(`/api/api/user/playlists/${playlistId}`, {
        headers: {
        Authorization: `Bearer ${token}`,
        },
    });
}

export { getPlaylists, deletePlaylist };