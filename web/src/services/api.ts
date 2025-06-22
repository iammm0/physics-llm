import axios from 'axios';

export interface ChatResponse {
    response: string;
}

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
    headers: { 'Content-Type': 'application/json' },
    timeout: 30_000,
});

/**
 * 向后端 /v1/chat 发送一次聊天请求
 */
export async function chat(query: string): Promise<ChatResponse> {
    const res = await api.post<ChatResponse>('/v1/chat', { query });
    return res.data;
}

