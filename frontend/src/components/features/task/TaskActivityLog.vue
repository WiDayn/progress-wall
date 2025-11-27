<template>
  <div class="w-96 border-l bg-muted/5 flex flex-col min-h-0">
    <div class="px-6 py-4 border-b bg-muted/10 flex items-center justify-between">
      <h3 class="font-semibold text-foreground">Activity</h3>
      <span class="text-xs text-muted-foreground">{{ activities.length }} 条记录</span>
    </div>
    
    <div class="flex-1 overflow-y-auto p-6">
      <div v-if="loading" class="flex justify-center py-4">
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
</template>

<script setup lang="ts">
import type { ActivityLog } from '@/services/task-api'
import Avatar from '@/components/ui/Avatar.vue'
import { getAvatarUrl } from '@/lib/utils'

defineProps<{
  activities: ActivityLog[]
  loading: boolean
}>()

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
</script>
