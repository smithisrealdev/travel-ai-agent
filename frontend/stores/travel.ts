import { defineStore } from 'pinia'
import type { TravelSearchRequest, TravelSearchResponse, TravelSearch } from '~/types'

interface TravelState {
  currentSearch: TravelSearchResponse | null
  searchHistory: TravelSearch[]
  loading: boolean
  error: string | null
}

/**
 * Pinia store for managing travel search state
 */
export const useTravelStore = defineStore('travel', {
  state: (): TravelState => ({
    currentSearch: null,
    searchHistory: [],
    loading: false,
    error: null,
  }),

  getters: {
    /**
     * Check if there's an active search result
     */
    hasSearchResult: (state): boolean => state.currentSearch !== null,

    /**
     * Get the current destination
     */
    currentDestination: (state): string | null => 
      state.currentSearch?.destination || null,

    /**
     * Get weather information from current search
     */
    currentWeather: (state) => state.currentSearch?.weather || null,

    /**
     * Get recommendations from current search
     */
    currentRecommendations: (state) => 
      state.currentSearch?.recommendations || [],
  },

  actions: {
    /**
     * Search for travel recommendations
     */
    async searchTravel(request: TravelSearchRequest) {
      this.loading = true
      this.error = null

      try {
        const api = useApi()
        const response = await api.searchTravel(request)
        this.currentSearch = response
        return response
      } catch (err: any) {
        this.error = err.message || 'Failed to search travel'
        throw err
      } finally {
        this.loading = false
      }
    },

    /**
     * Load search history for a user
     */
    async loadSearchHistory(userId: string) {
      this.loading = true
      this.error = null

      try {
        const api = useApi()
        const history = await api.getSearchHistory(userId)
        this.searchHistory = history
        return history
      } catch (err: any) {
        this.error = err.message || 'Failed to load search history'
        throw err
      } finally {
        this.loading = false
      }
    },

    /**
     * Clear the current search
     */
    clearCurrentSearch() {
      this.currentSearch = null
      this.error = null
    },

    /**
     * Clear the error message
     */
    clearError() {
      this.error = null
    },
  },
})
