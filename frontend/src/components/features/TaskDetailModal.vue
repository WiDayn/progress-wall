<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" @click.self="close">
    <div class="bg-background w-full max-w-5xl h-[85vh] rounded-xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b">
        <div class="flex items-center gap-3 flex-1">
          <span class="text-sm font-mono text-muted-foreground">#{{ task?.id }}</span>
          
          <!-- 可编辑标题 -->
          <div v-if="isEditing" class="flex-1 max-w-xl">
            <input 
              v-model="editForm.title"
              class="w-full px-2 py-1 text-xl font-bold border rounded focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="任务标题"
              @keyup.enter="saveChanges"
            />
          </div>
          <h2 v-else class="text-xl font-bold text-foreground truncate max-w-md cursor-pointer hover:text-primary transition-colors" :title="task?.title" @click="startEditing">
            {{ task?.title || '加载中...' }}
          </h2>
        </div>

        <div class="flex items-center gap-2">
           <div v-if="isEditing" class="flex gap-2 mr-4">
            <button @click="saveChanges" class="px-3 py-1 text-sm bg-primary text-primary-foreground rounded hover:bg-primary/90">保存</button>
            <button @click="cancelEditing" class="px-3 py-1 text-sm bg-secondary text-secondary-foreground rounded hover:bg-secondary/80">取消</button>
          </div>
          
          <!-- 删除按钮 -->
          <button 
            @click="handleDelete"
            class="text-muted-foreground hover:text-destructive p-1 rounded-md hover:bg-destructive/10 transition-colors mr-2"
            title="删除任务"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
          </button>

          <button @click="close" class="text-muted-foreground hover:text-foreground p-1 rounded-md hover:bg-muted transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
          </button>
        </div>
      </div>

      <div class="flex flex-1 min-h-0">
        <!-- Left: Task Detail -->
        <div class="flex-1 p-8 overflow-y-auto">
          <div v-if="loading && !task" class="flex items-center justify-center h-full">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
          
          <div v-else-if="task" class="space-y-8">
            <!-- Status & Priority -->
            <div class="flex flex-wrap gap-4 items-center p-3 bg-muted/10 rounded-lg border border-border/50">
              <!-- 状态选择 (Always editable) -->
              <!-- 状态由所在列决定，此处不再显示或编辑 -->
              <!-- 
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground font-medium uppercase">状态</span>
                <select 
                  :value="task.status"
                  @change="updateStatus(($event.target as HTMLSelectElement).value)"
                  :class="['h-7 rounded-md border bg-background px-2 py-0 text-xs shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring cursor-pointer', getStatusClass(task.status)]"
                >
                  <option :value="1">待办</option>
                  <option :value="2">进行中</option>
                  <option :value="3">已完成</option>
                  <option :value="4">已取消</option>
                </select>
              </div>

              <div class="h-4 w-px bg-border/50 mx-1"></div>
              -->

              <!-- 当前列信息 -->
              <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground font-medium uppercase">当前列</span>
                <span class="px-2 py-0.5 rounded-full text-xs font-medium border bg-muted text-muted-foreground">
                  {{ task.column?.name || '加载中...' }}
                </span>
              </div>

              <div class="h-4 w-px bg-border/50 mx-1"></div>

              <!-- 优先级选择 (Always editable) -->
               <div class="flex items-center gap-2">
                <span class="text-xs text-muted-foreground font-medium uppercase">优先级</span>
                <select 
                  :value="task.priority"
                  @change="updatePriority(($event.target as HTMLSelectElement).value)"
                  :class="['h-7 rounded-md border bg-background px-2 py-0 text-xs shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring cursor-pointer', getPriorityClass(task.priority)]"
                >
                  <option :value="1">低</option>
                  <option :value="2">中</option>
                  <option :value="3">高</option>
                  <option :value="4">紧急</option>
                </select>
              </div>

              <div class="h-4 w-px bg-border/50 mx-1" v-if="task.assignee"></div>

              <!-- 执行人 (暂只读) -->
              <div class="flex items-center gap-2" v-if="task.assignee">
                <span class="text-xs text-muted-foreground font-medium uppercase">执行人</span>
                <div class="flex items-center gap-1.5">
                  <Avatar :fallback="task.assignee.username.substring(0, 2).toUpperCase()" class="h-5 w-5 text-[10px]" />
                  <span class="text-sm font-medium">{{ task.assignee.nickname || task.assignee.username }}</span>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <h3 class="text-sm font-medium text-muted-foreground uppercase tracking-wider">Description</h3>
                <button v-if="!isEditing" @click="startEditing" class="text-xs text-primary hover:underline">编辑</button>
              </div>
              
              <div v-if="isEditing">
                <textarea
                  v-model="editForm.description"
                  class="w-full min-h-[150px] rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                  placeholder="添加任务描述..."
                ></textarea>
              </div>
              <div v-else 
                class="bg-muted/30 rounded-lg p-4 text-sm leading-relaxed text-foreground/90 min-h-[100px] prose prose-sm dark:prose-invert max-w-none"
                v-html="renderMarkdown(task.description || '暂无描述')">
              </div>
            </div>

            <!-- Meta Info -->
            <div class="grid grid-cols-2 gap-6 pt-4 border-t">
              <div>
                <span class="text-xs text-muted-foreground block mb-1">创建于</span>
                <span class="text-sm">{{ formatDateTime(task.created_at) }}</span>
              </div>
              <div>
                <span class="text-xs text-muted-foreground block mb-1">最后更新</span>
                <span class="text-sm">{{ formatDateTime(task.updated_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Activity Log -->
        <div class="w-96 border-l bg-muted/5 flex flex-col min-h-0">
          <div class="px-6 py-4 border-b bg-muted/10 flex items-center justify-between">
            <h3 class="font-semibold text-foreground">Activity</h3>
            <span class="text-xs text-muted-foreground">{{ activities.length }} 条记录</span>
          </div>
          
          <div class="flex-1 overflow-y-auto p-6">
            <div v-if="loadingActivities" class="flex justify-center py-4">
              <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary/50"></div>
            </div>
            
            <div v-else class="space-y-6">
              <div v-for="log in activities" :key="log.id" class="relative pl-4 group">
                <!-- Timeline Line -->
                <div class="absolute left-0 top-2 bottom-[-24px] w-px bg-border group-last:bottom-0"></div>
                
                <div class="flex gap-3">
                  <Avatar 
                    :src="getAvatarUrl(log.avatar)"
                    :fallback="(log.nickname || log.username).substring(0, 1).toUpperCase()" 
                    class="h-8 w-8 flex-shrink-0 ring-4 ring-background relative z-10"
                  />
                  <div class="flex flex-col gap-1 min-w-0">
                    <div class="text-sm text-foreground">
                      <span class="font-semibold hover:underline cursor-pointer">{{ log.nickname || log.username }}</span>
                      <span class="text-muted-foreground mx-1" v-html="formatAction(log)"></span>
                    </div>
                    <span class="text-xs text-muted-foreground" :title="formatDateTime(log.created_at)">
                      {{ formatRelativeTime(log.created_at) }}
                    </span>
                  </div>
                </div>
              </div>

              <div v-if="activities.length === 0" class="text-center py-8 text-muted-foreground text-sm">
                暂无活动记录
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useKanbanStore, type Task } from '@/stores/kanban'
import type { ActivityLog } from '@/services/task-api'
import Avatar from '@/components/ui/Avatar.vue'
import api from '@/lib/api'
import { getAvatarUrl } from '@/lib/utils'
import { marked } from 'marked'

const props = defineProps<{
  isOpen: boolean
  taskId?: number
}>()

const emit = defineEmits(['close'])

const store = useKanbanStore()
const task = ref<Task | null>(null)
const activities = ref<ActivityLog[]>([])
const loading = ref(false)
const loadingActivities = ref(false)

// 编辑状态
const isEditing = ref(false)
const editForm = reactive({
  title: '',
  description: '',
  // Status & Priority are now handled directly
})

const close = () => {
  if (isEditing.value) {
    if (!confirm('您有未保存的更改，确定要关闭吗？')) return
    cancelEditing()
  }
  emit('close')
}

const startEditing = () => {
  if (!task.value) return
  editForm.title = task.value.title
  editForm.description = task.value.description || ''
  isEditing.value = true
}

const cancelEditing = () => {
  isEditing.value = false
}

// Immediate update for Status
const updateStatus = async (value: string | number) => {
  if (!task.value) return
  const newStatus = Number(value)
  
  // Optimistic update
  const originalStatus = task.value.status
  task.value.status = newStatus
  
  try {
    await api.put(`/tasks/${task.value.id}`, { status: newStatus })
    store.updateTask(task.value.id, { status: newStatus })
    await loadData() // Reload to get new logs
  } catch (error) {
    console.error('Failed to update status:', error)
    task.value.status = originalStatus // Rollback
    alert('更新状态失败')
  }
}

// Immediate update for Priority
const updatePriority = async (value: string | number) => {
  if (!task.value) return
  const newPriority = Number(value)
  
  // Optimistic update
  const originalPriority = task.value.priority
  task.value.priority = newPriority
  
  try {
    await api.put(`/tasks/${task.value.id}`, { priority: newPriority })
    store.updateTask(task.value.id, { priority: newPriority })
    await loadData() // Reload to get new logs
  } catch (error) {
    console.error('Failed to update priority:', error)
    task.value.priority = originalPriority // Rollback
    alert('更新优先级失败')
  }
}

const saveChanges = async () => {
  if (!task.value) return
  
  try {
    await api.put(`/tasks/${task.value.id}`, {
      title: editForm.title,
      description: editForm.description
    })
    
    store.updateTask(task.value.id, {
      title: editForm.title,
      description: editForm.description
    })

    await loadData()
    isEditing.value = false
  } catch (error) {
    console.error('Failed to update task:', error)
    alert('保存失败，请重试')
  }
}

const handleDelete = async () => {
  if (!task.value) return
  if (!confirm('确定要删除这个任务吗？此操作不可恢复。')) return

  try {
    await store.deleteTask(task.value.id)
    emit('close')
  } catch (error) {
    console.error('Failed to delete task:', error)
    alert('删除任务失败，请重试')
  }
}

const loadData = async () => {
  if (!props.taskId) return

  loading.value = true
  loadingActivities.value = true
  
  try {
    const [taskRes, activityRes] = await Promise.all([
      store.fetchTaskDetail(props.taskId.toString()),
      store.fetchTaskActivities(props.taskId)
    ])

    if (taskRes.success) {
      task.value = taskRes.data
    }
    activities.value = activityRes.data
  } catch (error) {
    console.error('Failed to load task data:', error)
  } finally {
    loading.value = false
    loadingActivities.value = false
  }
}

watch(() => props.isOpen, (newVal) => {
  if (newVal && props.taskId) {
    loadData()
    isEditing.value = false
  } else {
    task.value = null
    activities.value = []
    isEditing.value = false
  }
})

// Helpers
const renderMarkdown = (text: string) => {
  return marked.parse(text || '')
}

const formatAction = (log: ActivityLog) => {
  let text = log.description
  text = text.replace(/moved this task/g, 'moved this task')
  text = text.replace(/from "(.+?)"/g, 'from <b>$1</b>')
  text = text.replace(/to "(.+?)"/g, 'to <b>$1</b>')
  text = text.replace(/created this task/g, 'created this task')
  text = text.replace(/updated task/g, 'updated task')
  
  return text
}

const formatDateTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatRelativeTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  
  if (diffInSeconds < 60) return '刚刚'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} 分钟前`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} 小时前`
  if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} 天前`
  
  return date.toLocaleDateString()
}

const getStatusClass = (status: number) => {
  const map: Record<number, string> = {
    1: 'bg-blue-50 text-blue-700 border-blue-200', // Todo
    2: 'bg-yellow-50 text-yellow-700 border-yellow-200', // In Progress
    3: 'bg-green-50 text-green-700 border-green-200', // Done
    4: 'bg-gray-50 text-gray-700 border-gray-200' // Cancelled
  }
  return map[status] || 'bg-gray-50 text-gray-700'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    1: '待办',
    2: '进行中',
    3: '已完成',
    4: '已取消'
  }
  return map[status] || '未知'
}

const getPriorityClass = (priority: number) => {
  const map: Record<number, string> = {
    1: 'bg-green-50 text-green-700 border-green-200', // Low
    2: 'bg-blue-50 text-blue-700 border-blue-200', // Medium
    3: 'bg-orange-50 text-orange-700 border-orange-200', // High
    4: 'bg-red-50 text-red-700 border-red-200' // Urgent
  }
  return map[priority] || 'bg-gray-50 text-gray-700'
}

const getPriorityText = (priority: number) => {
  const map: Record<number, string> = {
    1: '低',
    2: '中',
    3: '高',
    4: '紧急'
  }
  return map[priority] || '无'
}
</script>

<style scoped>
/* 自定义滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  @apply bg-muted-foreground/20 rounded-full;
}
::-webkit-scrollbar-thumb:hover {
  @apply bg-muted-foreground/40;
}
</style>
