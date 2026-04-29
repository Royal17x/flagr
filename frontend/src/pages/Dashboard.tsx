import { useQuery } from '@tanstack/react-query'
import { Flag, Activity, Layers, } from 'lucide-react'
import Layout from '../components/Layout'
import { flagsApi } from '../api/flags'
import { authStore } from '../store/auth'
import { useProject } from '../hooks/useProject'

export default function Dashboard() {
  const user = authStore.getUser()
  const project = useProject()
  const projectId = project?.id ?? ''

  const { data: flags } = useQuery({
    queryKey: ['flags', projectId],
    queryFn: () => flagsApi.list(projectId).then((response) => response.data),
    enabled: !!projectId,
  })

  const totalFlags = flags?.length ?? 0

  return (
    <Layout>
      <div className="p-8 max-w-6xl mx-auto">
        <div className="mb-10">
          <h2 className="text-3xl font-bold text-white tracking-tight">Dashboard</h2>
          <p className="text-slate-400 mt-2 text-sm">
            Welcome back
            {user ? (
              <span className="text-slate-300 font-medium"> · org {user.orgId.slice(0, 8)}…</span>
            ) : (
              ''
            )}
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mb-10">
          <div className="relative group overflow-hidden rounded-2xl bg-slate-900/40 backdrop-blur-xl border border-white/[0.06] p-6 hover:border-indigo-500/20 transition-all duration-300 hover:shadow-2xl hover:shadow-indigo-500/5">
            <div className="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
            <div className="relative flex items-center justify-between mb-4">
              <span className="text-slate-400 text-sm font-medium">Total Flags</span>
              <div className="w-9 h-9 rounded-xl bg-indigo-500/10 flex items-center justify-center border border-indigo-500/20">
                <Flag size={17} className="text-indigo-400" />
              </div>
            </div>
            <p className="relative text-4xl font-bold text-white tracking-tight">{totalFlags}</p>
          </div>

          <div className="relative group overflow-hidden rounded-2xl bg-slate-900/40 backdrop-blur-xl border border-white/[0.06] p-6 hover:border-emerald-500/20 transition-all duration-300 hover:shadow-2xl hover:shadow-emerald-500/5">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
            <div className="relative flex items-center justify-between mb-4">
              <span className="text-slate-400 text-sm font-medium">Project</span>
              <div className="w-9 h-9 rounded-xl bg-emerald-500/10 flex items-center justify-center border border-emerald-500/20">
                <Layers size={17} className="text-emerald-400" />
              </div>
            </div>
            <p className="relative text-xl font-semibold text-white truncate">
              {project?.name ?? '—'}
            </p>
          </div>

          <div className="relative group overflow-hidden rounded-2xl bg-slate-900/40 backdrop-blur-xl border border-white/[0.06] p-6 hover:border-emerald-500/20 transition-all duration-300 hover:shadow-2xl hover:shadow-emerald-500/5">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
            <div className="relative flex items-center justify-between mb-4">
              <span className="text-slate-400 text-sm font-medium">Status</span>
              <div className="w-9 h-9 rounded-xl bg-emerald-500/10 flex items-center justify-center border border-emerald-500/20">
                <Activity size={17} className="text-emerald-400" />
              </div>
            </div>
            <div className="relative flex items-center gap-2.5">
              <span className="relative flex h-2.5 w-2.5">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-40" />
                <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500" />
              </span>
              <span className="text-white font-semibold">Operational</span>
            </div>
          </div>
        </div>

        <div className="rounded-2xl bg-slate-900/40 backdrop-blur-xl border border-white/[0.06] overflow-hidden">
          <div className="px-6 py-5 border-b border-white/[0.06] flex items-center justify-between">
            <h3 className="text-white font-semibold text-sm">Recent Flags</h3>
            <span className="text-xs text-slate-500 font-medium">
              {totalFlags} total
            </span>
          </div>
          {!flags || flags.length === 0 ? (
            <div className="p-12 text-center">
              <div className="w-12 h-12 rounded-2xl bg-slate-900 border border-white/[0.06] flex items-center justify-center mx-auto mb-4">
                <Flag size={20} className="text-slate-600" />
              </div>
              <p className="text-slate-400 font-medium">No flags yet</p>
              <p className="text-slate-500 text-sm mt-1">Create your first feature flag to get started</p>
            </div>
          ) : (
            <div className="divide-y divide-white/[0.04]">
              {flags.slice(0, 5).map((item) => (
                <div
                  key={item.id}
                  className="px-6 py-4 flex items-center justify-between group hover:bg-white/[0.02] transition-colors"
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <span className="text-white font-medium text-sm truncate">{item.name}</span>
                    <span className="text-[11px] text-indigo-300 bg-indigo-500/10 border border-indigo-500/20 px-2 py-0.5 rounded-md font-mono shrink-0">
                      {item.key}
                    </span>
                  </div>
                  <span className="text-xs text-slate-500 font-medium shrink-0">
                    {new Date(item.created_at).toLocaleDateString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  )
}