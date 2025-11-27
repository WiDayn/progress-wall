import api from '@/lib/api'
import type { Task } from '@/stores/kanban'

export interface TaskDetailRequest {
  taskId: string
}

export interface TaskDetailResponse {
  success: boolean
  data: Task
  message?: string
}

export interface MoveTaskRequest {
  newColumnId: number
  newOrder: number
}

export interface MoveTaskResponse {
  message: string
}

export interface CreateTaskRequest {
  title: string
  description: string
  priority: number
  column_id: number
  project_id: number
  status: number
  position: number
}

// API 响应格式
export interface ApiResponse<T> {
  msg: string
  data: T
}

export class TaskApiService {
  
  /**
   * 获取任务详情
   * 需要认证：Bearer Token (JWT)
   */
  async getTaskDetail(taskId: string): Promise<ApiResponse<TaskDetailResponse>> {
    try {
      const response = await api.get(`/tasks/${taskId}`)
      return {
        msg: 'success',
        data: response.data // 直接返回 task 对象
      }
    } catch (error: any) {
      return {
        msg: error.response?.data?.msg || error.message || '获取任务详情失败',
        data: {
          success: false,
          data: {} as Task,
          message: error.response?.data?.msg || error.message
        } as unknown as TaskDetailResponse
      }
    }
  }

  /**
   * 创建任务
   */
  async createTask(columnId: string, request: CreateTaskRequest): Promise<ApiResponse<Task>> {
    try {
      const response = await api.post(`/columns/${columnId}/tasks`, request)
      return {
        msg: 'success',
        data: response.data
      }
    } catch (error: any) {
      console.error('创建任务失败:', error)
      throw new Error(error.response?.data?.error || '创建任务失败')
    }
  }

  /**
   * 移动任务（拖拽排序）
   * 需要认证：Bearer Token (JWT)
   * 
   * 支持：
   * - 跨列移动：将任务从一个列移动到另一个列
   * - 同列内移动：在同一列内调整任务顺序
   * 
   * @param taskId - 任务ID
   * @param request - 包含 newColumnId 和 newOrder
   * @returns 移动成功的消息
   */
  async moveTask(taskId: string | number, request: MoveTaskRequest): Promise<MoveTaskResponse> {
    try {
      const response = await api.patch(`/tasks/${taskId}/move`, request)
      return response.data
    } catch (error: any) {
      console.error('移动任务失败:', error.response?.data || error.message)
      
      // 处理不同的错误状态
      if (error.response?.status === 400) {
        throw new Error('请求参数错误（缺少 newColumnId 或 newOrder）')
      } else if (error.response?.status === 404) {
        throw new Error('任务不存在')
      } else {
        throw new Error(error.response?.data?.message || error.message || '移动任务失败')
      }
    }
  }
}

export const taskApiService = new TaskApiService()
