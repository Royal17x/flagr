import client from './client'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  org_name: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

export const authApi = {
  login: (data: LoginRequest) =>
    client.post<TokenPair>('/auth/login', data),

  register: (data: RegisterRequest) =>
    client.post<TokenPair>('/auth/register', data),

  logout: (refresh_token: string) =>
    client.post('/auth/logout', { refresh_token }),
}