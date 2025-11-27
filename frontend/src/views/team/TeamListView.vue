<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { teamApi, type Team } from '@/services/team-api'
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardContent from '@/components/ui/CardContent.vue'
import Button from '@/components/ui/Button.vue'
import { Plus, Users } from 'lucide-vue-next'

const router = useRouter()
const teams = ref<Team[]>([])
const loading = ref(true)
const showCreateModal = ref(false)
const newTeam = ref({ name: '', description: '' })

const fetchTeams = async () => {
  try {
    loading.value = true
    const res = await teamApi.getMyTeams()
    teams.value = res.data.teams
  } catch (err) {
    console.error('Failed to fetch teams', err)
  } finally {
    loading.value = false
  }
}

const handleCreateTeam = async () => {
  if (!newTeam.value.name) return
  try {
    await teamApi.createTeam(newTeam.value)
    showCreateModal.value = false
    newTeam.value = { name: '', description: '' }
    fetchTeams()
  } catch (err) {
    console.error('Failed to create team', err)
  }
}

const navigateToTeam = (teamId: number) => {
  router.push(`/teams/${teamId}/projects`)
}

onMounted(() => {
  fetchTeams()
})
</script>

<template>
  <div class="container mx-auto p-6">
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">My Teams</h1>
        <p class="text-gray-500 mt-1">Select a team to view projects</p>
      </div>
      <Button @click="showCreateModal = true">
        <Plus class="w-4 h-4 mr-2" />
        New Team
      </Button>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p>Loading teams...</p>
    </div>

    <div v-else-if="teams.length === 0" class="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
      <Users class="w-12 h-12 mx-auto text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900">No teams yet</h3>
      <p class="text-gray-500 mt-2 mb-6">Create your first team to get started</p>
      <Button @click="showCreateModal = true">Create Team</Button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card 
        v-for="team in teams" 
        :key="team.id" 
        class="cursor-pointer hover:shadow-lg transition-shadow"
        @click="navigateToTeam(team.id)"
      >
        <CardHeader>
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg">{{ team.name }}</h3>
            <Users class="w-4 h-4 text-gray-400" />
          </div>
        </CardHeader>
        <CardContent>
          <p class="text-gray-500 text-sm line-clamp-2">
            {{ team.description || 'No description provided' }}
          </p>
          <div class="mt-4 text-xs text-gray-400">
            Created by user #{{ team.creator_id }}
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Create Team Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Create New Team</h2>
        <form @submit.prevent="handleCreateTeam">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1">Team Name</label>
              <input 
                v-model="newTeam.name" 
                type="text" 
                class="w-full border rounded-md p-2 dark:bg-gray-700 dark:border-gray-600" 
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Description</label>
              <textarea 
                v-model="newTeam.description" 
                class="w-full border rounded-md p-2 dark:bg-gray-700 dark:border-gray-600" 
                rows="3"
              ></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <Button type="button" variant="ghost" @click="showCreateModal = false">Cancel</Button>
            <Button type="submit">Create Team</Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

