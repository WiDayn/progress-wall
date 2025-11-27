<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-background rounded-lg p-6 w-full max-w-md mx-4">
      <h2 class="text-xl font-semibold mb-4">创建新任务</h2>
      <form @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-foreground mb-2">任务标题</label>
            <Input
              v-model="form.title"
              type="text"
              required
              placeholder="输入任务标题"
              class="w-full"
              autoFocus
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-foreground mb-2">描述</label>
            <textarea
              v-model="form.description"
              class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring bg-background"
              rows="3"
              placeholder="输入任务描述（可选）"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-foreground mb-2">优先级</label>
              <select
                v-model="form.priority"
                class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring bg-background"
              >
                <option :value="1">低</option>
                <option :value="2">中</option>
                <option :value="3">高</option>
                <option :value="4">紧急</option>
              </select>
            </div>
            
            <!-- 未来可以添加负责人、截止日期等 -->
          </div>
        </div>

        <div class="flex justify-end space-x-2 mt-6">
          <Button type="button" variant="outline" @click="handleCancel">
            取消
          </Button>
          <Button type="submit" :disabled="loading">
            {{ loading ? '创建中...' : '创建' }}
          </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

interface Props {
  isOpen: boolean
  columnId?: number
  loading?: boolean
}

interface CreateTaskData {
  title: string
  description: string
  priority: number
  columnId: number
}

interface Emits {
  (e: 'submit', data: CreateTaskData): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

const form = reactive({
  title: '',
  description: '',
  priority: 2 // 默认中优先级
})

const handleSubmit = () => {
  if (form.title.trim() && props.columnId) {
    emit('submit', {
      title: form.title.trim(),
      description: form.description.trim(),
      priority: form.priority,
      columnId: props.columnId
    })
  }
}

const handleCancel = () => {
  emit('cancel')
}

// 当对话框关闭时重置表单
watch(() => props.isOpen, (isOpen) => {
  if (!isOpen) {
    form.title = ''
    form.description = ''
    form.priority = 2
  }
})
</script>

