<template>
  <div class="card">
    <!-- Header -->
    <button
      @click="isExpanded = !isExpanded"
      class="w-full flex items-center justify-between cursor-pointer"
    >
      <div class="flex items-center gap-3">
        <span class="text-2xl">📅</span>
        <div class="text-left">
          <h3 class="text-xl font-poppins font-semibold text-gray-900">
            Day {{ day }}
          </h3>
          <p class="text-sm text-gray-500">{{ date }}</p>
        </div>
      </div>
      <div class="flex items-center gap-4">
        <span class="text-lg font-semibold text-primary">
          ฿{{ formatPrice(budget) }}
        </span>
        <svg
          :class="{ 'rotate-180': isExpanded }"
          class="w-6 h-6 transition-transform"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </button>

    <!-- Expandable Content -->
    <div v-show="isExpanded" class="mt-6 space-y-6">
      <!-- Hotel Section -->
      <div v-if="hotel" class="border-l-4 border-primary pl-4">
        <div class="flex items-start gap-4">
          <span class="text-3xl">🏨</span>
          <div class="flex-1">
            <div class="flex items-start justify-between">
              <div>
                <h4 class="font-poppins font-semibold text-lg text-gray-900">{{ hotel.name }}</h4>
                <div class="flex items-center gap-2 mt-1">
                  <span class="text-yellow-500">★</span>
                  <span class="text-sm text-gray-600">{{ hotel.rating }}</span>
                  <span class="text-sm text-gray-400">•</span>
                  <span class="text-sm font-semibold text-primary">฿{{ formatPrice(hotel.price) }}</span>
                </div>
              </div>
              <button class="px-4 py-2 bg-primary text-white rounded-lg text-sm font-semibold hover:bg-blue-600 transition">
                Book Now
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Flight Section -->
      <div v-if="flight" class="border-l-4 border-accent pl-4">
        <div class="flex items-start gap-4">
          <span class="text-3xl">✈️</span>
          <div class="flex-1">
            <div class="flex items-start justify-between">
              <div>
                <h4 class="font-poppins font-semibold text-lg text-gray-900">{{ flight.airline }}</h4>
                <p class="text-sm text-gray-600 mt-1">{{ flight.route }}</p>
                <span class="text-sm font-semibold text-primary mt-1 inline-block">฿{{ formatPrice(flight.price) }}</span>
              </div>
              <button class="px-4 py-2 bg-accent text-white rounded-lg text-sm font-semibold hover:bg-yellow-600 transition">
                Book Now
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Activities Section -->
      <div v-if="activities && activities.length > 0" class="border-l-4 border-gray-300 pl-4">
        <div class="flex items-start gap-4">
          <span class="text-3xl">🎯</span>
          <div class="flex-1">
            <h4 class="font-poppins font-semibold text-lg text-gray-900 mb-3">Activities</h4>
            <ul class="space-y-2">
              <li v-for="(activity, idx) in activities" :key="idx" class="flex items-center gap-2">
                <input type="checkbox" class="rounded text-primary focus:ring-primary" />
                <span class="text-gray-700">{{ activity }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Hotel {
  name: string
  price: number
  rating: number
  image?: string
}

interface Flight {
  airline: string
  route: string
  price: number
}

defineProps<{
  day: number
  date?: string
  budget: number
  hotel?: Hotel
  flight?: Flight
  activities?: string[]
}>()

const isExpanded = ref(false)

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US').format(price)
}
</script>
