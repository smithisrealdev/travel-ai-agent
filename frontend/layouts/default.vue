<template>
  <div class="min-h-screen flex flex-col bg-background">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <!-- Logo -->
          <NuxtLink to="/" class="flex items-center gap-2">
            <span class="text-3xl">🌍</span>
            <span class="text-xl font-poppins font-semibold text-gray-900">Atravel</span>
          </NuxtLink>

          <!-- Navigation -->
          <nav class="hidden md:flex items-center gap-6">
            <NuxtLink to="/" class="text-gray-600 hover:text-primary transition font-medium">
              Home
            </NuxtLink>
            <NuxtLink to="/chat" class="text-gray-600 hover:text-primary transition font-medium">
              Chat
            </NuxtLink>
          </nav>

          <!-- Right Side -->
          <div class="flex items-center gap-4">
            <button class="p-2 hover:bg-gray-100 rounded-full transition">
              <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
            <button class="p-2 hover:bg-gray-100 rounded-full transition relative">
              <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
              <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
            </button>
            <div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center text-white font-semibold">
              U
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <!-- Logo & Tagline -->
          <div class="col-span-1">
            <div class="flex items-center gap-2 mb-4">
              <span class="text-3xl">🌍</span>
              <span class="text-xl font-poppins font-semibold text-gray-900">Atravel</span>
            </div>
            <p class="text-gray-600 text-sm">
              Plan your perfect trip with AI assistance
            </p>
          </div>

          <!-- Links -->
          <div>
            <h3 class="font-poppins font-semibold text-gray-900 mb-4">Company</h3>
            <ul class="space-y-2">
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">About</a></li>
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">How It Works</a></li>
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">Contact</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-poppins font-semibold text-gray-900 mb-4">Support</h3>
            <ul class="space-y-2">
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">Help Center</a></li>
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">Terms of Service</a></li>
              <li><a href="#" class="text-gray-600 hover:text-primary transition text-sm">Privacy Policy</a></li>
            </ul>
          </div>

          <div>
            <h3 class="font-poppins font-semibold text-gray-900 mb-4">Follow Us</h3>
            <div class="flex gap-3">
              <a href="#" class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-primary hover:text-white transition">
                <span class="text-xl">📘</span>
              </a>
              <a href="#" class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-primary hover:text-white transition">
                <span class="text-xl">📷</span>
              </a>
              <a href="#" class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-primary hover:text-white transition">
                <span class="text-xl">🐦</span>
              </a>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-200 mt-8 pt-8 text-center text-sm text-gray-600">
          <p>&copy; 2025 Atravel. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
</script>
