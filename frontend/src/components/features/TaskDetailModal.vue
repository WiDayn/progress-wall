<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm" @click.self="close">
    <div class="bg-background w-full max-w-5xl h-[85vh] rounded-xl shadow-2xl flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      
      <TaskDetailHeader 
        :taskId="task?.id"
        :title="isEditing ? editForm.title : (task?.title || '')"
        :isEditing="isEditing"
        @update:title="editForm.title = $event"
        @save="saveChanges"
        @cancel="cancelEditing"
        @start-editing="startEditing"
        @delete="handleDelete"
        @close="close"
      />

      <div class="flex flex-1 min-h-0">
        <div class="flex-1 overflow-y-auto" v-if="loading && !task">
           <div class="flex items-center justify-center h-full">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
        </div>
        
        <TaskDetailContent 
          v-else-if="task"
          :task="task"
          :isEditing="isEditing"
          :description="isEditing ? editForm.description : (task.description || '')"
          @update:priority="updatePriority"
          @update:description="editForm.description = $event"
          @start-editing="startEditing"
        />

        <TaskActivityLog 
          :activities="activities"
          :loading="loadingActivities"
        />
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useKanbanStore, type Task } from '@/stores/kanban'
import type { ActivityLog } from '@/services/task-api'
import api from '@/lib/api'

import TaskDetailHeader from './task/TaskDetailHeader.vue'
import TaskDetailContent from './task/TaskDetailContent.vue'
import TaskActivityLog from './task/TaskActivityLog.vue'

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

// Immediate update for Priority
const updatePriority = async (newPriority: number) => {
  if (!task.value) return
  
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