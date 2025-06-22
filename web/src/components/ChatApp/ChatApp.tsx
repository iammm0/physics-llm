import React, { useState } from 'react';
import Box from '@mui/material/Box';
import ChatDrawer from '../ChatDrawer/ChatDrawer';
import ChatWindow from '../ChatWindow/ChatWindow';
import styles from './ChatApp.module.css';

const ChatApp: React.FC = () => {
    const [activeChatId, setActiveChatId] = useState<string>('default');

    return (
        <Box className={styles.appContainer}>
            {/* 左侧会话列表抽屉 */}
            <ChatDrawer
                activeChatId={activeChatId}
                onSelectChat={setActiveChatId}
            />

            {/* 右侧聊天主区 */}
            <Box component="main" className={styles.chatContainer}>
                <ChatWindow sessionId={activeChatId} />
            </Box>
        </Box>
    );
};

export default ChatApp;
