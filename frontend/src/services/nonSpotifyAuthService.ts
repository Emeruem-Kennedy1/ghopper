import axios from "axios";
import { RegisterRequest, RegisterResponse, VerifyRequest } from "../types/non-spotify";

// Store the auth credentials in localStorage
const storeNonSpotifyCredentials = (userId: string, passphrase: string) => {
    localStorage.setItem('non_spotify_user_id', userId);
    localStorage.setItem('non_spotify_passphrase', passphrase);
};

// Remove the auth credentials from localStorage
const removeNonSpotifyCredentials = () => {
    localStorage.removeItem('non_spotify_user_id');
    localStorage.removeItem('non_spotify_passphrase');
};

// Get the auth credentials from localStorage
const getNonSpotifyCredentials = () => {
    const userId = localStorage.getItem('non_spotify_user_id');
    const passphrase = localStorage.getItem('non_spotify_passphrase');

    if (userId && passphrase) {
        return { userId, passphrase };
    }

    return null;
};

// Create authorization header for API calls
const getAuthHeader = () => {
    const credentials = getNonSpotifyCredentials();

    if (credentials) {
        const base64Credentials = btoa(`${credentials.userId}:${credentials.passphrase}`);
        return { Authorization: `Basic ${base64Credentials}` };
    }

    return {};
};

// Register a new non-Spotify user
const registerNonSpotifyUser = async (userId: string): Promise<RegisterResponse> => {
    const request: RegisterRequest = { user_id: userId };

    const response = await axios.post<RegisterResponse>(
        "/api/auth/non-spotify/register",
        request
    );

    // Store the credentials on successful registration
    storeNonSpotifyCredentials(response.data.user_id, response.data.passphrase);

    return response.data;
};

// Verify non-Spotify user credentials
const verifyNonSpotifyUser = async (userId: string, passphrase: string): Promise<boolean> => {
    const request: VerifyRequest = { user_id: userId, passphrase };

    try {
        await axios.post("/api/auth/non-spotify/verify", request);

        // Store the credentials on successful verification
        storeNonSpotifyCredentials(userId, passphrase);

        return true;
    } catch (error) {
        // Handle error if verification fails
        console.error("Verification failed:", error);
        return false;
    }
};

// Check if there is a logged in non-Spotify user
const isNonSpotifyUserLoggedIn = (): boolean => {
    return getNonSpotifyCredentials() !== null;
};

// Logout non-Spotify user
const logoutNonSpotifyUser = () => {
    removeNonSpotifyCredentials();
};

export {
    registerNonSpotifyUser,
    verifyNonSpotifyUser,
    isNonSpotifyUserLoggedIn,
    logoutNonSpotifyUser,
    getAuthHeader,
    getNonSpotifyCredentials
};