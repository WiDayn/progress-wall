<template>
  <div class="flex-1 p-8 overflow-y-auto">
    <div class="space-y-8">
      <!-- Status & Priority -->
      <div class="flex flex-wrap gap-4 items-center p-3 bg-muted/10 rounded-lg border border-border/50">
        <!-- 当前列信息 -->
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground font-medium uppercase">当前列</span>
          <span class="px-2 py-0.5 rounded-full text-xs font-medium border bg-muted text-muted-foreground">
            {{ task.column?.name || '加载中...' }}
          </span>
        </div>

        <div class="h-4 w-px bg-border/50 mx-1"></div>

        <!-- 优先级选择 -->
         <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground font-medium uppercase">优先级</span>
          <select 
            :value="task.priority"
            @change="$emit('update:priority', Number(($event.target as HTMLSelectElement).value))"
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
          <button v-if="!isEditing" @click="$emit('start-editing')" class="text-xs text-primary hover:underline">编辑</button>
        </div>
        
        <div v-if="isEditing">
          <textarea
            :value="description"
            @input="$emit('update:description', ($event.target as HTMLTextAreaElement).value)"
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
</template>

<script setup lang="ts">
import type { Task } from '@/stores/kanban'
import Avatar from '@/components/ui/Avatar.vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

defineProps<{
  task: Task
  isEditing: boolean
  description: string
}>()

defineEmits<{
  (e: 'update:priority', value: number): void
  (e: 'update:description', value: string): void
  (e: 'start-editing'): void
}>()

const renderMarkdown = (text: string) => {
  const rawMarkup = marked.parse(text || '')
  return DOMPurify.sanitize(rawMarkup as string)
}

const formatDateTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
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
</script>
