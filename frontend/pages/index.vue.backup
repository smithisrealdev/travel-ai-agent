<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 py-8">
    <div class="max-w-5xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-5xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
          🌍 Travel AI Agent
        </h1>
        <p class="text-gray-600 mt-3 text-lg">Plan your perfect trip with AI assistance</p>
      </div>

      <!-- Chat Container -->
      <div class="bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
        <!-- Messages Area -->
        <div ref="messagesContainer" class="h-[650px] overflow-y-auto p-8 space-y-6 bg-gradient-to-b from-gray-50 to-white">
          <!-- Welcome Message -->
          <div v-if="messages.length === 0" class="flex justify-center items-center h-full">
            <div class="text-center text-gray-500 max-w-2xl">
              <div class="text-7xl mb-6 animate-bounce">✈️</div>
              <h2 class="text-2xl font-bold mb-3 text-gray-700">Welcome to Travel AI Agent!</h2>
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
          <template v-for="(msg, index) in messages" :key="index">
            <!-- User Message -->
            <div v-if="msg.role === 'user'" class="flex justify-end animate-slide-in-right">
              <div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-2xl rounded-tr-sm px-6 py-4 max-w-2xl shadow-lg">
                <p class="whitespace-pre-wrap leading-relaxed">{{ msg.content }}</p>
                <p class="text-xs text-blue-100 mt-2 text-right">{{ formatTime(msg.timestamp) }}</p>
              </div>
            </div>

            <!-- AI Response -->
            <div v-else class="flex justify-start animate-slide-in-left">
              <div class="bg-white text-gray-800 rounded-2xl rounded-tl-sm px-6 py-4 max-w-3xl shadow-lg border border-gray-100">
                <div class="prose prose-sm max-w-none" v-html="renderMarkdown(msg.content)"></div>
                <p class="text-xs text-gray-400 mt-3 text-left">{{ formatTime(msg.timestamp) }}</p>
              </div>
            </div>
          </template>

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

        <!-- Input Area -->
        <div class="border-t border-gray-200 p-6 bg-gray-50">
          <form @submit.prevent="sendMessage" class="flex gap-3">
            <input
              v-model="message"
              type="text"
              placeholder="Where do you want to go? (e.g., I want to go to Tokyo for 5 days)"
              class="flex-1 px-6 py-4 border border-gray-300 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent shadow-sm"
              :disabled="loading"
            />
            <button
              type="submit"
              :disabled="loading || !message.trim()"
              class="px-8 py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-2xl hover:from-blue-600 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed transition font-semibold shadow-lg hover:shadow-xl transform hover:scale-105"
            >
              <span v-if="!loading">Send</span>
              <span v-else>...</span>
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

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

// Configure marked for better rendering
marked.setOptions({
  breaks: true,
  gfm: true,
})

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
      // Backend now returns formatted markdown
      aiContent = response.data.response
    } else if (response.data.destination) {
      // Fallback for old format
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

// Render markdown to HTML safely
const renderMarkdown = (content: string): string => {
  const rawHtml = marked.parse(content) as string
  return DOMPurify.sanitize(rawHtml)
}

// Legacy response formatter (fallback)
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

// Auto-scroll to bottom
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

<style scoped>
/* Animations */
@keyframes slide-in-right {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes slide-in-left {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate-slide-in-right {
  animation: slide-in-right 0.3s ease-out;
}

.animate-slide-in-left {
  animation: slide-in-left 0.3s ease-out;
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}

/* Prose Styling for Markdown */
:deep(.prose) {
  color: #374151;
  max-width: 100%;
}

:deep(.prose h1) {
  font-size: 1.875rem;
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 1rem;
  color: #1f2937;
  line-height: 1.2;
}

:deep(.prose h2) {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  color: #374151;
  line-height: 1.3;
}

:deep(.prose h3) {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 1.25rem;
  margin-bottom: 0.5rem;
  color: #4b5563;
  line-height: 1.4;
}

:deep(.prose p) {
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
  line-height: 1.75;
}

:deep(.prose ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
}

:deep(.prose ol) {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
}

:deep(.prose li) {
  margin-top: 0.375rem;
  margin-bottom: 0.375rem;
  line-height: 1.75;
}

:deep(.prose strong) {
  font-weight: 600;
  color: #1f2937;
}

:deep(.prose em) {
  font-style: italic;
}

:deep(.prose code) {
  background-color: #f3f4f6;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  font-family: monospace;
}

:deep(.prose pre) {
  background-color: #1f2937;
  color: #f9fafb;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin-top: 1rem;
  margin-bottom: 1rem;
}

:deep(.prose pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}

:deep(.prose blockquote) {
  border-left: 4px solid #e5e7eb;
  padding-left: 1rem;
  font-style: italic;
  color: #6b7280;
  margin-top: 1rem;
  margin-bottom: 1rem;
}

:deep(.prose hr) {
  border: none;
  border-top: 1px solid #e5e7eb;
  margin-top: 2rem;
  margin-bottom: 2rem;
}

:deep(.prose a) {
  color: #3b82f6;
  text-decoration: underline;
}

:deep(.prose a:hover) {
  color: #2563eb;
}

:deep(.prose table) {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
  margin-bottom: 1rem;
}

:deep(.prose th) {
  background-color: #f3f4f6;
  padding: 0.5rem;
  text-align: left;
  font-weight: 600;
  border: 1px solid #e5e7eb;
}

:deep(.prose td) {
  padding: 0.5rem;
  border: 1px solid #e5e7eb;
}
</style>
