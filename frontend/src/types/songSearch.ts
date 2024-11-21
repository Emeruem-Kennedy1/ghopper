interface Artist {
  id: number;
  name: string;
  isMain: boolean;
}

interface Song {
  id: number;
  title: string;
  artists: Artist[];
  genres: string[];
}

export interface SongSearchResponseData {
  adjacencyList: {
    [key: string]: {
      [key: string]: Record<string, never>;
    };
  };
  nodes: {
    [key: string]: Song;
  };
  paths: Array<{
    start: string;
    end: string;
    pathNodes: string[];
    distance: number;
  }>;
}


export interface TransformedSong {
  id: number;
  title: string;
  mainArtist: {
    id: number;
    name: string;
  };
  featuredArtists: Array<{
    id: number;
    name: string;
  }>;
  genres: string[];
}


interface Artist {
  id: number;
  name: string;
  isMain: boolean;
}

interface Song {
  id: number;
  title: string;
  artists: Artist[];
  genres: string[];
}

export interface GraphData {
  nodes: Array<{
    id: string;
    data: {
      label: string;
    };
    position: {
      x: number;
      y: number;
    };
    type?: string;
  }>;
  edges: Array<{
    id: string;
    source: string;
    target: string;
  }>;
}
