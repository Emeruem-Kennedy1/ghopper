// frontend/src/context/NonSpotifyAuthContext.tsx
import { createContext, ReactNode, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  getNonSpotifyCredentials,
  isNonSpotifyUserLoggedIn,
  logoutNonSpotifyUser,
  registerNonSpotifyUser,
  verifyNonSpotifyUser,
} from "../services/nonSpotifyAuthService";

type NonSpotifyAuthContextType = {
  isLoggedIn: boolean;
  userId: string | null;
  login: (userId: string, passphrase: string) => Promise<boolean>;
  register: (userId: string) => Promise<{ userId: string; passphrase: string }>;
  logout: () => void;
  isLoading: boolean;
};

export const NonSpotifyAuthContext = createContext<
  NonSpotifyAuthContextType | undefined
>(undefined);

export const NonSpotifyAuthProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [userId, setUserId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  // Check if user is already logged in on mount
  useEffect(() => {
    const checkLoginStatus = () => {
      const loggedIn = isNonSpotifyUserLoggedIn();
      setIsLoggedIn(loggedIn);

      if (loggedIn) {
        const credentials = getNonSpotifyCredentials();
        setUserId(credentials?.userId || null);
      }

      setIsLoading(false);
    };

    checkLoginStatus();
  }, []);

  const login = async (
    userId: string,
    passphrase: string
  ): Promise<boolean> => {
    try {
      const success = await verifyNonSpotifyUser(userId, passphrase);

      if (success) {
        setIsLoggedIn(true);
        setUserId(userId);
        return true;
      }

      return false;
    } catch (error) {
      console.error("Login failed:", error);
      return false;
    }
  };

  const register = async (userId: string) => {
    const response = await registerNonSpotifyUser(userId);
    const formattedResponse = {
      userId: response.user_id,
      passphrase: response.passphrase,
    };
    setIsLoggedIn(true);
    setUserId(formattedResponse.userId);
    return formattedResponse;
  };

  const logout = () => {
    logoutNonSpotifyUser();
    setIsLoggedIn(false);
    setUserId(null);
    navigate("/non-spotify/login");
  };

  return (
    <NonSpotifyAuthContext.Provider
      value={{
        isLoggedIn,
        userId,
        login,
        register,
        logout,
        isLoading,
      }}
    >
      {children}
    </NonSpotifyAuthContext.Provider>
  );
};
