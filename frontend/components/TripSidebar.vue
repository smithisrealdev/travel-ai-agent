<template>
  <aside class="w-80 bg-white border-r border-gray-200 flex flex-col h-screen">
    <!-- Header -->
    <div class="p-6 border-b border-gray-200">
      <button
        @click="$emit('new-chat')"
        class="w-full btn-primary"
      >
        + New Chat
      </button>
    </div>

    <!-- My Trips -->
    <div class="flex-1 overflow-y-auto">
      <div class="p-6">
        <h2 class="text-lg font-poppins font-semibold text-gray-900 mb-4">My Trips</h2>
        
        <div v-if="trips && trips.length > 0" class="space-y-3">
          <div
            v-for="trip in trips"
            :key="trip.id"
            @click="$emit('select-trip', trip.id)"
            :class="[
              'p-4 rounded-2xl cursor-pointer transition-all',
              activeId === trip.id 
                ? 'bg-primary text-white shadow-lg' 
                : 'bg-gray-50 hover:bg-gray-100'
            ]"
          >
            <!-- Trip Thumbnail -->
            <div class="flex items-start gap-3">
              <div class="w-16 h-16 rounded-xl overflow-hidden flex-shrink-0 bg-gray-200">
                <img 
                  v-if="trip.thumbnail"
                  :src="trip.thumbnail"
                  :alt="trip.destination"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-2xl">
                  🌍
                </div>
              </div>
              
              <div class="flex-1 min-w-0">
                <h3 
                  :class="[
                    'font-poppins font-semibold truncate',
                    activeId === trip.id ? 'text-white' : 'text-gray-900'
                  ]"
                >
                  {{ trip.destination }}
                </h3>
                <p 
                  :class="[
                    'text-sm truncate',
                    activeId === trip.id ? 'text-blue-100' : 'text-gray-600'
                  ]"
                >
                  {{ trip.dates }}
                </p>
                <span 
                  :class="[
                    'inline-block mt-1 px-2 py-1 rounded-full text-xs font-semibold',
                    getStatusClass(trip.status, activeId === trip.id)
                  ]"
                >
                  {{ trip.status }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-8 text-gray-500">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-sm">No trips yet</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
interface Trip {
  id: number
  destination: string
  dates: string
  status: string
  thumbnail?: string
}

defineProps<{
  trips?: Trip[]
  activeId?: number
}>()

defineEmits<{
  'new-chat': []
  'select-trip': [id: number]
}>()

const getStatusClass = (status: string, isActive: boolean) => {
  if (isActive) {
    return 'bg-white/20 text-white'
  }
  
  const statusMap: Record<string, string> = {
    upcoming: 'bg-green-100 text-green-700',
    completed: 'bg-gray-100 text-gray-700',
    cancelled: 'bg-red-100 text-red-700',
    planning: 'bg-blue-100 text-blue-700',
  }
  
  return statusMap[status] || 'bg-gray-100 text-gray-700'
}
</script>
