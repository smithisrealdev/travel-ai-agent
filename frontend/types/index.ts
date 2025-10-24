/**
 * Type definitions for the Travel AI Agent application
 */

export interface TravelSearchRequest {
  destination: string
  startDate?: string
  endDate?: string
  budget?: number
  preferences?: Record<string, any>
  userId?: string
}

export interface TravelSearchResponse {
  searchId: number
  destination: string
  summary: string
  recommendations: TravelRecommendation[]
  weather?: WeatherInfo
  flights?: FlightInfo[]
  estimatedCost: number
  createdAt: string
}

export interface TravelRecommendation {
  id: number
  type: string
  title: string
  description: string
  price?: number
  rating?: number
  location?: string
  imageUrl?: string
  url?: string
  metadata?: Record<string, any>
}

export interface WeatherInfo {
  temperature: number
  description: string
  humidity: number
  windSpeed: number
  icon: string
  forecast?: DayForecast[]
}

export interface DayForecast {
  date: string
  tempMin: number
  tempMax: number
  description: string
  icon: string
}

export interface FlightInfo {
  flightNumber: string
  airline: string
  departure: string
  arrival: string
  departTime: string
  arriveTime: string
  duration: string
  price: number
  stops: number
}

export interface TravelSearch {
  id: number
  userId: string
  destination: string
  startDate?: string
  endDate?: string
  budget: number
  preferences: Record<string, any>
  results: Record<string, any>
  createdAt: string
  updatedAt: string
}

export interface HealthCheckResponse {
  status: string
  services: Record<string, string>
  time: string
}
