import axios from "axios";
import { Artist, TopArtistResponse } from "../types";
import { getToken } from "../utils/auth";

export async function getTopArtists(): Promise<Artist[]> {
  const token = getToken();
  const response = await axios.get("/api/api/user/top-artists", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  const artists = response.data.artists.items;

  return artists.map((artist: TopArtistResponse) => ({
    id: artist.id,
    name: artist.name,
    uri: artist.external_urls.spotify,
    image: artist.images[0].url,
    genres: artist.genres,
  }));
}
