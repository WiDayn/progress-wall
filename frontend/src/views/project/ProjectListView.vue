<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { teamApi, type Team } from '@/services/team-api'
import api from '@/lib/api'
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardContent from '@/components/ui/CardContent.vue'
import Button from '@/components/ui/Button.vue'
import { Plus, Briefcase, ArrowLeft } from 'lucide-vue-next'

interface Project {
  id: number
  name: string
  description: string
  team_id: number
  created_at: string
}

const route = useRoute()
const router = useRouter()
const teamId = Number(route.params.teamId)

const team = ref<Team | null>(null)
const projects = ref<Project[]>([])
const loading = ref(true)
const showCreateModal = ref(false)
const newProject = ref({ name: '', description: '' })

const fetchData = async () => {
  try {
    loading.value = true
    // Parallel fetch
    const [teamRes, projectsRes] = await Promise.all([
      teamApi.getTeam(teamId),
      api.get<{ projects: Project[] }>(`/teams/${teamId}/projects`)
    ])
    team.value = teamRes.data
    projects.value = projectsRes.data.projects
  } catch (err) {
    console.error('Failed to fetch data', err)
  } finally {
    loading.value = false
  }
}

const handleCreateProject = async () => {
  if (!newProject.value.name) return
  try {
    await api.post(`/teams/${teamId}/projects`, newProject.value)
    showCreateModal.value = false
    newProject.value = { name: '', description: '' }
    fetchData()
  } catch (err) {
    console.error('Failed to create project', err)
  }
}

const navigateToBoard = (projectId: number) => {
  // Assuming project detail view lists boards, or directly to first board?
  // Requirement: "Click project -> Project View (lists boards)"
  router.push(`/projects/${projectId}/boards`)
}

onMounted(() => {
  if (!teamId) {
    router.push('/teams')
    return
  }
  fetchData()
})
</script>

<template>
  <div class="container mx-auto p-6">
    <!-- Breadcrumb / Header -->
    <div class="flex items-center gap-2 text-sm text-gray-500 mb-6">
      <router-link to="/teams" class="hover:text-blue-600">My Teams</router-link>
      <span>/</span>
      <span class="text-gray-900 font-medium">{{ team?.name || 'Loading...' }}</span>
    </div>

    <div class="flex justify-between items-center mb-8">
      <div class="flex items-center gap-4">
        <Button variant="ghost" size="sm" @click="router.push('/teams')">
          <ArrowLeft class="w-4 h-4" />
        </Button>
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ team?.name }}</h1>
          <p class="text-gray-500 mt-1">{{ team?.description }}</p>
        </div>
      </div>
      <Button @click="showCreateModal = true">
        <Plus class="w-4 h-4 mr-2" />
        New Project
      </Button>
    </div>

    <div v-if="loading" class="text-center py-12">
      <p>Loading projects...</p>
    </div>

    <div v-else-if="projects.length === 0" class="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
      <Briefcase class="w-12 h-12 mx-auto text-gray-400 mb-4" />
      <h3 class="text-lg font-medium text-gray-900">No projects yet</h3>
      <p class="text-gray-500 mt-2 mb-6">Start a new project in this team</p>
      <Button @click="showCreateModal = true">Create Project</Button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card 
        v-for="project in projects" 
        :key="project.id" 
        class="cursor-pointer hover:shadow-lg transition-shadow"
        @click="navigateToBoard(project.id)"
      >
        <CardHeader>
          <div class="flex items-center justify-between">
            <h3 class="font-semibold text-lg">{{ project.name }}</h3>
            <Briefcase class="w-4 h-4 text-gray-400" />
          </div>
        </CardHeader>
        <CardContent>
          <p class="text-gray-500 text-sm line-clamp-2">
            {{ project.description || 'No description provided' }}
          </p>
          <div class="mt-4 text-xs text-gray-400">
            Created at {{ new Date(project.created_at).toLocaleDateString() }}
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Create Project Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Create New Project</h2>
        <p class="text-sm text-gray-500 mb-4">In team: {{ team?.name }}</p>
        <form @submit.prevent="handleCreateProject">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1">Project Name</label>
              <input 
                v-model="newProject.name" 
                type="text" 
                class="w-full border rounded-md p-2 dark:bg-gray-700 dark:border-gray-600" 
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Description</label>
              <textarea 
                v-model="newProject.description" 
                class="w-full border rounded-md p-2 dark:bg-gray-700 dark:border-gray-600" 
                rows="3"
              ></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <Button type="button" variant="ghost" @click="showCreateModal = false">Cancel</Button>
            <Button type="submit">Create Project</Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

