import { Routes, Route } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import AuthCallback from "../components/auth/AuthCallback";
import UserProfile from "../pages/UserProfilePage";
import ProtectedRoute from "../components/auth/ProtectedRoute";
import { AnalysisPage } from "../pages/AnalysisPage";
import NotFoundPage from "../pages/NotFoundPage";
import HomePage from "../pages/Homepage";
import NonSpotifyLoginPage from "../pages/NonSpotifyLoginPage";
import NonSpotifyProtectedRoute from "../components/auth/NonSpotifyProtectedRoute";
import CreateNonSpotifyPlaylistPage from "../pages/CreateNonSpotifyPlaylistPage";
import NonSpotifyDashboardPage from "../pages/NonSpotifyDashboardPage";
import NonSpotifyPlaylistDetailsPage from "../pages/NonSpotifyPlaylistDetailsPage";
import NonSpotifyRegisterPage from "../pages/NonSpotifyRegisterPage";
import { useLocation } from "react-router-dom";

const AppRoutes = () => {
  const location = useLocation();
  const isNonSpotifyRoute = location.pathname.startsWith("/non-spotify");

  // Non-Spotify routes
  if (isNonSpotifyRoute) {
    return (
      <Routes>
        <Route path="/login" element={<NonSpotifyLoginPage />} />
        <Route path="/register" element={<NonSpotifyRegisterPage />} />
        <Route
          path="/dashboard"
          element={
            <NonSpotifyProtectedRoute>
              <NonSpotifyDashboardPage />
            </NonSpotifyProtectedRoute>
          }
        />
        <Route
          path="/create-playlist"
          element={
            <NonSpotifyProtectedRoute>
              <CreateNonSpotifyPlaylistPage />
            </NonSpotifyProtectedRoute>
          }
        />
        <Route
          path="/playlists/:playlistId"
          element={
            <NonSpotifyProtectedRoute>
              <NonSpotifyPlaylistDetailsPage />
            </NonSpotifyProtectedRoute>
          }
        />
        <Route path="*" element={<NonSpotifyLoginPage />} />
      </Routes>
    );
  }

  // Regular routes
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/auth-callback" element={<AuthCallback />} />
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <UserProfile />
          </ProtectedRoute>
        }
      />
      <Route
        path="/analysis"
        element={
          <ProtectedRoute>
            <AnalysisPage />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

export default AppRoutes;
