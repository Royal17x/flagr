import client from './client'

export interface Project {
  id: string
  organization_id: string
  name: string
  description: string
}

export const projectsApi = {
  list: () => client.get<Project[]>('/projects'),
}