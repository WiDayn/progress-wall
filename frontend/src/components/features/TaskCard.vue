<template>
  <Card
    class="p-4 cursor-pointer hover:shadow-md transition-all hover:scale-[1.02] bg-card text-card-foreground border-border"
    :data-task-id="task.id"
    @click="$emit('select', task)"
  >
    <div class="space-y-2">
      <div class="flex justify-between items-start">
        <h4 class="font-medium text-sm leading-tight">{{ task.title }}</h4>
        <Button
          @click.stop="$emit('delete', task.id)"
          variant="ghost"
          size="sm"
          class="h-6 w-6 p-0 text-muted-foreground hover:text-destructive ml-2"
        >
          ×
        </Button>
      </div>
      
      <div class="flex justify-between items-center pt-1">
        <span
          :class="[
            'px-2 py-0.5 rounded-full text-[10px] font-medium',
            getPriorityClass(task.priority)
          ]"
        >
          {{ getPriorityText(task.priority) }}
        </span>
        
        <span class="text-[10px] text-muted-foreground">
          {{ formatDate(task.updated_at || task.created_at) }}
        </span>
      </div>
    </div>
  </Card>
</template>

<script setup lang="ts">
import type { Task } from '@/stores/kanban'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

interface Props {
  task: Task
}

defineProps<Props>()

defineEmits<{
  select: [task: Task]
  delete: [taskId: number]
}>()

const getPriorityText = (priority: number) => {
  const priorityMap: Record<number, string> = {
    1: '低',
    2: '中',
    3: '高',
    4: '紧急'
  }
  return priorityMap[priority] || '未知'
}

const getPriorityClass = (priority: number) => {
  const map: Record<number, string> = {
    1: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
    2: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    3: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
    4: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
  }
  return map[priority] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    month: 'short',
    day: 'numeric'
  }).format(new Date(dateStr))
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
