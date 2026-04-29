import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Key, Trash2, Copy, Check } from 'lucide-react'
import Layout from '../components/Layout'
import { useProject } from '../hooks/useProject'
import client from '../api/client'
import { environmentsApi } from '../hooks/useEnvironment'

interface SDKKey {
  id: string
  name: string
  created_at: string
  expires_at: string | null
}

export default function SDKKeys() {
  const queryClient = useQueryClient()
  const project = useProject()
  const projectId = project?.id ?? ''

  const [newKeyName, setNewKeyName] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [createdKey, setCreatedKey] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const { data: keys, isLoading } = useQuery({
    queryKey: ['sdk-keys', projectId],
    queryFn: () => client.get<SDKKey[]>(`/sdk-keys?project_id=${projectId}`).then((r) => r.data),
    enabled: !!projectId,
  })

  const { data: environments } = useQuery({
    queryKey: ['environments', projectId],
    queryFn: () => environmentsApi.list(projectId).then((r) => r.data),
    enabled: !!projectId,
  })

  const defaultEnv = environments?.[0]



  const createMutation = useMutation({
    mutationFn: (name: string) =>
    client.post('/sdk-keys', {
        project_id: projectId,
        environment_id: defaultEnv?.id ?? '',
        name,
    }).then((r) => r.data),
    onSuccess: (data) => {
      setCreatedKey(data.key)
      setShowCreate(false)
      setNewKeyName('')
      queryClient.invalidateQueries({ queryKey: ['sdk-keys'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => client.delete(`/sdk-keys/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sdk-keys'] }),
  })

  const handleCopy = () => {
    if (createdKey) {
      navigator.clipboard.writeText(createdKey)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <Layout>
      <div className="p-8 max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-white tracking-tight">SDK Keys</h2>
            <p className="text-slate-400 mt-2 text-sm">
              API keys for your SDK integrations
            </p>
          </div>
          <button
            onClick={() => setShowCreate(true)}
            className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 text-white pl-4 pr-5 py-2.5 rounded-xl font-medium text-sm shadow-lg shadow-indigo-500/20 hover:-translate-y-0.5 transition-all"
          >
            <Plus size={17} />
            New Key
          </button>
        </div>

        {createdKey && (
          <div className="mb-6 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 p-6">
            <p className="text-emerald-400 font-semibold mb-1">Key created — save it now!</p>
            <p className="text-slate-400 text-sm mb-3">This key will only be shown once.</p>
            <div className="flex items-center gap-3">
              <code className="flex-1 bg-slate-950/50 border border-slate-800 text-emerald-300 px-4 py-2.5 rounded-xl font-mono text-sm break-all">
                {createdKey}
              </code>
              <button onClick={handleCopy}
                className="flex items-center gap-2 bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2.5 rounded-xl text-sm font-medium shrink-0">
                {copied ? <Check size={16} /> : <Copy size={16} />}
                {copied ? 'Copied!' : 'Copy'}
              </button>
            </div>
          </div>
        )}

        {showCreate && (
          <div className="mb-6 rounded-2xl bg-slate-900/60 border border-white/[0.08] p-6">
            <h3 className="text-white font-semibold mb-4">Create SDK Key</h3>
            <div className="flex gap-3">
              <input
                placeholder="Key name (e.g. Production SDK)"
                value={newKeyName}
                onChange={(e) => setNewKeyName(e.target.value)}
                className="flex-1 bg-slate-950/50 border border-slate-800 text-white px-4 py-2.5 rounded-xl focus:outline-none focus:border-indigo-500/50 placeholder:text-slate-600"
              />
              <button
                onClick={() => createMutation.mutate(newKeyName)}
                disabled={!newKeyName || createMutation.isPending}
                className="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-40 text-white px-5 py-2.5 rounded-xl text-sm font-medium"
              >
                Create
              </button>
              <button onClick={() => setShowCreate(false)}
                className="bg-white/[0.04] text-slate-300 px-4 py-2.5 rounded-xl text-sm border border-white/[0.06]">
                Cancel
              </button>
            </div>
          </div>
        )}

        {isLoading ? (
          <div className="space-y-3">
            {[1, 2].map(i => (
              <div key={i} className="bg-slate-900/40 border border-white/[0.05] rounded-2xl p-5 animate-pulse h-16" />
            ))}
          </div>
        ) : !keys || keys.length === 0 ? (
          <div className="text-center py-16 border border-white/[0.06] rounded-2xl bg-slate-900/20">
            <Key size={32} className="mx-auto mb-4 text-slate-700" />
            <p className="text-slate-300 font-medium">No SDK keys yet</p>
          </div>
        ) : (
          <div className="space-y-3">
            {keys.map((key) => (
              <div key={key.id}
                className="group bg-slate-900/40 border border-white/[0.06] rounded-2xl p-5 flex items-center justify-between hover:border-slate-700 transition-all">
                <div className="flex items-center gap-4">
                  <div className="w-9 h-9 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center">
                    <Key size={16} className="text-indigo-400" />
                  </div>
                  <div>
                    <p className="text-white font-medium text-sm">{key.name}</p>
                    <p className="text-slate-500 text-xs mt-0.5">
                      Created {new Date(key.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => deleteMutation.mutate(key.id)}
                  className="opacity-0 group-hover:opacity-100 text-slate-500 hover:text-red-400 p-2 rounded-xl hover:bg-red-500/10 transition-all"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </Layout>
  )
}