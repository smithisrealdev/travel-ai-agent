<template>
  <div class="flex flex-col h-screen max-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-gradient-to-r from-blue-600 to-blue-700 text-white p-4 shadow-lg">
      <div class="flex items-center justify-between max-w-4xl mx-auto">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-white rounded-full flex items-center justify-center">
            <span class="text-2xl">✈️</span>
          </div>
          <div>
            <h1 class="text-xl font-bold">Atravel AI</h1>
            <p class="text-sm text-blue-100">Your Travel Companion</p>
          </div>
        </div>
        <button
          @click="clearChat"
          class="px-4 py-2 bg-blue-800 hover:bg-blue-900 rounded-lg text-sm transition-colors"
        >
          New Chat
        </button>
      </div>
    </div>

    <!-- Chat Messages -->
    <div
      ref="chatContainer"
      class="flex-1 overflow-y-auto p-4 space-y-4 max-w-4xl mx-auto w-full"
    >
      <div v-if="messages.length === 0" class="text-center py-12">
        <span class="text-6xl mb-4 block">🌍</span>
        <h2 class="text-2xl font-semibold text-gray-700 mb-2">Welcome to Atravel!</h2>
        <p class="text-gray-500">Ask me anything about your travel plans, weather, or packing tips.</p>
      </div>

      <div
        v-for="(message, index) in messages"
        :key="index"
        :class="[
          'flex',
          message.role === 'user' ? 'justify-end' : 'justify-start'
        ]"
      >
        <div
          :class="[
            'max-w-[70%] rounded-2xl px-4 py-3 shadow-sm',
            message.role === 'user'
              ? 'bg-blue-600 text-white rounded-br-none'
              : 'bg-white text-gray-800 rounded-bl-none border border-gray-200'
          ]"
        >
          <div class="flex items-start space-x-2">
            <div
              v-if="message.role === 'assistant'"
              class="flex-shrink-0 w-6 h-6 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white text-xs font-bold mt-1"
            >
              A
            </div>
            <div class="flex-1">
              <p class="whitespace-pre-wrap break-words">{{ message.content }}</p>
              <span
                :class="[
                  'text-xs mt-1 block',
                  message.role === 'user' ? 'text-blue-100' : 'text-gray-400'
                ]"
              >
                {{ formatTime(message.timestamp) }}
              </span>
            </div>
            <div
              v-if="message.role === 'user'"
              class="flex-shrink-0 w-6 h-6 rounded-full bg-white text-blue-600 flex items-center justify-center text-xs font-bold mt-1"
            >
              U
            </div>
          </div>
        </div>
      </div>

      <!-- Loading Indicator -->
      <div v-if="isLoading" class="flex justify-start">
        <div class="max-w-[70%] rounded-2xl rounded-bl-none px-4 py-3 bg-white border border-gray-200 shadow-sm">
          <div class="flex items-center space-x-2">
            <div class="w-6 h-6 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white text-xs font-bold">
              A
            </div>
            <div class="flex space-x-1">
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0ms"></div>
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 150ms"></div>
              <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 300ms"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        <p class="font-semibold">Error:</p>
        <p>{{ error }}</p>
      </div>
    </div>

    <!-- Input Area -->
    <div class="border-t border-gray-200 bg-white p-4">
      <div class="max-w-4xl mx-auto">
        <form @submit.prevent="sendMessage" class="flex space-x-2">
          <input
            v-model="userInput"
            type="text"
            :disabled="isLoading"
            placeholder="Ask about your travel plans, weather, packing..."
            class="flex-1 px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
          />
          <button
            type="submit"
            :disabled="isLoading || !userInput.trim()"
            class="px-6 py-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-semibold"
          >
            <span v-if="!isLoading">Send</span>
            <span v-else>...</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { collection, addDoc, query, where, orderBy, getDocs, updateDoc, doc } from 'firebase/firestore'
import { db } from '~/firebase/config'

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: Date
}

const ATRAVEL_SYSTEM_PROMPT = `You are Atravel, a friendly AI travel companion.\nYou help travelers with trip planning, adjusting itineraries, weather updates, and packing advice.\nYour tone is warm, human-like, and concise.\nIf a user mentions 'weather', 'packing', or 'change plan', respond proactively with real data.`

const messages = ref<Message[]>([])
const userInput = ref('')
const isLoading = ref(false)
const error = ref('')
const chatContainer = ref<HTMLElement | null>(null)
const sessionId = ref('')
const userId = ref('user-' + Date.now())

onMounted(async () => {
  sessionId.value = 'session-' + Date.now()
  await loadChatHistory()
  
  if (messages.value.length === 0) {
    messages.value.push({
      role: 'assistant',
      content: 'Hello! I\'m Atravel, your AI travel companion. How can I help you with your travel plans today?',
      timestamp: new Date()
    })
  }
})

async function loadChatHistory() {
  try {
    const q = query(
      collection(db, 'chat_sessions'),
      where('userId', '==', userId.value),
      orderBy('createdAt', 'desc')
    )
    
    const querySnapshot = await getDocs(q)
    
    if (!querySnapshot.empty) {
      const latestSession = querySnapshot.docs[0]
      sessionId.value = latestSession.id
      const data = latestSession.data()
      
      if (data.messages && Array.isArray(data.messages)) {
        messages.value = data.messages.map((msg: any) => ({
          ...msg,
          timestamp: msg.timestamp?.toDate() || new Date()
        }))
      }
    }
  } catch (err) {
    console.error('Error loading chat history:', err)
  }
}

async function sendMessage() {
  if (!userInput.value.trim() || isLoading.value) return

  const userMessage: Message = {
    role: 'user',
    content: userInput.value.trim(),
    timestamp: new Date()
  }

  messages.value.push(userMessage)
  userInput.value = ''
  error.value = ''
  isLoading.value = true

  await scrollToBottom()

  try {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        messages: [
          { role: 'system', content: ATRAVEL_SYSTEM_PROMPT },
          ...messages.value.map(msg => ({
            role: msg.role,
            content: msg.content
          }))
        ]
      })
    })

    if (!response.ok) {
      throw new Error('Failed to get response from AI')
    }

    const data = await response.json()
    
    const assistantMessage: Message = {
      role: 'assistant',
      content: data.message || data.content || 'Sorry, I couldn\'t generate a response.',
      timestamp: new Date()
    }

    messages.value.push(assistantMessage)
    await saveChatHistory()
    await scrollToBottom()
  } catch (err: any) {
    error.value = err.message || 'An error occurred. Please try again.'
    console.error('Chat error:', err)
  } finally {
    isLoading.value = false
  }
}

async function saveChatHistory() {
  try {
    const chatData = {
      userId: userId.value,
      sessionId: sessionId.value,
      messages: messages.value.map(msg => ({
        role: msg.role,
        content: msg.content,
        timestamp: msg.timestamp
      })),
      updatedAt: new Date()
    }

    if (sessionId.value.startsWith('session-')) {
      const docRef = await addDoc(collection(db, 'chat_sessions'), {
        ...chatData,
        createdAt: new Date()
      })
      sessionId.value = docRef.id
    } else {
      const docRef = doc(db, 'chat_sessions', sessionId.value)
      await updateDoc(docRef, chatData)
    }
  } catch (err) {
    console.error('Error saving chat history:', err)
  }
}

async function clearChat() {
  messages.value = []
  sessionId.value = 'session-' + Date.now()
  error.value = ''
  
  messages.value.push({
    role: 'assistant',
    content: 'Hello! I\'m Atravel, your AI travel companion. How can I help you with your travel plans today?',
    timestamp: new Date()
  })
}

async function scrollToBottom() {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

function formatTime(date: Date): string {
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>