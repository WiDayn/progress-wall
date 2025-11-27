import api from '@/lib/api'
import type { ApiResponse } from './board-api'

export interface Column {
  id: number
  name: string
  description: string
  color: string
  position: number
  board_id: number
  status: number
  created_at: string
  tasks?: any[]
}

export interface CreateColumnRequest {
  name: string
  description?: string
  color?: string
  position?: number
}

export class ColumnApiService {
  // 获取看板的所有列
  async getColumns(boardId: number): Promise<ApiResponse<Column[]>> {
    try {
      const response = await api.get(`/boards/${boardId}/columns`)
      return {
        data: response.data.data // 后端返回格式通常是 { data: [...] }
      }
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '获取列列表失败',
        data: []
      }
    }
  }

  // 创建列
  async createColumn(boardId: number, data: CreateColumnRequest): Promise<ApiResponse<Column>> {
    try {
      const response = await api.post(`/boards/${boardId}/columns`, data)
      return response.data
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '创建列失败',
        data: undefined
      }
    }
  }
}

export const columnApiService = new ColumnApiService()

