import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useNonSpotifyAuth } from "../../hooks/useNonSpotifyAuth";
import { Flex, Spin } from "antd";

interface NonSpotifyProtectedRouteProps {
  children: React.ReactNode;
}

const NonSpotifyProtectedRoute: React.FC<NonSpotifyProtectedRouteProps> = ({
  children,
}) => {
  const { isLoggedIn, isLoading } = useNonSpotifyAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <Flex justify="center" align="center" style={{ height: "100vh" }}>
        <Spin size="large" />
      </Flex>
    );
  }

  if (!isLoggedIn) {
    // Save the current path to redirect back after login
    return (
      <Navigate
        to="/non-spotify/login"
        state={{ from: location.pathname }}
        replace
      />
    );
  }

  return <>{children}</>;
};

export default NonSpotifyProtectedRoute;
