import { useContext } from "react";
import { NonSpotifyAuthContext } from "../context/NonSpotifyAuthContext";


export const useNonSpotifyAuth = () => {
    const context = useContext(NonSpotifyAuthContext);

    if (context === undefined) {
        throw new Error(
            "useNonSpotifyAuth must be used within a NonSpotifyAuthProvider"
        );
    }

    return context;
};
