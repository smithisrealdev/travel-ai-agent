<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="text-4xl">✈️</div>
            <div>
              <h1 class="text-3xl font-bold text-gray-900">Travel AI Agent</h1>
              <p class="text-sm text-gray-600">AI-powered travel planning assistant</p>
            </div>
          </div>
          <button
            v-if="hasSearchResult"
            @click="clearSearch"
            class="btn-secondary"
          >
            New Search
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- Hero Section (only show when no search result) -->
      <div v-if="!hasSearchResult" class="text-center mb-12">
        <h2 class="text-5xl font-extrabold text-gray-900 mb-4">
          Discover Your Next Adventure
        </h2>
        <p class="text-xl text-gray-600 max-w-2xl mx-auto">
          Let our AI-powered assistant help you plan the perfect trip with personalized
          recommendations, weather forecasts, and travel insights.
        </p>
      </div>

      <!-- Features Grid (only show when no search result) -->
      <div v-if="!hasSearchResult" class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div class="card text-center">
          <div class="text-4xl mb-4">🤖</div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">AI-Powered</h3>
          <p class="text-sm text-gray-600">
            Get personalized recommendations based on your preferences and budget
          </p>
        </div>
        <div class="card text-center">
          <div class="text-4xl mb-4">🌤️</div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Live Weather</h3>
          <p class="text-sm text-gray-600">
            Real-time weather information to help you plan the perfect trip
          </p>
        </div>
        <div class="card text-center">
          <div class="text-4xl mb-4">💰</div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Budget Planning</h3>
          <p class="text-sm text-gray-600">
            Get cost estimates and stay within your travel budget
          </p>
        </div>
      </div>

      <!-- Search Component -->
      <TravelSearch />

      <!-- Connection Status -->
      <div class="mt-8 text-center">
        <button
          @click="checkConnection"
          class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
        >
          <span v-if="connectionChecking">Checking connection...</span>
          <span v-else-if="connectionStatus === 'connected'" class="text-green-600">
            ✓ Connected to backend
          </span>
          <span v-else-if="connectionStatus === 'error'" class="text-red-600">
            ✗ Connection error - Click to retry
          </span>
          <span v-else>Click to check connection</span>
        </button>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center text-gray-600">
          <p class="text-sm">
            &copy; {{ new Date().getFullYear() }} Travel AI Agent. 
            Powered by OpenAI, Weather API, and Flight API.
          </p>
          <p class="text-xs mt-2">
            Built with Nuxt 3, Go Fiber, PostgreSQL, and Redis
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTravelStore } from '~/stores/travel'

const travelStore = useTravelStore()

// Connection status
const connectionStatus = ref<'idle' | 'connected' | 'error'>('idle')
const connectionChecking = ref(false)

// Computed properties
const hasSearchResult = computed(() => travelStore.hasSearchResult)

// Check backend connection
const checkConnection = async () => {
  connectionChecking.value = true
  try {
    const api = useApi()
    await api.checkHealth()
    connectionStatus.value = 'connected'
  } catch (err) {
    connectionStatus.value = 'error'
    console.error('Connection check failed:', err)
  } finally {
    connectionChecking.value = false
  }
}

// Clear current search
const clearSearch = () => {
  travelStore.clearCurrentSearch()
}

// Set page metadata
useHead({
  title: 'Travel AI Agent - Plan Your Perfect Trip',
  meta: [
    {
      name: 'description',
      content: 'AI-powered travel planning assistant with personalized recommendations, weather forecasts, and budget planning.',
    },
  ],
})
</script>
