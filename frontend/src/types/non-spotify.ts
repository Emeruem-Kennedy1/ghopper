export type NonSpotifyUser = {
    id: string;
    created_at: string;
    updated_at: string;
};

export type NonSpotifyPlaylist = {
    id: string;
    user_id: string;
    name: string;
    genre: string;
    description: string;
    image_url: string;
    created_at: string;
    updated_at: string;
};

export type NonSpotifyPlaylistTrack = {
    id: string;
    playlist_id: string;
    title: string;
    artist: string;
    added_to_playlist: boolean;
    created_at: string;
    updated_at: string;
};

export type NonSpotifyPlaylistSeedTrack = {
    id: string;
    playlist_id: string;
    title: string;
    artist: string;
    created_at: string;
    updated_at: string;
};

export type NonSpotifyPlaylistWithTracks = NonSpotifyPlaylist & {
    tracks: NonSpotifyPlaylistTrack[];
    seed_tracks: NonSpotifyPlaylistSeedTrack[];
};

export type RegisterRequest = {
    user_id: string;
};

export type RegisterResponse = {
    user_id: string;
    passphrase: string;
};

export type VerifyRequest = {
    user_id: string;
    passphrase: string;
};

export type SeedTrack = {
    title: string;
    artist: string;
};

export type GeneratePlaylistRequest = {
    seed_tracks: SeedTrack[];
    genre: string;
};

export type UpdateTrackStatusRequest = {
    added_to_playlist: boolean;
};