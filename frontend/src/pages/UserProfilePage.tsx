import { useAuth } from "../hooks/useAuth";

const UserProfile = () => {
  const { user, isLoading, logout } = useAuth();

  if (isLoading) return <div>Loading profile...</div>;

  return (
    <div>
      <h2>Welcome, {user?.display_name}</h2>
      <p>Email: {user?.email}</p>
      <p>Country: {user?.country}</p>
      <p>Spotify URI: {user?.uri}</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
};

export default UserProfile;
