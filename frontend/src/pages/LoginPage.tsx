// LoginPage.tsx
import { useLocation } from 'react-router-dom';

const LoginPage = () => {
  const location = useLocation();
  const error = location.state?.error;

  const handleSpotifyLogin = () => {
    window.location.href = `api/auth/spotify/login`;
  };

  return (
    <div>
      <h1>Login</h1>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <button onClick={handleSpotifyLogin}>Login with Spotify</button>
    </div>
  );
};

export default LoginPage;