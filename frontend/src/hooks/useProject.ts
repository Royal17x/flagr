import { useQuery } from '@tanstack/react-query'
import { projectsApi, type Project } from '../api/projects'

export function useProject() {
  const { data: projects } = useQuery({
    queryKey: ['projects'],
    queryFn: (): Promise<Project[]> => projectsApi.list().then((r) => r.data),
  })
  return projects?.[0] ?? null
}