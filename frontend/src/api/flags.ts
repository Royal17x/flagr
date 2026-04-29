import client from './client'

export interface Flag {
  id: string
  project_id: string
  key: string
  name: string
  description: string
  type: string
  created_at: string
}

export interface CreateFlagRequest {
  project_id: string
  key: string
  name: string
  description: string
  type: string
}

export interface UpdateFlagRequest {
  key: string
  name: string
  description: string
  type: string
}

export const flagsApi = {
  list: (projectId: string) =>
    client.get<Flag[]>('/flags', { params: { project_id: projectId } }),

  create: (data: CreateFlagRequest) =>
    client.post<Flag>('/flags', data),

  update: (id: string, data: UpdateFlagRequest) =>
    client.put<Flag>(`/flags/${id}`, data),

  remove: (id: string) =>
    client.delete(`/flags/${id}`),
}