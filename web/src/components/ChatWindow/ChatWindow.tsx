import React, { useState, useRef, useEffect, type FormEvent } from 'react';
import Paper from '@mui/material/Paper';
import InputBase from '@mui/material/InputBase';
import IconButton from '@mui/material/IconButton';
import CircularProgress from '@mui/material/CircularProgress';
import SendIcon from '@mui/icons-material/Send';
import Typography from '@mui/material/Typography';
import styles from './ChatWindow.module.css';
import {chat} from "../../services/api.ts";
import Divider from "@mui/material/Divider";
import MessageContent from "../MessageContent/MessageContent.tsx";

type Role = 'user' | 'assistant' | 'system';
interface Message { role: Role; content: string; }

interface ChatWindowProps {
    sessionId: string;
}

const ChatWindow: React.FC<ChatWindowProps> = ({ sessionId }) => {
    const [messages, setMessages] = useState<Message[]>([
        { role: 'system', content: '欢迎使用天津城建大学理学院物理研究社所研发的 Physics-LLM v0.0.1，您可以问我任何物理问题。' }
    ]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

    // 滚动到底部
    useEffect(() => {
        const c = containerRef.current;
        if (c) c.scrollTop = c.scrollHeight;
    }, [messages]);

    // 切换会话时重置
    useEffect(() => {
        setMessages([{ role: 'system', content: '欢迎使用天津城建大学理学院物理研究社所研发的 Physics-LLM v0.0.1，您可以问我任何物理问题。' }]);
        setInput('');
        setLoading(false);
    }, [sessionId]);

    async function handleSubmit(e: FormEvent) {
        e.preventDefault();
        const question = input.trim();
        if (!question) return;
        setMessages(m => [...m, { role: 'user', content: question }]);
        setInput('');
        setLoading(true);
        try {
            const res = await chat(question);
            setMessages(m => [...m, { role: 'assistant', content: res.response }]);
        } catch {
            setMessages(m => [...m, { role: 'assistant', content: '❗️ 发生错误，请稍后重试。' }]);
        } finally {
            setLoading(false);
        }
    }
    return (
        <div className={styles.chatWindow}>
            <div className={styles.messages} ref={containerRef}>
                {messages.map((msg, idx) => (
                    <React.Fragment key={idx}>
                        <div className={`${styles.message} ${styles[msg.role]}`}>
                            <Typography
                                variant="body2"
                                component="div"
                                className={styles.content}
                            >
                                <MessageContent content={msg.content} />
                            </Typography>
                        </div>
                        {idx < messages.length - 1 && (
                            <Divider variant="fullWidth" className={styles.divider} />
                        )}
                    </React.Fragment>
                ))}
            </div>

            {/* === ChatGPT-style 输入区 === */}
            <Paper
                component="form"
                className={styles.inputArea}
                elevation={0}
                onSubmit={handleSubmit}
            >
                <InputBase
                    className={styles.inputField}
                    placeholder="输入你的问题…（Shift+Enter 换行，Enter 发送）"
                    multiline
                    maxRows={6}
                    value={input}
                    onChange={e => setInput(e.target.value)}
                    disabled={loading}
                    onKeyDown={e => {
                        if (e.key === 'Enter' && !e.shiftKey) {
                            e.preventDefault();
                            handleSubmit(e);
                        }
                    }}
                />

                <IconButton
                    type="submit"
                    disabled={!input.trim() || loading}
                    className={styles.sendButton}
                >
                    {loading
                        ? <CircularProgress size={20} />
                        : <SendIcon className={styles.sendIcon} />
                    }
                </IconButton>
            </Paper>
        </div>
    );
};

export default ChatWindow;
