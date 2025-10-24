<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-8">
    <div class="max-w-4xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-800">🌍 Travel AI Agent</h1>
        <p class="text-gray-600 mt-2">Plan your perfect trip with AI</p>
      </div>

      <!-- Chat Container -->
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <!-- Messages Area -->
        <div ref="messagesContainer" class="h-[600px] overflow-y-auto p-6 space-y-4">
          <!-- Welcome Message -->
          <div v-if="messages.length === 0" class="flex justify-center items-center h-full">
            <div class="text-center text-gray-500">
              <div class="text-6xl mb-4">✈️</div>
              <p class="text-lg font-semibold mb-2">Welcome to Travel AI Agent!</p>
              <p class="text-sm">Ask me anything about your travel plans.</p>
              <p class="text-xs mt-4 text-gray-400">
                Example: "I want to go to Tokyo for 5 days"
              </p>
            </div>
          </div>

          <!-- Messages -->
          <template v-for="(msg, index) in messages" :key="index">
            <!-- User Message -->
            <div v-if="msg.role === 'user'" class="flex justify-end">
              <div class="bg-blue-500 text-white rounded-lg px-4 py-2 max-w-md shadow">
                <p class="whitespace-pre-wrap">{{ msg.content }}</p>
                <p class="text-xs text-blue-100 mt-1">{{ formatTime(msg.timestamp) }}</p>
              </div>
            </div>

            <!-- AI Response -->
            <div v-else class="flex justify-start">
              <div class="bg-gray-100 text-gray-800 rounded-lg px-4 py-2 max-w-md shadow">
                <pre class="whitespace-pre-wrap font-sans text-sm">{{ msg.content }}</pre>
                <p class="text-xs text-gray-500 mt-1">{{ formatTime(msg.timestamp) }}</p>
              </div>
            </div>
          </template>

          <!-- Loading Indicator -->
          <div v-if="loading" class="flex justify-start">
            <div class="bg-gray-100 rounded-lg px-4 py-2 shadow">
              <div class="flex space-x-2">
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Input Area -->
        <div class="border-t p-4 bg-gray-50">
          <form @submit.prevent="sendMessage" class="flex gap-2">
            <input
              v-model="message"
              type="text"
              placeholder="Where do you want to go? (e.g., I want to go to Tokyo for 5 days)"
              class="flex-1 px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :disabled="loading"
            />
            <button
              type="submit"
              :disabled="loading || !message.trim()"
              class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
            >
              Send
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

const config = useRuntimeConfig()
const message = ref('')
const messages = ref<Message[]>([])
const loading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

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

    // Format the AI response
    let aiContent = ''
    if (response.data.destination) {
      aiContent += `🌍 Destination: ${response.data.destination}\n`
    }
    if (response.data.duration_days) {
      aiContent += `📅 Duration: ${response.data.duration_days} days\n`
    }
    if (response.data.budget) {
      aiContent += `💰 Budget: $${response.data.budget.toFixed(2)}\n`
    }
    if (response.data.weather) {
      aiContent += `\n🌤️ Weather:\n`
      aiContent += `   Average Temperature: ${response.data.weather.avg_temp}°C\n`
      aiContent += `   Condition: ${response.data.weather.condition}\n`
    }
    if (response.data.flight_price) {
      aiContent += `\n✈️ Estimated Flight Price: $${response.data.flight_price.toFixed(2)}\n`
    }
    if (response.data.hotel_price) {
      aiContent += `🏨 Estimated Hotel Price: $${response.data.hotel_price.toFixed(2)}\n`
    }
    if (response.data.itinerary && response.data.itinerary.length > 0) {
      aiContent += `\n📋 Itinerary:\n`
      response.data.itinerary.forEach((day: any) => {
        aiContent += `   Day ${day.day}: ${day.activity}\n`
      })
    }

    // If no structured data, fallback to raw response
    if (!aiContent) {
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

// Auto-scroll to bottom
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// Format timestamp
const formatTime = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  }).format(date)
}

watch(messages, scrollToBottom, { deep: true })

// Set page metadata
useHead({
  title: 'Travel AI Agent - Plan Your Perfect Trip',
  meta: [
    {
      name: 'description',
      content: 'AI-powered travel planning assistant with personalized recommendations, weather forecasts, and budget planning.',
    },
  ],
})
</script>
