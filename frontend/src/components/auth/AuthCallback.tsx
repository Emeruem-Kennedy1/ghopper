
// AuthCallback.tsx
import { useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';

const AuthCallback = () => {
  const { handleAuthCallback } = useAuth();

  useEffect(() => {
    handleAuthCallback();
  }, [handleAuthCallback]);

  return <div>Authenticating...</div>;
};

export default AuthCallback;