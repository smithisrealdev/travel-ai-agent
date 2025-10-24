<template>
  <div>
    <NuxtLayout name="default">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <!-- Back Button -->
        <button @click="$router.back()" class="flex items-center gap-2 text-gray-600 hover:text-primary mb-6 transition">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          <span class="font-medium">Back to Chat</span>
        </button>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- Main Content -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Trip Overview Card -->
            <div class="card">
              <div class="flex items-start justify-between mb-6">
                <div>
                  <h1 class="heading-2 mb-2">{{ trip.destination }}</h1>
                  <div class="flex items-center gap-4 text-gray-600">
                    <span class="flex items-center gap-1">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                      {{ trip.duration }}
                    </span>
                    <span class="flex items-center gap-1">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      {{ trip.budget }}
                    </span>
                  </div>
                </div>
                <div class="text-right">
                  <div class="flex items-center gap-2 text-gray-600">
                    <span class="text-2xl">🌤️</span>
                    <div>
                      <p class="font-semibold text-gray-900">{{ trip.weather.condition }}</p>
                      <p class="text-sm">{{ trip.weather.temp }}</p>
                    </div>
                  </div>
                </div>
              </div>

              <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
                <p class="text-sm text-blue-800">
                  <strong>Best time to visit:</strong> {{ trip.bestTime }}
                </p>
              </div>
            </div>

            <!-- Itinerary -->
            <div>
              <h2 class="text-2xl font-poppins font-semibold text-gray-900 mb-4">Suggested Itinerary</h2>
              <div class="space-y-4">
                <ItineraryCard
                  v-for="day in trip.itinerary"
                  :key="day.day"
                  :day="day.day"
                  :date="day.date"
                  :budget="day.budget"
                  :hotel="day.hotel"
                  :flight="day.flight"
                  :activities="day.activities"
                />
              </div>
            </div>
          </div>

          <!-- Sidebar -->
          <aside class="lg:col-span-1">
            <div class="sticky top-24 space-y-6">
              <!-- Price Breakdown -->
              <div class="card">
                <h3 class="text-xl font-poppins font-semibold text-gray-900 mb-4">Price Breakdown</h3>
                
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600">Flights</span>
                    <span class="font-semibold text-gray-900">฿{{ formatPrice(priceBreakdown.flights) }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600">Hotels</span>
                    <span class="font-semibold text-gray-900">฿{{ formatPrice(priceBreakdown.hotels) }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600">Activities</span>
                    <span class="font-semibold text-gray-900">฿{{ formatPrice(priceBreakdown.activities) }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600">Food & Misc</span>
                    <span class="font-semibold text-gray-900">฿{{ formatPrice(priceBreakdown.food) }}</span>
                  </div>
                  
                  <div class="border-t border-gray-200 pt-3">
                    <div class="flex items-center justify-between">
                      <span class="font-poppins font-semibold text-lg text-gray-900">Total</span>
                      <span class="font-poppins font-semibold text-2xl text-primary">
                        ฿{{ formatPrice(priceBreakdown.total) }}
                      </span>
                    </div>
                  </div>
                </div>

                <button class="w-full btn-primary mt-6">
                  Book Now
                </button>
              </div>

              <!-- Alternative Plan -->
              <div class="card bg-green-50 border border-green-200">
                <div class="flex items-start gap-3 mb-4">
                  <span class="text-3xl">💡</span>
                  <div>
                    <h4 class="font-poppins font-semibold text-gray-900 mb-1">Save More!</h4>
                    <p class="text-sm text-gray-600">Alternative budget plan available</p>
                  </div>
                </div>
                
                <div class="bg-white rounded-lg p-4 mb-4">
                  <div class="flex items-center justify-between mb-2">
                    <span class="text-sm text-gray-600">Cheaper hotels</span>
                    <span class="text-sm font-semibold text-green-600">-฿{{ formatPrice(alternativeSavings.hotels) }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">Different activities</span>
                    <span class="text-sm font-semibold text-green-600">-฿{{ formatPrice(alternativeSavings.activities) }}</span>
                  </div>
                  <div class="border-t border-gray-200 mt-3 pt-3">
                    <div class="flex items-center justify-between">
                      <span class="font-semibold text-gray-900">Total Savings</span>
                      <span class="font-semibold text-lg text-green-600">฿{{ formatPrice(alternativeSavings.total) }}</span>
                    </div>
                  </div>
                </div>

                <button class="w-full btn-secondary text-sm">
                  View Alternative Plan
                </button>
              </div>
            </div>
          </aside>
        </div>
      </div>
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const tripId = route.params.id

// Sample trip data
const trip = {
  id: tripId,
  destination: 'Chiang Mai, Thailand',
  duration: '7 days',
  budget: '50,000 THB',
  weather: {
    condition: 'Sunny',
    temp: '32°C'
  },
  bestTime: 'November to February offers the best weather with cooler temperatures and less rain.',
  itinerary: [
    {
      day: 1,
      date: '2025-11-01',
      budget: 6000,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5,
        image: '/hotels/riverside.jpg'
      },
      activities: [
        'Arrival in Chiang Mai',
        'Check-in at hotel',
        'Explore Old City Temples',
        'Visit Warorot Market',
        'Dinner at local restaurant'
      ]
    },
    {
      day: 2,
      date: '2025-11-02',
      budget: 5500,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5
      },
      activities: [
        'Visit Doi Suthep Temple',
        'Lunch at mountain view restaurant',
        'Explore Nimman neighborhood',
        'Evening at Night Bazaar'
      ]
    },
    {
      day: 3,
      date: '2025-11-03',
      budget: 7000,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5
      },
      activities: [
        'Full day at Elephant Nature Park',
        'Learn about elephant conservation',
        'Feed and bathe elephants',
        'Traditional Thai dinner'
      ]
    },
    {
      day: 4,
      date: '2025-11-04',
      budget: 6000,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5
      },
      activities: [
        'Thai cooking class',
        'Visit local market for ingredients',
        'Afternoon at Bua Thong Waterfalls',
        'Dinner at cooking class venue'
      ]
    },
    {
      day: 5,
      date: '2025-11-05',
      budget: 5500,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5
      },
      activities: [
        'White Temple (Wat Rong Khun) day trip',
        'Blue Temple (Wat Rong Suea Ten)',
        'Black House (Baan Dam Museum)',
        'Return to Chiang Mai'
      ]
    },
    {
      day: 6,
      date: '2025-11-06',
      budget: 6000,
      hotel: {
        name: 'Riverside Hotel Chiang Mai',
        price: 2500,
        rating: 4.5
      },
      activities: [
        'Visit Art in Paradise 3D Museum',
        'Shopping at Central Festival',
        'Thai massage at spa',
        'Farewell dinner at Khao Soi restaurant'
      ]
    },
    {
      day: 7,
      date: '2025-11-07',
      budget: 4000,
      flight: {
        airline: 'Thai Airways',
        route: 'Chiang Mai → Bangkok',
        price: 3000
      },
      activities: [
        'Last minute shopping',
        'Check-out from hotel',
        'Transfer to airport',
        'Departure flight'
      ]
    }
  ]
}

const priceBreakdown = {
  flights: 15000,
  hotels: 17500,
  activities: 12000,
  food: 5500,
  total: 50000
}

const alternativeSavings = {
  hotels: 5000,
  activities: 3000,
  total: 8000
}

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US').format(price)
}

// Set page metadata
useHead({
  title: `Booking - ${trip.destination} - Travel AI Agent`,
  meta: [
    {
      name: 'description',
      content: `Review and book your ${trip.duration} trip to ${trip.destination}`,
    },
  ],
})
</script>
