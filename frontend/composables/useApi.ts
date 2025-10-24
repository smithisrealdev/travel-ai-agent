import type { TravelSearchRequest, TravelSearchResponse, TravelSearch } from '~/types'

/**
 * API composable for making requests to the backend
 */
export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  /**
   * Search for travel recommendations
   */
  const searchTravel = async (request: TravelSearchRequest): Promise<TravelSearchResponse> => {
    try {
      const response = await $fetch<TravelSearchResponse>(`${baseURL}/api/v1/travel/search`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: request,
      })
      return response
    } catch (error) {
      console.error('Failed to search travel:', error)
      throw error
    }
  }

  /**
   * Get travel search history for a user
   */
  const getSearchHistory = async (userId: string): Promise<TravelSearch[]> => {
    try {
      const response = await $fetch<TravelSearch[]>(`${baseURL}/api/v1/travel/history`, {
        method: 'GET',
        params: { userId },
      })
      return response
    } catch (error) {
      console.error('Failed to get search history:', error)
      throw error
    }
  }

  /**
   * Check API health status
   */
  const checkHealth = async (): Promise<any> => {
    try {
      const response = await $fetch(`${baseURL}/health`)
      return response
    } catch (error) {
      console.error('Health check failed:', error)
      throw error
    }
  }

  return {
    searchTravel,
    getSearchHistory,
    checkHealth,
  }
}
