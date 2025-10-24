<template>
  <section class="relative overflow-hidden gradient-sky py-20 md:py-32">
    <div class="absolute inset-0 opacity-10">
      <div 
        v-if="backgroundImage"
        :style="{ backgroundImage: `url(${backgroundImage})` }"
        class="w-full h-full bg-cover bg-center"
      />
    </div>
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
      <div class="text-center">
        <h1 class="heading-1 mb-6">
          {{ title }}
        </h1>
        <p class="text-xl md:text-2xl text-gray-600 mb-8 max-w-3xl mx-auto">
          {{ subtitle }}
        </p>
        <div class="flex flex-col items-center gap-4">
          <slot />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subtitle: string
  backgroundImage?: string
}>()
</script>
