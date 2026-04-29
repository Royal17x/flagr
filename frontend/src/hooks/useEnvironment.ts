import client from '../api/client'

export interface Environment {
  id: string
  project_id: string
  name: string
  slug: string
}

export const environmentsApi = {
  list: (projectId: string) =>
    client.get<Environment[]>('/environments', { params: { project_id: projectId } }),
}