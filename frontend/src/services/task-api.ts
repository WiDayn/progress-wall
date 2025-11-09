import { BaseApiService, type ApiResponse } from './base-api'
import { getEndpointUrl, getEndpointMethod } from '@/config/api'
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
  newColumnId: string
  newOrder: number
}

export interface MoveTaskResponse {
  success: boolean
  message?: string
  data?: {
    taskId: string
    columnId: string
    position: number
  }
}

export class TaskApiService extends BaseApiService {
  
  async getTaskDetail(taskId: string): Promise<ApiResponse<TaskDetailResponse>> {
    const url = `${getEndpointUrl('TASK_DETAIL')}/${taskId}`
    
    return this.request<TaskDetailResponse>(url, {
      method: getEndpointMethod('TASK_DETAIL')
    })
  }

  async moveTask(taskId: string, request: MoveTaskRequest): Promise<ApiResponse<MoveTaskResponse>> {
    const url = `${getEndpointUrl('TASK_MOVE')}/${taskId}/move`
    
    return this.request<MoveTaskResponse>(url, {
      method: getEndpointMethod('TASK_MOVE'),
      body: JSON.stringify(request),
    })
  }
}

export const taskApiService = new TaskApiService()
