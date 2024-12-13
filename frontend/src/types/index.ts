export type Artist = {
  id: string;
  name: string;
  uri: string;
  image: string;
  genres: string[];
};

export type TopArtistResponse = {
  id: string;
  name: string;
  external_urls: { spotify: string };
  images: { url: string }[];
  genres: string[];
};

export type Track = {
  id: string;
  name: string;
  image: string;
  artists: Artist[];
};

export type TopTrackResponse = {
  id: string;
  name: string;
  image: string;
  artists: string[];
};

export type ArtistSearchData = {
  songs: { title: string; artist: string }[];
  genre: string;
  maxDepth: number;
};

export type PlaylistResponse = {
  id: string;
  name: string;
  description: string;
  url: string;
  image: string;
};