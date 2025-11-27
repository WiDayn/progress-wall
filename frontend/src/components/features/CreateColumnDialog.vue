<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-background rounded-lg p-6 w-full max-w-sm mx-4">
      <h2 class="text-xl font-semibold mb-4">创建新列</h2>
      <form @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-foreground mb-2">列名称</label>
            <Input
              v-model="form.name"
              type="text"
              required
              placeholder="输入列名称"
              class="w-full"
              autoFocus
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-foreground mb-2">描述（可选）</label>
            <textarea
              v-model="form.description"
              class="w-full px-3 py-2 border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-ring bg-background"
              rows="3"
              placeholder="输入描述"
            ></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-foreground mb-2">颜色</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="color in presetColors"
                :key="color"
                type="button"
                @click="form.color = color"
                class="w-6 h-6 rounded-full border-2 transition-all hover:scale-110"
                :class="form.color === color ? 'border-foreground shadow-md' : 'border-transparent hover:border-muted-foreground'"
                :style="{ backgroundColor: color }"
                :title="color"
              ></button>
            </div>
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
  loading?: boolean
}

interface Emits {
  (e: 'submit', data: { name: string; description: string; color: string }): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

const form = reactive({
  name: '',
  description: '',
  color: '#9CA3AF' // 默认灰色
})

const presetColors = [
  '#9CA3AF', // 灰色
  '#3B82F6', // 蓝色
  '#10B981', // 绿色
  '#EF4444', // 红色
  '#F59E0B', // 黄色
  '#8B5CF6', // 紫色
]

const handleSubmit = () => {
  if (form.name.trim()) {
    emit('submit', {
      name: form.name.trim(),
      description: form.description.trim(),
      color: form.color
    })
  }
}

const handleCancel = () => {
  emit('cancel')
}

// 当对话框关闭时重置表单
watch(() => props.isOpen, (isOpen) => {
  if (!isOpen) {
    form.name = ''
    form.description = ''
    form.color = '#9CA3AF'
  }
})
</script>

