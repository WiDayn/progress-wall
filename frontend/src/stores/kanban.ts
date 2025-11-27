import { defineStore } from 'pinia'
import { ref } from 'vue'
import { taskApiService, boardApiService, columnApiService } from '@/services'
import type { Column as ApiColumn, Task as ApiTask } from '@/services/board-api'
import type { CreateColumnRequest } from '@/services/column-api'
import type { MoveTaskRequest, MoveTaskResponse, CreateTaskRequest } from '@/services/task-api'

// 前端使用的 Task 类型 (基于 API 类型扩展或适配)
export interface Task extends ApiTask {
  // 可以添加前端特有的字段
}

// 前端使用的 Column 类型
export interface Column extends ApiColumn {
  tasks: Task[]
}

// 任务详情响应
export interface TaskDetailResponse {
  success: boolean
  data: Task
  message?: string
}

export const useKanbanStore = defineStore('kanban', () => {
  const columns = ref<Column[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const currentBoardId = ref<string | number | null>(null)
  const currentProjectId = ref<number | null>(null)

  // 获取看板详情（包含列和任务）
  const fetchBoardDetail = async (boardId: string | number) => {
    isLoading.value = true
    error.value = null
    currentBoardId.value = boardId
    
    try {
      const boardDetail = await boardApiService.getKanbanDetail(boardId)
      if (boardDetail) {
        currentProjectId.value = boardDetail.project_id // 保存项目ID，创建任务时需要
        // 确保 tasks 数组存在
        columns.value = (boardDetail.columns || []).map(col => ({
          ...col,
          tasks: col.tasks || []
        }))
      } else {
        columns.value = []
        currentProjectId.value = null
      }
    } catch (err: any) {
      console.error('Fetch board detail failed:', err)
      error.value = err.message || '获取看板详情失败'
      columns.value = []
      currentProjectId.value = null
    } finally {
      isLoading.value = false
    }
  }

  // 创建新列
  const createColumn = async (data: CreateColumnRequest) => {
    if (!currentBoardId.value) return

    try {
      const response = await columnApiService.createColumn(Number(currentBoardId.value), data)
      if (response) {
        // 重新获取看板数据以保持同步（或者手动添加到 columns）
        // 简单起见，如果返回了新列数据，手动添加
        const newColumn: Column = {
          ...response as unknown as ApiColumn, // 类型转换
          tasks: []
        }
        columns.value.push(newColumn)
      }
    } catch (err: any) {
      console.error('Create column failed:', err)
      throw err
    }
  }

  const addTask = (task: Task) => {
    const column = columns.value.find(col => col.id === task.column_id)
    if (column) {
      column.tasks.push(task)
    }
  }

  // 创建任务
  const createTask = async (columnId: number, data: { title: string; description: string; priority: number }) => {
    if (!currentProjectId.value) {
      throw new Error('无法创建任务：缺少项目ID')
    }

    const request: CreateTaskRequest = {
      title: data.title,
      description: data.description,
      priority: data.priority,
      column_id: columnId,
      project_id: currentProjectId.value,
      status: 1, // 默认待办
      position: 0 // 默认位置
    }

    try {
      const response = await taskApiService.createTask(columnId.toString(), request)
      if (response.data) {
        // 添加到本地状态
        const newTask = response.data as unknown as Task
        addTask(newTask)
      }
    } catch (err: any) {
      console.error('Create task failed:', err)
      throw err
    }
  }

  const updateTask = (taskId: number, updates: Partial<Task>) => {
    for (const column of columns.value) {
      const task = column.tasks.find(t => t.id === taskId)
      if (task) {
        Object.assign(task, updates, { updated_at: new Date().toISOString() })
        break
      }
    }
  }

  // 拖拽移动任务（乐观更新）
  const moveTaskWithDrag = async (taskId: number, newColumnId: number, newOrder: number) => {
    // 保存原始状态用于回滚
    const originalColumns = JSON.parse(JSON.stringify(columns.value))
    
    try {
      // 1. 乐观更新：立即更新本地状态
      let taskToMove: Task | null = null
      let sourceColumnIndex = -1
      let sourceTaskIndex = -1

      // 找到要移动的任务
      for (let i = 0; i < columns.value.length; i++) {
        const taskIndex = columns.value[i].tasks.findIndex(t => t.id === taskId)
        if (taskIndex !== -1) {
          taskToMove = columns.value[i].tasks[taskIndex]
          sourceColumnIndex = i
          sourceTaskIndex = taskIndex
          break
        }
      }

      if (!taskToMove) {
        throw new Error('任务不存在')
      }

      // 从源列移除任务
      columns.value[sourceColumnIndex].tasks.splice(sourceTaskIndex, 1)

      // 找到目标列
      const targetColumn = columns.value.find(col => col.id === newColumnId)
      if (!targetColumn) {
        throw new Error('目标列不存在')
      }

      // 更新任务列ID
      taskToMove.column_id = newColumnId
      // taskToMove.updated_at = new Date().toISOString()

      // 插入到目标列的指定位置
      if (newOrder >= targetColumn.tasks.length) {
        targetColumn.tasks.push(taskToMove)
      } else {
        targetColumn.tasks.splice(newOrder, 0, taskToMove)
      }

      // 2. 调用API
      // 注意：这里不使用 await，让 API 在后台异步执行，从而避免 UI 阻塞
      moveTaskAPI(taskId.toString(), { newColumnId: newColumnId, newOrder })
        .catch(error => {
          // 只有在 API 失败时才回滚
          columns.value = originalColumns
          console.error('Move task failed:', error)
          // 可以添加一个全局提示，告诉用户同步失败
        })

    } catch (error) {
      // 如果是同步逻辑出错（比如找不到列），立即回滚
      columns.value = originalColumns
      console.error('Move task local update failed:', error)
      throw error
    }
  }

  // 调用真实API移动任务
  const moveTaskAPI = async (taskId: string, request: MoveTaskRequest): Promise<MoveTaskResponse> => {
    return await taskApiService.moveTask(taskId, request)
  }

  const deleteTask = (taskId: number) => {
    for (const column of columns.value) {
      const taskIndex = column.tasks.findIndex(t => t.id === taskId)
      if (taskIndex !== -1) {
        column.tasks.splice(taskIndex, 1)
        break
      }
    }
  }

  // 调用真实API获取任务详情
  const fetchTaskDetail = async (taskId: string): Promise<TaskDetailResponse> => {
    const response = await taskApiService.getTaskDetail(taskId)
    
    if (response.data) {
      return {
        success: true,
        data: response.data as unknown as Task
      }
    }
    
    return {
      success: false,
      data: {} as Task,
      message: response.msg || '获取任务详情失败'
    }
  }

  return {
    columns,
    isLoading,
    error,
    currentBoardId,
    fetchBoardDetail,
    createColumn,
    createTask,
    addTask,
    updateTask,
    moveTaskWithDrag,
    deleteTask,
    fetchTaskDetail
  }
})
