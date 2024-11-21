import axios from "axios";
import { TopTrackResponse, Track } from "../types";
import { getToken } from "../utils/auth";

export async function getTopTracks(): Promise<Track[]> {
  const token = getToken();

  const response = await axios.get("/api/api/user/top-tracks", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  const tracks = response.data.tracks;

  return tracks.map((track: TopTrackResponse) => ({
    id: track.id,
    name: track.name,
    image: track.image,
    artists: track.artists,
  }));
}
