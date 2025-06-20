<!-- src/components/ChatWindow.vue -->
<template>
  <div class="chat-window">
    <div class="messages">
      <div
          v-for="(msg, i) in messages"
          :key="i"
          :class="['message', msg.role]"
      >
        <strong>{{ msg.role === 'user' ? '👤' : '🤖' }}</strong>
        <span>{{ msg.content }}</span>
      </div>
    </div>

    <form @submit.prevent="onSend" class="input-area">
      <textarea
          v-model="input"
          placeholder="输入你的问题…"
          :disabled="loading"
          rows="2"
      ></textarea>
      <button type="submit" :disabled="!input.trim() || loading">
        {{ loading ? '…' : '发送' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const messages = ref([
  { role: 'system', content: '欢迎使用 Physics-LLM，您可以问我任何物理问题。' }
])
const input = ref('')
const loading = ref(false)

async function onSend() {
  const question = input.value.trim()
  if (!question) return

  // 先把用户消息推到界面
  messages.value.push({ role: 'user', content: question })
  input.value = ''
  loading.value = true

  try {
    // 调用后端
    const res = await axios.post('/api/v1/chat', { query: question })
    const answer = res.data.response

    // 推入模型回复
    messages.value.push({ role: 'assistant', content: answer })
  } catch (err) {
    console.error(err)
    messages.value.push({
      role: 'assistant',
      content: '❗️ 发生错误，请稍后重试。'
    })
  } finally {
    loading.value = false
    // 自动滚动到底部
    scrollToBottom()
  }
}

function scrollToBottom() {
  // 下次 DOM 更新后执行
  setTimeout(() => {
    const container = document.querySelector('.messages')
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  })
}
</script>

<style scoped>
.chat-window {
  border: 1px solid #ddd;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  height: 80vh;
}
.messages {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  background: #f9f9f9;
}
.message {
  margin-bottom: 0.75rem;
  display: flex;
  align-items: flex-start;
}
.message.system {
  color: #666;
  font-style: italic;
}
.message.user span {
  background: #e0f7fa;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  margin-left: 0.5rem;
}
.message.assistant span {
  background: #fff9c4;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  margin-left: 0.5rem;
}
.input-area {
  display: flex;
  gap: 0.5rem;
  padding: 0.5rem;
  border-top: 1px solid #ddd;
}
textarea {
  flex: 1;
  resize: none;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}
button {
  padding: 0 1rem;
  border: none;
  background: #1976d2;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}
button:disabled {
  background: #90caf9;
  cursor: not-allowed;
}
</style>
