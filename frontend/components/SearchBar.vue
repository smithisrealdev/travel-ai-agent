<template>
  <div class="bg-white rounded-3xl shadow-xl p-6 max-w-5xl w-full">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <!-- Destination Input -->
      <div class="md:col-span-1">
        <label class="block text-sm font-medium text-gray-700 mb-2">Where</label>
        <input
          v-model="localDestination"
          type="text"
          placeholder="Destination"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-primary focus:border-transparent transition"
        />
      </div>

      <!-- Date Picker -->
      <div class="md:col-span-1">
        <label class="block text-sm font-medium text-gray-700 mb-2">When</label>
        <input
          v-model="localDates"
          type="text"
          placeholder="Select dates"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-primary focus:border-transparent transition"
        />
      </div>

      <!-- Budget Input -->
      <div class="md:col-span-1">
        <label class="block text-sm font-medium text-gray-700 mb-2">Budget</label>
        <input
          v-model="localBudget"
          type="text"
          placeholder="Budget (THB)"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-primary focus:border-transparent transition"
        />
      </div>

      <!-- Search Button -->
      <div class="md:col-span-1 flex items-end">
        <button
          @click="handleSearch"
          class="w-full btn-primary"
        >
          Search
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const localDestination = ref('')
const localDates = ref('')
const localBudget = ref('')

const emit = defineEmits<{
  search: [{ destination: string; dates: string; budget: string }]
}>()

const handleSearch = () => {
  emit('search', {
    destination: localDestination.value,
    dates: localDates.value,
    budget: localBudget.value
  })
}
</script>
