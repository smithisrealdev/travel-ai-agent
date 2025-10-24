<template>
  <div class="flex h-screen overflow-hidden bg-background">
    <!-- Left Sidebar -->
    <TripSidebar 
      :trips="myTrips"
      :active-id="activeTrip?.id"
      @new-chat="handleNewChat"
      @select-trip="handleSelectTrip"
    />

    <!-- Main Chat Panel -->
    <div class="flex-1 flex flex-col">
      <!-- Chat Header -->
      <div class="bg-white border-b border-gray-200 px-6 py-4">
        <div class="flex items-center gap-3">
          <span class="text-3xl">🤖</span>
          <div>
            <h1 class="text-xl font-poppins font-semibold text-gray-900">Travel AI Assistant</h1>
            <p class="text-sm text-gray-500">Your personal travel planning companion</p>
          </div>
        </div>
      </div>

      <!-- Messages Area -->
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-6 space-y-4">
        <!-- Welcome Message -->
        <div v-if="messages.length === 0" class="flex justify-center items-center h-full">
          <div class="text-center text-gray-500 max-w-2xl">
            <div class="text-7xl mb-6 animate-bounce">✈️</div>
            <h2 class="text-2xl font-poppins font-semibold mb-3 text-gray-700">Welcome to Travel AI Agent!</h2>
            <p class="text-lg mb-6">Ask me anything about your travel plans in Thai or English.</p>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-8">
              <button 
                @click="sendExampleMessage('ฉันอยากไปเที่ยวเชียงใหม่ 7 วัน งบ 50,000 บาท')"
                class="p-4 bg-blue-50 hover:bg-blue-100 rounded-xl text-left transition border border-blue-200"
              >
                <div class="text-sm text-blue-600 font-semibold mb-1">🇹🇭 Thai Example</div>
                <div class="text-xs text-gray-600">"ฉันอยากไปเที่ยวเชียงใหม่ 7 วัน งบ 50,000 บาท"</div>
              </button>
              
              <button 
                @click="sendExampleMessage('Plan a 5-day trip to Tokyo with 80,000 THB budget')"
                class="p-4 bg-purple-50 hover:bg-purple-100 rounded-xl text-left transition border border-purple-200"
              >
                <div class="text-sm text-purple-600 font-semibold mb-1">🇬🇧 English Example</div>
                <div class="text-xs text-gray-600">"Plan a 5-day trip to Tokyo with 80,000 THB budget"</div>
              </button>
            </div>
          </div>
        </div>

        <!-- Messages -->
        <ChatMessage 
          v-for="(msg, index) in messages"
          :key="index"
          :message="msg"
        />

        <!-- Loading Indicator -->
        <div v-if="loading" class="flex justify-start animate-fade-in">
          <div class="bg-white rounded-2xl rounded-tl-sm px-6 py-4 shadow-lg border border-gray-100">
            <div class="flex items-center space-x-3">
              <div class="flex space-x-2">
                <div class="w-2 h-2 bg-blue-400 rounded-full animate-bounce"></div>
                <div class="w-2 h-2 bg-blue-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                <div class="w-2 h-2 bg-blue-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
              </div>
              <span class="text-sm text-gray-500">AI is thinking...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Smart Suggestions -->
      <div v-if="smartSuggestions.length > 0" class="px-6 py-3 border-t border-gray-200 bg-gray-50">
        <div class="flex gap-2 overflow-x-auto">
          <button
            v-for="suggestion in smartSuggestions"
            :key="suggestion"
            @click="sendExampleMessage(suggestion)"
            class="px-4 py-2 bg-white border border-gray-300 rounded-full text-sm font-medium text-gray-700 hover:bg-gray-100 transition whitespace-nowrap"
          >
            {{ suggestion }}
          </button>
        </div>
      </div>

      <!-- Input Area -->
      <div class="border-t border-gray-200 p-6 bg-white">
        <form @submit.prevent="sendMessage" class="flex gap-3">
          <input
            v-model="message"
            type="text"
            placeholder="Where do you want to go? (e.g., I want to go to Tokyo for 5 days)"
            class="flex-1 px-6 py-4 border border-gray-300 rounded-2xl focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent shadow-sm"
            :disabled="loading"
          />
          <button
            type="submit"
            :disabled="loading || !message.trim()"
            class="px-8 py-4 gradient-primary text-white rounded-2xl hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition font-poppins font-semibold"
          >
            <span v-if="!loading">Send</span>
            <span v-else>...</span>
          </button>
        </form>
      </div>
    </div>

    <!-- Right Sidebar -->
    <aside class="w-80 bg-white border-l border-gray-200 p-6 overflow-y-auto">
      <!-- Active Trip Card -->
      <div v-if="activeTrip" class="card mb-6">
        <h3 class="text-lg font-poppins font-semibold text-gray-900 mb-4">Active Trip</h3>
        <div class="space-y-3">
          <div>
            <p class="text-sm text-gray-500">Destination</p>
            <p class="font-semibold text-gray-900">{{ activeTrip.destination }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Dates</p>
            <p class="font-semibold text-gray-900">{{ activeTrip.dates }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Status</p>
            <span class="inline-block px-3 py-1 rounded-full text-xs font-semibold bg-green-100 text-green-700">
              {{ activeTrip.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <QuickActions 
        :actions="quickActions"
        @action="handleQuickAction"
      />
    </aside>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

interface Trip {
  id: number
  destination: string
  dates: string
  status: string
  thumbnail?: string
}

const config = useRuntimeConfig()
const message = ref('')
const messages = ref<Message[]>([])
const loading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

const route = useRoute()

// Sample data
const myTrips = ref<Trip[]>([
  {
    id: 1,
    destination: 'Chiang Mai',
    dates: 'Nov 1-7, 2025',
    status: 'upcoming',
    thumbnail: 'https://images.unsplash.com/photo-1598965675045-4c4e8f9fbb0b?w=200'
  }
])

const activeTrip = ref<Trip | null>(null)

const smartSuggestions = ref([
  'Show me hotels under $100',
  'Find local restaurants',
  'Check weather forecast',
  'Add shopping to itinerary'
])

const quickActions = [
  { icon: '✈️', label: 'Reschedule flight', action: 'reschedule-flight' },
  { icon: '🍽️', label: 'Add restaurant', action: 'add-restaurant' },
  { icon: '🏨', label: 'Change hotel', action: 'change-hotel' },
  { icon: '🌤️', label: 'Get weather', action: 'check-weather' }
]

// Handle initial query params
onMounted(() => {
  const { destination, dates, budget } = route.query
  if (destination) {
    const queryMessage = `I want to visit ${destination}${dates ? ' on ' + dates : ''}${budget ? ' with a budget of ' + budget : ''}`
    message.value = queryMessage
    sendMessage()
  }
})

const handleNewChat = () => {
  messages.value = []
  activeTrip.value = null
}

const handleSelectTrip = (id: number) => {
  const trip = myTrips.value.find(t => t.id === id)
  if (trip) {
    activeTrip.value = trip
  }
}

const handleQuickAction = (action: string) => {
  console.log('Quick action:', action)
  // Handle quick actions
  const actionMessages: Record<string, string> = {
    'reschedule-flight': 'I want to reschedule my flight',
    'add-restaurant': 'Suggest some good restaurants',
    'change-hotel': 'Show me alternative hotels',
    'check-weather': 'What\'s the weather forecast?'
  }
  
  if (actionMessages[action]) {
    message.value = actionMessages[action]
    sendMessage()
  }
}

const sendExampleMessage = (exampleText: string) => {
  message.value = exampleText
  sendMessage()
}

const sendMessage = async () => {
  if (!message.value.trim()) return

  const userMessage = message.value
  message.value = ''

  // Add user message
  messages.value.push({
    role: 'user',
    content: userMessage,
    timestamp: new Date()
  })

  loading.value = true

  try {
    const response = await axios.post(`${config.public.apiBase}/api/plan`, {
      message: userMessage
    })

    // Get AI response
    let aiContent = ''
    if (response.data.response) {
      aiContent = response.data.response
    } else if (response.data.destination) {
      aiContent = formatLegacyResponse(response.data)
    } else {
      aiContent = JSON.stringify(response.data, null, 2)
    }

    // Add AI response
    messages.value.push({
      role: 'assistant',
      content: aiContent,
      timestamp: new Date()
    })
  } catch (error: any) {
    console.error('Error:', error)
    let errorMessage = '❌ Sorry, something went wrong. Please try again.'
    
    if (error.response?.data?.message) {
      errorMessage = `❌ Error: ${error.response.data.message}`
    } else if (error.message) {
      errorMessage = `❌ Error: ${error.message}`
    }

    messages.value.push({
      role: 'assistant',
      content: errorMessage,
      timestamp: new Date()
    })
  } finally {
    loading.value = false
  }
}

const formatLegacyResponse = (data: any): string => {
  let md = ''
  
  if (data.destination) {
    md += `🌍 **Destination:** ${data.destination}\n\n`
  }
  if (data.duration_days) {
    md += `📅 **Duration:** ${data.duration_days} days\n\n`
  }
  if (data.budget) {
    md += `💰 **Budget:** $${data.budget.toFixed(2)}\n\n`
  }
  if (data.weather) {
    md += `## 🌤️ Weather\n\n`
    md += `- **Temperature:** ${data.weather.avg_temp}°C\n`
    md += `- **Condition:** ${data.weather.condition}\n\n`
  }
  if (data.flight_price) {
    md += `✈️ **Estimated Flight Price:** $${data.flight_price.toFixed(2)}\n\n`
  }
  if (data.hotel_price) {
    md += `🏨 **Estimated Hotel Price:** $${data.hotel_price.toFixed(2)}\n\n`
  }
  if (data.itinerary && data.itinerary.length > 0) {
    md += `## 📋 Itinerary\n\n`
    data.itinerary.forEach((day: any) => {
      md += `**Day ${day.day}:** ${day.activity}\n\n`
    })
  }
  
  return md || JSON.stringify(data, null, 2)
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTo({
        top: messagesContainer.value.scrollHeight,
        behavior: 'smooth'
      })
    }
  })
}

watch(messages, scrollToBottom, { deep: true })

// Set page metadata
useHead({
  title: 'Chat - Travel AI Agent',
  meta: [
    {
      name: 'description',
      content: 'Chat with our AI travel assistant to plan your perfect trip.',
    },
  ],
})

definePageMeta({
  layout: false
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}
</style>
