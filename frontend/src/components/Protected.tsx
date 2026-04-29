import { Navigate } from 'react-router-dom'
import { authStore } from '../store/auth'

export default function Protected({ children }: { children: React.ReactNode }) {
  if (!authStore.isAuthenticated()) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}