import React, { useState, useRef, useEffect, type FormEvent } from 'react';
import { chat } from '../services/api';
import './ChatWindow.css';

type Role = 'user' | 'assistant' | 'system';
interface Message {
    role: Role;
    content: string;
}

export const ChatWindow: React.FC = () => {
    const [messages, setMessages] = useState<Message[]>([
        { role: 'system', content: '欢迎使用 Physics-LLM，您可以问我任何物理问题。' }
    ]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

    // 自动滚动到底部
    useEffect(() => {
        const c = containerRef.current;
        if (c) c.scrollTop = c.scrollHeight;
    }, [messages]);

    function formatContent(text: string) {
        return text.split('\n').map((line, i) => (
            <React.Fragment key={i}>
                {line}
                {i < text.split('\n').length - 1 && <br />}
            </React.Fragment>
        ));
    }

    async function handleSubmit(e: FormEvent) {
        e.preventDefault();
        const question = input.trim();
        if (!question) return;

        // 添加用户消息
        setMessages(m => [...m, { role: 'user', content: question }]);
        setInput('');
        setLoading(true);

        try {
            const res = await chat(question);
            setMessages(m => [...m, { role: 'assistant', content: res.response }]);
        } catch (err) {
            console.error(err);
            setMessages(m => [...m, { role: 'assistant', content: '❗️ 发生错误，请稍后重试。' }]);
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="chat-window">
            <div className="messages" ref={containerRef}>
                {messages.map((msg, idx) => (
                    <div key={idx} className={`message ${msg.role}`}>
                        <strong>
                            {msg.role === 'user' ? '👤' :
                                msg.role === 'system' ? '⚙️' : '🤖'}
                        </strong>
                        <div className="content">
                            {formatContent(msg.content)}
                        </div>
                    </div>
                ))}
            </div>

            <form className="input-area" onSubmit={handleSubmit}>
        <textarea
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder="输入你的问题…"
            disabled={loading}
            rows={2}
        />
                <button type="submit" disabled={!input.trim() || loading}>
                    {loading ? '…' : '发送'}
                </button>
            </form>
        </div>
    );
};
