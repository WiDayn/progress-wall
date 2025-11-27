<template>
  <div class="bg-muted/50 rounded-lg p-4 min-h-[500px] w-80 flex-shrink-0 flex flex-col" :data-column-id="column.id">
    <div class="flex justify-between items-center mb-4 flex-none">
      <h3 class="text-lg font-semibold truncate mr-2">
        {{ column.name || column.title }}
      </h3>
      <span class="text-sm text-muted-foreground bg-muted px-2 py-1 rounded-full flex-none">
        {{ column.tasks?.length || 0 }}
      </span>
    </div>
    
    <div class="flex-1 min-h-0 overflow-y-auto pr-1 -mr-1 scrollbar-thin">
      <VueDraggable
        v-model="column.tasks"
        :group="{ name: 'tasks', pull: true, put: true }"
        :animation="200"
        ghost-class="ghost-task"
        chosen-class="chosen-task"
        drag-class="drag-task"
        class="space-y-3 min-h-[100px]"
        @end="onDragEnd"
      >
        <TaskCard
          v-for="task in column.tasks"
          :key="task.id"
          :task="task"
          :data-task-id="task.id"
          @select="$emit('select-task', task)"
          @delete="$emit('delete-task', task.id)"
        />
        
        <div
          v-if="!column.tasks || column.tasks.length === 0"
          class="text-center text-muted-foreground py-8 pointer-events-none"
        >
          暂无任务
        </div>
      </VueDraggable>
    </div>

    <!-- 底部添加按钮 -->
    <button 
      @click="$emit('add-task', column.id)"
      class="mt-2 w-full py-2 flex items-center justify-center text-sm text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors border border-transparent hover:border-border border-dashed flex-none"
    >
      + 添加任务
    </button>
  </div>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import type { Column } from '@/stores/kanban'
import { useKanbanStore } from '@/stores/kanban'
import TaskCard from '@/components/features/TaskCard.vue'

interface Props {
  column: Column
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'select-task': [task: any]
  'delete-task': [taskId: number]
  'add-task': [columnId: number]
}>()

const kanbanStore = useKanbanStore()

// 拖拽结束事件处理
const onDragEnd = async (event: any) => {
  const { item, to, from, newIndex, oldIndex } = event
  
  // 获取被拖拽的任务ID
  const taskIdStr = item.getAttribute('data-task-id')
  if (!taskIdStr) {
    console.error('未找到 taskId')
    return
  }
  const taskId = parseInt(taskIdStr, 10)

  // 获取目标列ID
  const targetColumnElement = to.closest('[data-column-id]')
  const targetColumnIdStr = targetColumnElement?.getAttribute('data-column-id')
  if (!targetColumnIdStr) {
    console.error('未找到 targetColumnId')
    return
  }
  const targetColumnId = parseInt(targetColumnIdStr, 10)

  // 如果是同一列内的重新排序，或者跨列移动
  if (from !== to || newIndex !== oldIndex) {
    console.log('拖拽移动 - taskId:', taskId, 'to column:', targetColumnId, 'index:', newIndex)
    try {
      // 在 store 中已经处理了 API 调用和乐观更新
      // 如果 store 中的同步逻辑出错（比如找不到列），这里会捕获到异常
      // 但异步的 API 错误会在 store 内部捕获并处理回滚，这里不需要 await
      await kanbanStore.moveTaskWithDrag(taskId, targetColumnId, newIndex)
    } catch (error) {
      console.error('拖拽移动任务失败:', error)
      
      // 如果是同步逻辑出错，手动回滚 UI（VueDraggable 会自动修改 v-model，需要撤销）
      // 注意：VueDraggable 的 v-model 双向绑定已经修改了数据
      // 如果 store 抛出错误，说明 store 的数据没有更新或者更新失败
      // 但因为 v-model 是直接绑定到 column.tasks 的，VueDraggable 可能已经修改了数组
      // 最简单的回滚方式是重新获取看板数据，或者利用 store 的回滚机制
      
      // 由于 store.moveTaskWithDrag 内部已经有 try-catch 处理同步错误并回滚
      // 这里其实主要捕获的是 store 抛出的同步错误
      // 我们可以选择刷新看板数据来确保一致性
      // kanbanStore.fetchBoardDetail(kanbanStore.currentBoardId!)
    }
  }
}
</script>

<style scoped>
/* 拖拽样式 */
:deep(.ghost-task) {
  opacity: 0.5;
  background: #f0f0f0;
  border: 2px dashed #ccc;
}

:deep(.chosen-task) {
  transform: rotate(2deg);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

:deep(.drag-task) {
  transform: rotate(2deg);
  opacity: 0.9;
}
</style>
