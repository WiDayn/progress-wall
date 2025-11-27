<template>
  <div class="min-h-screen bg-background flex flex-col h-screen">
    <div class="flex-none container mx-auto px-4 py-4">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-foreground">{{ boardName }}</h1>
          <p v-if="error" class="text-destructive text-sm mt-1">{{ error }}</p>
        </div>
        <div class="flex gap-2">
          <Button @click="refreshBoard" variant="outline" :disabled="isLoading">
            {{ isLoading ? '加载中...' : '刷新' }}
          </Button>
          <Button @click="goBack" variant="secondary">
            返回
          </Button>
        </div>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-x-auto px-4 pb-4 container mx-auto mb-4">
      <div class="flex h-full gap-6">
        <!-- 列渲染 -->
        <KanbanColumn
          v-for="column in kanbanStore.columns"
          :key="column.id"
          :column="column"
          class="w-80 flex-shrink-0 flex flex-col h-full max-h-full"
          @select-task="selectTask"
          @delete-task="deleteTask"
          @add-task="openCreateTaskDialog"
        />
        
        <!-- 添加新列按钮 -->
        <div class="w-80 flex-shrink-0 pt-2 h-full">
          <button 
            @click="openCreateDialog"
            class="w-full h-12 border-2 border-dashed border-border rounded-lg flex items-center justify-center text-muted-foreground hover:border-primary hover:text-primary transition-colors bg-muted/30 hover:bg-muted/50"
          >
            + 添加新列
          </button>
        </div>
      </div>
    </div>

    <!-- 创建列对话框 -->
    <CreateColumnDialog
      :is-open="isCreateDialogOpen"
      :loading="isCreating"
      @submit="handleCreateColumn"
      @cancel="closeCreateDialog"
    />

    <!-- 创建任务对话框 -->
    <CreateTaskDialog
      :is-open="isCreateTaskDialogOpen"
      :column-id="currentColumnId"
      :loading="isCreatingTask"
      @submit="handleCreateTask"
      @cancel="closeCreateTaskDialog"
    />

    <!-- 任务详情模态框 -->
    <TaskDetailModal
      :is-open="isTaskDetailModalOpen"
      :task-id="selectedTaskId"
      @close="closeTaskDetailModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useKanbanStore } from '@/stores/kanban'
import Button from '@/components/ui/Button.vue'
import KanbanColumn from '@/components/features/KanbanColumn.vue'
import CreateColumnDialog from '@/components/features/CreateColumnDialog.vue'
import CreateTaskDialog from '@/components/features/CreateTaskDialog.vue'
import TaskDetailModal from '@/components/features/TaskDetailModal.vue'

const route = useRoute()
const router = useRouter()
const kanbanStore = useKanbanStore()

const isLoading = computed(() => kanbanStore.isLoading)
const error = computed(() => kanbanStore.error)
const boardName = ref('项目看板')

// 创建列相关状态
const isCreateDialogOpen = ref(false)
const isCreating = ref(false)

// 创建任务相关状态
const isCreateTaskDialogOpen = ref(false)
const isCreatingTask = ref(false)
const currentColumnId = ref<number | undefined>(undefined)

// 任务详情模态框状态
const isTaskDetailModalOpen = ref(false)
const selectedTaskId = ref<number | undefined>(undefined)

const loadData = async () => {
  const boardId = route.params.boardId
  if (boardId) {
    await kanbanStore.fetchBoardDetail(boardId as string)
  }
}

onMounted(() => {
  loadData()
})

const refreshBoard = () => {
  loadData()
}

const goBack = () => {
  router.back()
}

const selectTask = (task: any) => {
  selectedTaskId.value = task.id
  isTaskDetailModalOpen.value = true
}

const closeTaskDetailModal = () => {
  isTaskDetailModalOpen.value = false
  selectedTaskId.value = undefined
}

const deleteTask = (taskId: number) => {
  if (confirm('确定要删除这个任务吗？')) {
    kanbanStore.deleteTask(taskId)
  }
}

// 列操作
const openCreateDialog = () => {
  isCreateDialogOpen.value = true
}

const closeCreateDialog = () => {
  isCreateDialogOpen.value = false
}

const handleCreateColumn = async (data: { name: string; description: string; color: string }) => {
  isCreating.value = true
  try {
    await kanbanStore.createColumn({
      name: data.name,
      description: data.description,
      color: data.color,
      position: kanbanStore.columns.length * 1000 + 1000
    })
    closeCreateDialog()
  } catch (e) {
    console.error(e)
    alert('创建列失败')
  } finally {
    isCreating.value = false
  }
}

// 任务操作
const openCreateTaskDialog = (columnId: number) => {
  currentColumnId.value = columnId
  isCreateTaskDialogOpen.value = true
}

const closeCreateTaskDialog = () => {
  isCreateTaskDialogOpen.value = false
  currentColumnId.value = undefined
}

const handleCreateTask = async (data: { title: string; description: string; priority: number; columnId: number }) => {
  isCreatingTask.value = true
  try {
    await kanbanStore.createTask(data.columnId, {
      title: data.title,
      description: data.description,
      priority: data.priority
    })
    closeCreateTaskDialog()
  } catch (e) {
    console.error(e)
    alert('创建任务失败')
  } finally {
    isCreatingTask.value = false
  }
}
</script>
