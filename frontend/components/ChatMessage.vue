<template>
  <div :class="['animate-slide-in', isUser ? 'flex justify-end' : 'flex justify-start']">
    <div 
      :class="[
        'max-w-3xl rounded-2xl px-6 py-4 shadow-lg',
        isUser 
          ? 'bg-gradient-to-r from-primary to-blue-600 text-white rounded-tr-sm' 
          : 'bg-white text-gray-800 rounded-tl-sm border border-gray-100'
      ]"
    >
      <!-- Message Content -->
      <div v-if="isUser" class="whitespace-pre-wrap leading-relaxed">
        {{ message.content }}
      </div>
      <div v-else class="prose prose-sm max-w-none" v-html="renderMarkdown(message.content)" />
      
      <!-- Timestamp -->
      <p 
        :class="[
          'text-xs mt-2',
          isUser ? 'text-blue-100 text-right' : 'text-gray-400 text-left'
        ]"
      >
        {{ formatTime(message.timestamp) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import DOMPurify from 'dompurify'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

const props = defineProps<{
  message: Message
}>()

const isUser = computed(() => props.message.role === 'user')

// Configure marked
marked.setOptions({
  breaks: true,
  gfm: true,
})

const renderMarkdown = (content: string): string => {
  const rawHtml = marked.parse(content) as string
  return DOMPurify.sanitize(rawHtml)
}

const formatTime = (date: Date) => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  }).format(date)
}
</script>

<style scoped>
@keyframes slide-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-in {
  animation: slide-in 0.3s ease-out;
}

/* Prose Styling for Markdown */
:deep(.prose) {
  color: #374151;
  max-width: 100%;
}

:deep(.prose h1),
:deep(.prose h2),
:deep(.prose h3) {
  font-weight: 700;
  margin-top: 1rem;
  margin-bottom: 0.5rem;
  color: #1f2937;
}

:deep(.prose h1) { font-size: 1.875rem; }
:deep(.prose h2) { font-size: 1.5rem; }
:deep(.prose h3) { font-size: 1.25rem; }

:deep(.prose p) {
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
  line-height: 1.75;
}

:deep(.prose ul),
:deep(.prose ol) {
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

:deep(.prose a) {
  color: #3b82f6;
  text-decoration: underline;
}
</style>
