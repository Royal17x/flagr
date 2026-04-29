import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Flag as FlagIcon, Trash2, Sparkles } from 'lucide-react'
import Layout from '../components/Layout'
import { flagsApi, type Flag } from '../api/flags'
import { environmentsApi } from '../hooks/useEnvironment'
import { useProject } from '../hooks/useProject'

export default function Flags() {
  const queryClient = useQueryClient()
  const project = useProject()
  const projectId = project?.id ?? ''

  const [showCreate, setShowCreate] = useState(false)
  const [newFlag, setNewFlag] = useState({ key: '', name: '', description: '' })

  const { data: flags, isLoading } = useQuery({
    queryKey: ['flags', projectId],
    queryFn: () => flagsApi.list(projectId).then((r) => r.data),
    enabled: !!projectId,
  })

  const { data: environments } = useQuery({
    queryKey: ['environments', projectId],
    queryFn: () => environmentsApi.list(projectId).then((r) => r.data),
    enabled: !!projectId,
  })

  const defaultEnv = environments?.[0]

  const createMutation = useMutation({
    mutationFn: flagsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['flags'] })
      setShowCreate(false)
      setNewFlag({ key: '', name: '', description: '' })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: flagsApi.remove,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['flags'] }),
  })

  const toggleMutation = useMutation({
    mutationFn: ({ flagId, enabled }: { flagId: string; enabled: boolean }) =>
      flagsApi.toggle(flagId, defaultEnv?.id ?? '', enabled),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['flags'] }),
  })

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate({ ...newFlag, project_id: projectId, type: 'boolean' })
  }

  return (
    <Layout>
      <div className="p-8 max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-white tracking-tight">Feature Flags</h2>
            <p className="text-slate-400 mt-2 text-sm">
              {project ? (
                <span className="inline-flex items-center gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-indigo-500" />
                  {project.name}
                  {defaultEnv && (
                    <span className="text-slate-500">· {defaultEnv.name}</span>
                  )}
                </span>
              ) : 'Loading...'}
            </p>
          </div>
          <button
            onClick={() => setShowCreate(true)}
            disabled={!projectId}
            className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-40 text-white pl-4 pr-5 py-2.5 rounded-xl transition-all font-medium text-sm shadow-lg shadow-indigo-500/20 hover:-translate-y-0.5 active:translate-y-0"
          >
            <Plus size={17} strokeWidth={2.5} />
            New Flag
          </button>
        </div>

        {showCreate && (
          <div className="mb-8 rounded-2xl bg-slate-900/60 backdrop-blur-xl border border-white/[0.08] p-6">
            <div className="flex items-center gap-2 mb-5">
              <div className="w-7 h-7 rounded-lg bg-indigo-500/10 flex items-center justify-center border border-indigo-500/20">
                <Sparkles size={14} className="text-indigo-400" />
              </div>
              <h3 className="text-white font-semibold text-sm">Create Flag</h3>
            </div>
            <form onSubmit={handleCreate} className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-1.5">
                <label className="block text-xs font-medium text-slate-400 uppercase tracking-wider">Key</label>
                <input placeholder="checkout-v2" value={newFlag.key}
                  onChange={(e) => setNewFlag({ ...newFlag, key: e.target.value })}
                  className="w-full bg-slate-950/50 border border-slate-800 text-white px-3.5 py-2.5 rounded-xl focus:outline-none focus:border-indigo-500/50 font-mono text-sm placeholder:text-slate-600 transition-all"
                  required />
              </div>
              <div className="space-y-1.5">
                <label className="block text-xs font-medium text-slate-400 uppercase tracking-wider">Name</label>
                <input placeholder="Checkout V2" value={newFlag.name}
                  onChange={(e) => setNewFlag({ ...newFlag, name: e.target.value })}
                  className="w-full bg-slate-950/50 border border-slate-800 text-white px-3.5 py-2.5 rounded-xl focus:outline-none focus:border-indigo-500/50 placeholder:text-slate-600 transition-all"
                  required />
              </div>
              <div className="space-y-1.5">
                <label className="block text-xs font-medium text-slate-400 uppercase tracking-wider">Description</label>
                <input placeholder="Optional" value={newFlag.description}
                  onChange={(e) => setNewFlag({ ...newFlag, description: e.target.value })}
                  className="w-full bg-slate-950/50 border border-slate-800 text-white px-3.5 py-2.5 rounded-xl focus:outline-none focus:border-indigo-500/50 placeholder:text-slate-600 transition-all" />
              </div>
              <div className="md:col-span-3 flex gap-3 pt-1">
                <button type="submit" disabled={createMutation.isPending}
                  className="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-40 text-white px-5 py-2.5 rounded-xl text-sm font-medium shadow-lg shadow-indigo-500/20">
                  {createMutation.isPending ? 'Creating...' : 'Create Flag'}
                </button>
                <button type="button" onClick={() => setShowCreate(false)}
                  className="bg-white/[0.04] hover:bg-white/[0.08] text-slate-300 px-5 py-2.5 rounded-xl text-sm font-medium border border-white/[0.06]">
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        {isLoading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-slate-900/40 border border-white/[0.05] rounded-2xl p-6 animate-pulse h-20" />
            ))}
          </div>
        ) : !flags || flags.length === 0 ? (
          <div className="text-center py-20 border border-white/[0.06] rounded-2xl bg-slate-900/20">
            <div className="w-16 h-16 rounded-2xl bg-slate-900 border border-white/[0.06] flex items-center justify-center mx-auto mb-5">
              <FlagIcon size={28} className="text-slate-700" />
            </div>
            <p className="text-lg text-slate-300 font-medium">No flags yet</p>
            <p className="text-sm text-slate-500 mt-1.5">Create your first feature flag to get started</p>
          </div>
        ) : (
          <div className="space-y-3">
            {flags.map((flag: Flag) => (
              <FlagRow
                key={flag.id}
                flag={flag}
                envId={defaultEnv?.id ?? ''}
                onToggle={(enabled) => toggleMutation.mutate({ flagId: flag.id, enabled })}
                onDelete={() => deleteMutation.mutate(flag.id)}
              />
            ))}
          </div>
        )}
      </div>
    </Layout>
  )
}

interface FlagRowProps {
  flag: Flag
  envId: string
  onToggle: (enabled: boolean) => void
  onDelete: () => void
}

function FlagRow({ flag, envId, onToggle, onDelete }: FlagRowProps) {
  const { data: flagEnv } = useQuery({
    queryKey: ['flag-env', flag.id, envId],
    queryFn: () => flagsApi.getEnvironment(flag.id, envId).then((r) => r.data),
    enabled: !!envId,
  })

  const [optimisticEnabled, setOptimisticEnabled] = useState<boolean | null>(null)
  const isEnabled = optimisticEnabled ?? flagEnv?.enabled ?? false

  const handleToggle = () => {
    const next = !isEnabled
    setOptimisticEnabled(next)
    onToggle(next)
  }

  return (
    <div className="group bg-slate-900/40 backdrop-blur-sm border border-white/6 rounded-2xl p-5 flex items-center justify-between hover:border-indigo-500/20 hover:bg-slate-900/60 transition-all duration-300">
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-3 flex-wrap">
          <span className="text-white font-semibold text-sm">{flag.name}</span>
          <span className="text-[11px] text-indigo-300 bg-indigo-500/10 border border-indigo-500/20 px-2 py-0.5 rounded-md font-mono">
            {flag.key}
          </span>
          <span className="text-[11px] text-slate-400 bg-white/4 border border-white/6 px-2 py-0.5 rounded-md">
            {flag.type}
          </span>
        </div>
        {flag.description && (
          <p className="text-slate-500 text-sm mt-1.5">{flag.description}</p>
        )}
      </div>
      <div className="flex items-center gap-3 ml-4 shrink-0">
        {envId && (
          <button
            onClick={handleToggle}
            className={`relative w-11 h-6 rounded-full transition-all duration-300 focus:outline-none ${
              isEnabled ? 'bg-indigo-500 shadow-lg shadow-indigo-500/30' : 'bg-slate-700'
            }`}
          >
            <span className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-sm transition-transform duration-300 ${
              isEnabled ? 'translate-x-5' : 'translate-x-0'
            }`} />
          </button>
        )}
        <button onClick={onDelete}
          className="opacity-0 group-hover:opacity-100 text-slate-500 hover:text-red-400 transition-all p-2 rounded-xl hover:bg-red-500/10">
          <Trash2 size={16} />
        </button>
      </div>
    </div>
  )
}