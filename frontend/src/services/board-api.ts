import api from '@/lib/api'
import type { Board } from '@/stores/board'

// API 响应格式
export interface ApiResponse<T> {
  msg?: string
  data?: T
}

// 看板列表响应接口
export interface KanbanListResponse extends Array<Board> {}

// 用户信息接口
export interface User {
  id: number
  username: string
  email: string
}

// 任务接口（后端返回格式）
export interface Task {
  id: number
  title: string
  description: string
  priority: number
  status: number
  position: number
  column_id: number
  creator_id: number
  assignee_id: number
  project_id: number
  created_at: string
  assignee?: User
  creator?: User
}

// 列接口（后端返回格式）
export interface Column {
  id: number
  name: string
  description: string
  color: string
  position: number
  board_id: number
  status: number
  created_at: string
  tasks: Task[]
}

// 看板详情接口（后端返回格式 - 包含嵌套的列和任务）
export interface BoardDetail {
  id: number
  name: string
  description: string
  color: string
  status: number
  project_id: number
  owner_id: number
  position: number
  created_at: string
  updated_at: string
  owner: User
  columns: Column[]
}

// 看板 API 服务
export class BoardApiService {
  // 获取看板列表
  async getKanbanList(): Promise<ApiResponse<KanbanListResponse>> {
    try {
      const response = await api.get('/boards')
      return response.data
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '获取看板列表失败',
        data: []
      }
    }
  }

  /**
   * 获取单个看板详情（包含嵌套的列和任务）
   * 需要认证: Bearer Token (JWT)
   * 
   * @param boardId - 看板ID
   * @returns 完整的看板数据，包括 columns 和 tasks
   */
  async getKanbanDetail(boardId: string | number): Promise<BoardDetail | null> {
    try {
      const response = await api.get(`/boards/${boardId}`)
      return response.data
    } catch (error: any) {
      console.error('获取看板详情失败:', error.response?.data || error.message)
      
      // 处理不同的错误状态
      if (error.response?.status === 400) {
        throw new Error('无效的看板ID')
      } else if (error.response?.status === 404) {
        throw new Error('看板不存在')
      } else {
        throw new Error(error.response?.data?.msg || error.message || '获取看板详情失败')
      }
    }
  }

  // 创建看板
  async createKanban(board: Omit<Board, 'id' | 'createdAt' | 'updatedAt'>, projectId: string): Promise<ApiResponse<Board>> {
    try {
      const response = await api.post(`/projects/${projectId}/boards`, board)
      return response.data
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '创建看板失败',
        data: undefined
      }
    }
  }

  // 更新看板
  async updateKanban(boardId: string, updates: Partial<Board>): Promise<ApiResponse<Board>> {
    try {
      const response = await api.put(`/boards/${boardId}`, updates)
      return response.data
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '更新看板失败',
        data: undefined
      }
    }
  }

  // 删除看板
  async deleteKanban(boardId: string): Promise<ApiResponse<void>> {
    try {
      const response = await api.delete(`/boards/${boardId}`)
      return response.data
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '删除看板失败',
        data: undefined
      }
    }
  }
}

// 导出单例实例
export const boardApiService = new BoardApiService()
