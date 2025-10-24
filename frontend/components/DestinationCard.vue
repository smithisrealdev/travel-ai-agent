<template>
  <div class="group cursor-pointer transform transition-transform duration-300 hover:scale-105">
    <div class="relative overflow-hidden rounded-3xl shadow-lg">
      <!-- Image -->
      <div class="aspect-w-16 aspect-h-12 bg-gray-200">
        <img 
          :src="destination.image || 'https://via.placeholder.com/400x300?text=' + destination.name"
          :alt="destination.name"
          class="w-full h-64 object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent"></div>
      </div>
      
      <!-- Content Overlay -->
      <div class="absolute bottom-0 left-0 right-0 p-6 text-white">
        <div class="flex items-center gap-2 mb-2">
          <h3 class="text-2xl font-poppins font-semibold">{{ destination.name }}</h3>
        </div>
        <div class="flex items-center justify-between">
          <span class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-white/20 backdrop-blur-sm">
            📍 {{ destination.country }}
          </span>
          <span class="text-lg font-semibold">
            ฿{{ formatPrice(destination.avgCost) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Destination {
  id: number
  name: string
  country: string
  image?: string
  avgCost: number
}

defineProps<{
  destination: Destination
}>()

const formatPrice = (price: number) => {
  return new Intl.NumberFormat('en-US').format(price)
}
</script>
