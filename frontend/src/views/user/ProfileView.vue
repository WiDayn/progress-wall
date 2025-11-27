<template>
  <div class="min-h-screen bg-background">
    <div class="container mx-auto px-4 py-8">
      <div class="max-w-4xl mx-auto">
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-foreground">个人资料</h1>
          <p class="text-muted-foreground mt-2">管理您的个人信息和账户设置</p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- 个人信息卡片 -->
          <div class="lg:col-span-1">
            <Card class="p-6">
              <div class="text-center">
                <div class="relative inline-block group">
                  <Avatar
                    :src="avatarUrl"
                    :fallback="userStore.currentUser?.nickname?.charAt(0).toUpperCase() || userStore.currentUser?.username?.charAt(0).toUpperCase() || 'U'"
                    class="h-24 w-24 mx-auto mb-4 cursor-pointer hover:opacity-80 transition-opacity"
                    @click="triggerUpload"
                  />
                  <div class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 rounded-full opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
                     <span class="text-white text-xs">更换头像</span>
                  </div>
                </div>
                
                <input
                  type="file"
                  ref="fileInput"
                  class="hidden"
                  accept="image/jpeg,image/png,image/gif"
                  @change="handleFileChange"
                />

                <h2 class="text-xl font-semibold text-foreground">
                  {{ userStore.currentUser?.nickname || userStore.currentUser?.username }}
                </h2>
                <p class="text-muted-foreground">
                  {{ userStore.currentUser?.email }}
                </p>
                <div class="mt-4">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      userStore.currentUser?.role === 'admin'
                        ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                        : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                    ]"
                  >
                    {{ userStore.currentUser?.role === 'admin' ? '管理员' : '普通用户' }}
                  </span>
                </div>
              </div>
            </Card>
          </div>

          <!-- 详细信息 -->
          <div class="lg:col-span-2">
            <Card class="p-6">
              <CardHeader>
                <h3 class="text-lg font-semibold">基本信息</h3>
              </CardHeader>
              <CardContent class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-foreground mb-2">
                      昵称
                    </label>
                    <input
                      v-model="profileForm.nickname"
                      type="text"
                      class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-foreground mb-2">
                      邮箱
                    </label>
                    <input
                      v-model="profileForm.email"
                      type="email"
                      class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                    />
                  </div>
                   <div>
                    <label class="block text-sm font-medium text-foreground mb-2">
                      手机号
                    </label>
                    <input
                      v-model="profileForm.phone"
                      type="tel"
                      class="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                    />
                  </div>
                </div>
                
                <div class="flex justify-end space-x-4 pt-4">
                  <Button @click="resetForm" variant="outline">
                    重置
                  </Button>
                  <Button @click="saveProfile">
                    保存更改
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { getAvatarUrl } from '@/lib/utils'
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardContent from '@/components/ui/CardContent.vue'
import Button from '@/components/ui/Button.vue'
import Avatar from '@/components/ui/Avatar.vue'

const userStore = useUserStore()
const fileInput = ref<HTMLInputElement | null>(null)

const profileForm = ref({
  nickname: '',
  email: '',
  phone: ''
})

const avatarUrl = computed(() => {
  return getAvatarUrl(userStore.currentUser?.avatar)
})

const resetForm = () => {
  if (userStore.currentUser) {
    profileForm.value = {
      nickname: userStore.currentUser.nickname || '',
      email: userStore.currentUser.email || '',
      phone: userStore.currentUser.phone || ''
    }
  }
}

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileChange = async (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files && input.files[0]) {
    try {
      await userStore.uploadAvatar(input.files[0])
    } catch (e: any) {
      console.error(e)
      alert(e.message || '上传失败')
    } finally {
      // Reset input so same file can be selected again if needed
      input.value = ''
    }
  }
}

const saveProfile = async () => {
  try {
    await userStore.updateProfile({
      nickname: profileForm.value.nickname,
      email: profileForm.value.email,
      phone: profileForm.value.phone
    })
    alert('保存成功')
  } catch (e: any) {
    alert(e.message || '保存失败')
  }
}

onMounted(() => {
  resetForm()
})
</script>
