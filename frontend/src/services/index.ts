// API 服务统一导出
export { BaseApiService, type ApiResponse } from './base-api'
export { 
  BoardApiService, 
  boardApiService, 
  type KanbanListResponse,
  type BoardDetail,
  type Column,
  type Task as BoardTask,
  type User
} from './board-api'
export { 
  TaskApiService, 
  taskApiService, 
  type TaskDetailRequest,
  type TaskDetailResponse,
  type MoveTaskRequest,
  type MoveTaskResponse 
} from './task-api'
export * from './team-api'
