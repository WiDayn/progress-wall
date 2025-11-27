import api from '@/lib/api'

export interface Team {
  id: number
  name: string
  description: string
  creator_id: number
  created_at: string
  updated_at: string
}

export interface TeamMember {
  id: number
  team_id: number
  user_id: number
  role: number // 1: Member, 2: Admin
  joined_at: string
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

export interface CreateTeamDto {
  name: string
  description?: string
}

export const teamApi = {
  // Get all teams for the current user
  getMyTeams: () => {
    return api.get<{ teams: Team[] }>('/teams')
  },

  // Create a new team
  createTeam: (data: CreateTeamDto) => {
    return api.post<Team>('/teams', data)
  },

  // Get team details
  getTeam: (teamId: number) => {
    return api.get<Team>(`/teams/${teamId}`)
  },

  // Get team members
  getMembers: (teamId: number) => {
    return api.get<{ members: TeamMember[] }>(`/teams/${teamId}/members`)
  },

  // Add member to team
  addMember: (teamId: number, userId: number, role: number = 1) => {
    return api.post(`/teams/${teamId}/members`, { user_id: userId, role })
  }
}

