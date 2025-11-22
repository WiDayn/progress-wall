<template>
  <input
    :class="cn(
      'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
      typeof $attrs.class === 'string' ? $attrs.class : undefined
    )"
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    v-bind="inputAttrs"
  />
</template>

<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from "@/lib/utils"

interface Props {
  modelValue?: string | number
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const attrs = useAttrs()
const inputAttrs = computed(() => {
  const { class: _, ...rest } = attrs as Record<string, any>
  return rest
})
</script>

