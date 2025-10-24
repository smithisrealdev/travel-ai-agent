<template>
  <div class="max-w-4xl mx-auto">
    <!-- Search Form -->
    <div class="card mb-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-6">Plan Your Next Adventure</h2>
      
      <form @submit.prevent="handleSearch" class="space-y-6">
        <!-- Destination Input -->
        <div>
          <label for="destination" class="block text-sm font-medium text-gray-700 mb-2">
            Where would you like to go?
          </label>
          <input
            id="destination"
            v-model="searchForm.destination"
            type="text"
            placeholder="e.g., Paris, Tokyo, New York"
            class="input-field"
            required
          />
        </div>

        <!-- Date Range -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label for="startDate" class="block text-sm font-medium text-gray-700 mb-2">
              Start Date
            </label>
            <input
              id="startDate"
              v-model="searchForm.startDate"
              type="date"
              class="input-field"
            />
          </div>
          <div>
            <label for="endDate" class="block text-sm font-medium text-gray-700 mb-2">
              End Date
            </label>
            <input
              id="endDate"
              v-model="searchForm.endDate"
              type="date"
              class="input-field"
            />
          </div>
        </div>

        <!-- Budget -->
        <div>
          <label for="budget" class="block text-sm font-medium text-gray-700 mb-2">
            Budget (USD)
          </label>
          <input
            id="budget"
            v-model.number="searchForm.budget"
            type="number"
            min="0"
            step="100"
            placeholder="e.g., 2000"
            class="input-field"
          />
        </div>

        <!-- Preferences -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Travel Preferences
          </label>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
            <label class="flex items-center space-x-2 cursor-pointer">
              <input
                v-model="preferences.culture"
                type="checkbox"
                class="rounded text-primary-600 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-700">Culture</span>
            </label>
            <label class="flex items-center space-x-2 cursor-pointer">
              <input
                v-model="preferences.adventure"
                type="checkbox"
                class="rounded text-primary-600 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-700">Adventure</span>
            </label>
            <label class="flex items-center space-x-2 cursor-pointer">
              <input
                v-model="preferences.relaxation"
                type="checkbox"
                class="rounded text-primary-600 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-700">Relaxation</span>
            </label>
            <label class="flex items-center space-x-2 cursor-pointer">
              <input
                v-model="preferences.food"
                type="checkbox"
                class="rounded text-primary-600 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-700">Food</span>
            </label>
          </div>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="loading"
          class="btn-primary w-full"
        >
          <span v-if="loading">Searching...</span>
          <span v-else>Search Destinations</span>
        </button>
      </form>

      <!-- Error Message -->
      <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
        <p class="text-red-800 text-sm">{{ error }}</p>
      </div>
    </div>

    <!-- Search Results -->
    <div v-if="searchResult" class="space-y-6">
      <!-- Weather Information -->
      <div v-if="searchResult.weather" class="card">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Current Weather</h3>
        <div class="flex items-center space-x-4">
          <div class="text-4xl">
            <img
              v-if="searchResult.weather.icon"
              :src="`https://openweathermap.org/img/wn/${searchResult.weather.icon}@2x.png`"
              :alt="searchResult.weather.description"
              class="w-16 h-16"
            />
          </div>
          <div>
            <p class="text-3xl font-bold text-gray-900">
              {{ Math.round(searchResult.weather.temperature) }}°C
            </p>
            <p class="text-gray-600 capitalize">{{ searchResult.weather.description }}</p>
            <p class="text-sm text-gray-500">
              Humidity: {{ searchResult.weather.humidity }}% | 
              Wind: {{ searchResult.weather.windSpeed }} m/s
            </p>
          </div>
        </div>
      </div>

      <!-- AI Recommendations -->
      <div v-if="searchResult.summary" class="card">
        <h3 class="text-xl font-bold text-gray-900 mb-4">AI Travel Recommendations</h3>
        <div class="prose prose-sm max-w-none text-gray-700">
          <p class="whitespace-pre-line">{{ searchResult.summary }}</p>
        </div>
      </div>

      <!-- Recommendations Grid -->
      <div v-if="searchResult.recommendations.length > 0" class="card">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Suggested Activities & Places</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="rec in searchResult.recommendations"
            :key="rec.id"
            class="border border-gray-200 rounded-lg p-4 hover:border-primary-500 transition-colors"
          >
            <div class="flex items-start justify-between mb-2">
              <span
                class="inline-block px-2 py-1 text-xs font-semibold text-primary-700 bg-primary-100 rounded"
              >
                {{ rec.type }}
              </span>
              <span v-if="rec.rating" class="text-sm font-semibold text-yellow-600">
                ⭐ {{ rec.rating }}
              </span>
            </div>
            <h4 class="font-semibold text-gray-900 mb-2">{{ rec.title }}</h4>
            <p class="text-sm text-gray-600 mb-2">{{ rec.description }}</p>
            <p v-if="rec.price" class="text-lg font-bold text-primary-600">
              ${{ rec.price.toFixed(2) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Estimated Cost -->
      <div class="card bg-primary-50 border border-primary-200">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Estimated Trip Cost</h3>
            <p class="text-sm text-gray-600">Based on your preferences</p>
          </div>
          <p class="text-3xl font-bold text-primary-600">
            ${{ searchResult.estimatedCost.toFixed(2) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTravelStore } from '~/stores/travel'
import type { TravelSearchRequest } from '~/types'

const travelStore = useTravelStore()

// Form state
const searchForm = ref({
  destination: '',
  startDate: '',
  endDate: '',
  budget: undefined as number | undefined,
})

const preferences = ref({
  culture: false,
  adventure: false,
  relaxation: false,
  food: false,
})

// Computed properties
const loading = computed(() => travelStore.loading)
const error = computed(() => travelStore.error)
const searchResult = computed(() => travelStore.currentSearch)

// Handle search submission
const handleSearch = async () => {
  const request: TravelSearchRequest = {
    destination: searchForm.value.destination,
    startDate: searchForm.value.startDate || undefined,
    endDate: searchForm.value.endDate || undefined,
    budget: searchForm.value.budget || undefined,
    preferences: {
      culture: preferences.value.culture,
      adventure: preferences.value.adventure,
      relaxation: preferences.value.relaxation,
      food: preferences.value.food,
    },
  }

  try {
    await travelStore.searchTravel(request)
  } catch (err) {
    console.error('Search failed:', err)
  }
}
</script>
