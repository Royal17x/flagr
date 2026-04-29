import { useNavigate, NavLink } from 'react-router-dom'
import { Flag, LayoutDashboard, LogOut, Shield } from 'lucide-react'
import { authApi } from '../api/auth'
import { authStore } from '../store/auth'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const navigate = useNavigate()
  const user = authStore.getUser()

  const handleLogout = async () => {
    const refreshToken = authStore.getRefreshToken()
    if (refreshToken) {
      await authApi.logout(refreshToken).catch(() => {})
    }
    authStore.clearTokens()
    navigate('/login')
  }

  const navBase =
    'group flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-400 transition-all duration-200 hover:text-white hover:bg-white/[0.04]'
  const navActive =
    'group flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-white bg-white/[0.06] shadow-[inset_0_1px_0_0_rgba(255,255,255,0.04)] relative'

  return (
    <div className="min-h-screen bg-slate-950 flex selection:bg-indigo-500/30">
      <aside className="w-64 bg-slate-950/80 backdrop-blur-xl border-r border-white/[0.06] flex flex-col fixed h-full z-20">
        <div className="p-6 pb-4">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/20">
              <Shield size={18} className="text-white" />
            </div>
            <div>
              <h1 className="text-lg font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent tracking-tight">
                Flagr
              </h1>
            </div>
          </div>
          <p className="text-[11px] font-medium text-slate-500 mt-1.5 uppercase tracking-wider">
            Feature Flags
          </p>
        </div>

        <nav className="flex-1 px-3 space-y-1 mt-2">
          <NavLink
            to="/dashboard"
            className={({ isActive }: { isActive: boolean }) =>
              isActive ? navActive : navBase
            }
          >
            <LayoutDashboard size={17} className="transition-transform group-hover:scale-110" />
            Dashboard
          </NavLink>

          <NavLink
            to="/flags"
            className={({ isActive }) => (isActive ? navActive : navBase)}
          >
            <Flag size={17} className="transition-transform group-hover:scale-110" />
            Flags
          </NavLink>
        </nav>

        <div className="p-4 mt-auto">
          {user && (
            <div className="mb-3 px-3 py-2 rounded-lg bg-white/[0.03] border border-white/[0.05]">
              <p className="text-[11px] text-slate-500 uppercase tracking-wider font-semibold">Organization</p>
              <p className="text-xs text-slate-300 font-mono mt-0.5 truncate">
                {user.orgId.slice(0, 12)}…
              </p>
            </div>
          )}
          <button
            onClick={handleLogout}
            className="group flex items-center gap-3 w-full px-3 py-2.5 rounded-xl text-sm font-medium text-slate-400 hover:text-red-400 hover:bg-red-500/10 transition-all duration-200"
          >
            <LogOut size={17} className="transition-transform group-hover:scale-110" />
            Sign out
          </button>
        </div>
      </aside>

      <main className="flex-1 overflow-auto ml-64 relative">
        <div className="absolute inset-0 pointer-events-none">
          <div className="absolute top-0 left-0 w-[500px] h-[500px] bg-indigo-500/[0.03] rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2" />
          <div className="absolute bottom-0 right-0 w-[600px] h-[600px] bg-purple-500/[0.02] rounded-full blur-3xl translate-x-1/3 translate-y-1/3" />
        </div>
        <div className="relative z-10">{children}</div>
      </main>
    </div>
  )
}