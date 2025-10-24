<template>
  <div>
    <NuxtLayout name="default">
      <!-- Hero Section -->
      <HeroSection 
        title="Plan Your Next Adventure with AI"
        subtitle="Tell us your dream trip, we'll handle the rest."
      >
        <SearchBar @search="handleSearch" />
        <NuxtLink to="/chat" class="btn-primary mt-4">
          Start Planning
        </NuxtLink>
      </HeroSection>

      <!-- Top Destinations -->
      <section class="py-16 bg-white">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 class="heading-2 text-center mb-12">Top Destinations</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <DestinationCard 
              v-for="destination in destinations"
              :key="destination.id"
              :destination="destination"
            />
          </div>
        </div>
      </section>

      <!-- How It Works -->
      <section class="py-16 gradient-sky">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 class="heading-2 text-center mb-12">How It Works</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <!-- Step 1 -->
            <div class="text-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
                <span class="text-4xl">✍️</span>
              </div>
              <h3 class="text-xl font-poppins font-semibold text-gray-900 mb-3">
                Describe Your Trip
              </h3>
              <p class="text-gray-600">
                Tell our AI about your dream destination, dates, and budget
              </p>
            </div>

            <!-- Step 2 -->
            <div class="text-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
                <span class="text-4xl">🤖</span>
              </div>
              <h3 class="text-xl font-poppins font-semibold text-gray-900 mb-3">
                AI Builds a Plan
              </h3>
              <p class="text-gray-600">
                Get personalized itineraries, hotel recommendations, and activities
              </p>
            </div>

            <!-- Step 3 -->
            <div class="text-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
                <span class="text-4xl">✈️</span>
              </div>
              <h3 class="text-xl font-poppins font-semibold text-gray-900 mb-3">
                Book & Go!
              </h3>
              <p class="text-gray-600">
                Review your plan and book everything with just a few clicks
              </p>
            </div>
          </div>

          <!-- CTA -->
          <div class="text-center mt-12">
            <NuxtLink to="/chat" class="btn-primary">
              Get Started Now
            </NuxtLink>
          </div>
        </div>
      </section>
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
const destinations = [
  {
    id: 1,
    name: 'Tokyo',
    country: 'Japan',
    image: 'https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=400',
    avgCost: 80000
  },
  {
    id: 2,
    name: 'Paris',
    country: 'France',
    image: 'https://images.unsplash.com/photo-1502602898657-3e91760cbb34?w=400',
    avgCost: 120000
  },
  {
    id: 3,
    name: 'Chiang Mai',
    country: 'Thailand',
    image: 'https://images.unsplash.com/photo-1598965675045-4c4e8f9fbb0b?w=400',
    avgCost: 50000
  },
  {
    id: 4,
    name: 'New York',
    country: 'USA',
    image: 'https://images.unsplash.com/photo-1496442226666-8d4d0e62e6e9?w=400',
    avgCost: 150000
  }
]

const handleSearch = (searchData: any) => {
  console.log('Search:', searchData)
  // Navigate to chat with search data
  navigateTo({
    path: '/chat',
    query: {
      destination: searchData.destination,
      dates: searchData.dates,
      budget: searchData.budget
    }
  })
}

// Set page metadata
useHead({
  title: 'Atravel - Plan Your Next Adventure with AI',
  meta: [
    {
      name: 'description',
      content: 'AI-powered travel planning assistant. Get personalized recommendations, itineraries, and book your perfect trip.',
    },
  ],
})
</script>
