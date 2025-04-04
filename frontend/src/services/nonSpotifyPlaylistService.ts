import axios from "axios";
import {
    GeneratePlaylistRequest,
    NonSpotifyPlaylist,
    NonSpotifyPlaylistWithTracks,
    SeedTrack
} from "../types/non-spotify";
import { getAuthHeader } from "./nonSpotifyAuthService";

// Generate a new playlist
const generatePlaylist = async (
    seedTracks: SeedTrack[],
    genre: string
): Promise<NonSpotifyPlaylistWithTracks> => {
    const request: GeneratePlaylistRequest = {
        seed_tracks: seedTracks,
        genre
    };

    const response = await axios.post<{ playlist: NonSpotifyPlaylistWithTracks }>(
        "/api/api/non-spotify/playlists",
        request,
        { headers: getAuthHeader() }
    );

    return response.data.playlist;
};

// Get all playlists for the current user
const getUserPlaylists = async (): Promise<NonSpotifyPlaylist[]> => {
    const response = await axios.get<{ playlists: NonSpotifyPlaylist[] }>(
        "/api/api/non-spotify/playlists",
        { headers: getAuthHeader() }
    );

    return response.data.playlists;
};

// Get playlist details including tracks
const getPlaylistDetails = async (playlistId: string): Promise<NonSpotifyPlaylistWithTracks> => {
    const response = await axios.get<{ playlist: NonSpotifyPlaylistWithTracks }>(
        `/api/api/non-spotify/playlists/${playlistId}`,
        { headers: getAuthHeader() }
    );

    return response.data.playlist;
};

// Update track status (added to playlist)
const updateTrackStatus = async (
    trackId: string,
    addedToPlaylist: boolean
): Promise<void> => {
    await axios.patch(
        `/api/api/non-spotify/tracks/${trackId}`,
        { added_to_playlist: addedToPlaylist },
        { headers: getAuthHeader() }
    );
};

// Delete a playlist
const deletePlaylist = async (playlistId: string): Promise<void> => {
    await axios.delete(
        `/api/api/non-spotify/playlists/${playlistId}`,
        { headers: getAuthHeader() }
    );
};

export {
    generatePlaylist,
    getUserPlaylists,
    getPlaylistDetails,
    updateTrackStatus,
    deletePlaylist
};