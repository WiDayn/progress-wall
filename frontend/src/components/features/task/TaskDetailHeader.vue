<template>
  <div class="flex items-center justify-between px-6 py-4 border-b">
    <div class="flex items-center gap-3 flex-1">
      <span class="text-sm font-mono text-muted-foreground">#{{ taskId }}</span>
      
      <!-- 可编辑标题 -->
      <div v-if="isEditing" class="flex-1 max-w-xl">
        <input 
          :value="title"
          @input="$emit('update:title', ($event.target as HTMLInputElement).value)"
          class="w-full px-2 py-1 text-xl font-bold border rounded focus:outline-none focus:ring-2 focus:ring-primary"
          placeholder="任务标题"
          @keyup.enter="$emit('save')"
        />
      </div>
      <h2 v-else class="text-xl font-bold text-foreground truncate max-w-md cursor-pointer hover:text-primary transition-colors" :title="title" @click="$emit('start-editing')">
        {{ title || '加载中...' }}
      </h2>
    </div>

    <div class="flex items-center gap-2">
       <div v-if="isEditing" class="flex gap-2 mr-4">
        <button @click="$emit('save')" class="px-3 py-1 text-sm bg-primary text-primary-foreground rounded hover:bg-primary/90">保存</button>
        <button @click="$emit('cancel')" class="px-3 py-1 text-sm bg-secondary text-secondary-foreground rounded hover:bg-secondary/80">取消</button>
      </div>
      
      <!-- 删除按钮 -->
      <button 
        @click="$emit('delete')"
        class="text-muted-foreground hover:text-destructive p-1 rounded-md hover:bg-destructive/10 transition-colors mr-2"
        title="删除任务"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
      </button>

      <button @click="$emit('close')" class="text-muted-foreground hover:text-foreground p-1 rounded-md hover:bg-muted transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  taskId?: number
  title: string
  isEditing: boolean
}>()

defineEmits<{
  (e: 'update:title', value: string): void
  (e: 'save'): void
  (e: 'cancel'): void
  (e: 'start-editing'): void
  (e: 'delete'): void
  (e: 'close'): void
}>()
</script>
