import { Routes, Route } from 'react-router-dom'
import LoginPage from '../pages/LoginPage'
import AuthCallback from '../components/auth/AuthCallback'
import UserProfile from '../pages/UserProfilePage'
import ProtectedRoute from '../components/auth/ProtectedRoute'

const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<h1>Welcome to the G-Hopper App</h1>} />
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
    </Routes>
  )
}

export default AppRoutes