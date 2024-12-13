import { useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createContext, ReactNode, useCallback } from "react";
import { getToken, removeToken, storeToken } from "../utils/auth";
import axios from "axios";
import { useNavigate, useLocation } from "react-router-dom";
import { AuthContextType, UserProfile } from "../types/auth";
import fetchUser from "../services/userService";

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined
);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const location = useLocation();

  const { data: user, isLoading } = useQuery({
    queryKey: ["user"],
    queryFn: fetchUser,
    retry: false,
    staleTime: Infinity,
  });

  useEffect(() => {
    if (getToken()) {
      queryClient.prefetchQuery({
        queryKey: ["user"],
        queryFn: fetchUser,
      });
    }
  }, [queryClient]);

  const login = useCallback(
    (userData: UserProfile, token: string) => {
      queryClient.setQueryData(["user"], userData);
      storeToken(token);
    },
    [queryClient]
  );

  const logout = useCallback(() => {
    removeToken();
    queryClient.setQueryData(["user"], null);
    navigate("/login");
  }, [queryClient, navigate]);

  // Axios interceptor for 401 errors
  useEffect(() => {
    const interceptor = axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Clear all queries and cached data
          queryClient.clear();
          // Handle the logout
          logout();
          // preserve the intended destination
          const destination = location.pathname;
          if (destination !== "/login") {
            navigate("/login", {
              state: {
                returnTo: destination,
                error: "Your session has expired. Please log in again.",
              },
            });
          }
        }
        return Promise.reject(error);
      }
    );

    // Cleanup interceptor on unmount
    return () => {
      axios.interceptors.response.eject(interceptor);
    };
  }, [logout, navigate, location, queryClient]);

  const handleAuthCallback = useCallback(() => {
    const params = new URLSearchParams(location.search);
    const encodedData = params.get("data");
    const error = params.get("error");

    if (error) {
      console.error("Authentication error:", error);
      navigate("/login", { state: { error } });
      return;
    }

    if (encodedData) {
      try {
        const decodedData = atob(encodedData);
        const data = JSON.parse(decodedData);

        if (data.user && data.token) {
          login(data.user, data.token);
          navigate("/dashboard");
        } else {
          throw new Error("Invalid data structure");
        }
      } catch (error) {
        console.error("Failed to process authentication data:", error);
        navigate("/login", {
          state: { error: "Failed to process authentication data" },
        });
      }
    } else {
      navigate("/login", {
        state: { error: "No authentication data received" },
      });
    }
  }, [location, login, navigate]);

  return (
    <AuthContext.Provider
      value={{
        user: user ?? null,
        login,
        logout,
        isLoading,
        handleAuthCallback,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
