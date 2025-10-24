<template>
  <div class="bg-white rounded-2xl shadow-lg p-6">
    <h3 class="text-lg font-poppins font-semibold text-gray-900 mb-4">Quick Actions</h3>
    
    <div class="space-y-2">
      <button
        v-for="action in actions"
        :key="action.action"
        @click="$emit('action', action.action)"
        class="w-full flex items-center gap-3 p-3 rounded-xl hover:bg-gray-50 transition text-left"
      >
        <span class="text-2xl">{{ action.icon }}</span>
        <span class="text-sm font-medium text-gray-700">{{ action.label }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface QuickAction {
  icon: string
  label: string
  action: string
}

defineProps<{
  actions: QuickAction[]
}>()

defineEmits<{
  action: [action: string]
}>()
</script>
