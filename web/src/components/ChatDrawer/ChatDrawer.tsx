// ChatDrawer.tsx
import React, { useState } from 'react';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import Button from '@mui/material/Button';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Divider from '@mui/material/Divider';
import MenuIcon from '@mui/icons-material/Menu';
import CloseIcon from '@mui/icons-material/Close';
import ChatBubbleOutlineIcon from '@mui/icons-material/ChatBubbleOutline';
import Typography from '@mui/material/Typography';
import styles from './ChatDrawer.module.css';

interface ChatSession {
    id: string;
    name: string;
}

interface ChatDrawerProps {
    activeChatId: string;
    onSelectChat: (id: string) => void;
}

const ChatDrawer: React.FC<ChatDrawerProps> = ({ activeChatId, onSelectChat }) => {
    const [open, setOpen] = useState(false);
    const [sessions, setSessions] = useState<ChatSession[]>([
        { id: 'default', name: '默认聊天' },
    ]);

    const toggleDrawer = (flag: boolean) => () => setOpen(flag);

    const createNewSession = () => {
        const newId = crypto.randomUUID();
        const newName = `聊天 ${sessions.length + 1}`;
        setSessions(prev => [...prev, { id: newId, name: newName }]);
        onSelectChat(newId);
        setOpen(false);
    };

    return (
        <>
            {!open && (
                <IconButton
                    size="small"
                    color="primary"
                    onClick={toggleDrawer(true)}
                    className={styles.triggerButton}
                >
                    <MenuIcon fontSize="small" />
                </IconButton>
            )}

            <Drawer
                anchor="left"
                open={open}
                onClose={toggleDrawer(false)}
                classes={{ paper: styles.drawerPaper }}
            >
                {/* header */}
                <div className={styles.header}>
                    <Typography className={styles.title}>Physics-LLM</Typography>
                    <IconButton
                        size="small"
                        onClick={toggleDrawer(false)}
                        className={styles.closeButton}
                    >
                        <CloseIcon fontSize="small" />
                    </IconButton>
                </div>

                {/* session list */}
                <List className={styles.sessionList}>
                    {sessions.map(s => (
                        <ListItemButton
                            key={s.id}
                            selected={s.id === activeChatId}
                            onClick={() => {
                                onSelectChat(s.id);
                                setOpen(false);
                            }}
                            className={styles.sessionItem}
                        >
                            <ListItemIcon>
                                <ChatBubbleOutlineIcon fontSize="small" />
                            </ListItemIcon>
                            <ListItemText primary={s.name} />
                        </ListItemButton>
                    ))}
                </List>
                <Divider />

                {/* new session */}
                <div className={styles.newSessionContainer}>
                    <Button
                        variant="outlined"
                        onClick={createNewSession}
                        className={styles.newSessionButton}
                    >
                        + 新建聊天
                    </Button>
                </div>
            </Drawer>
        </>
    );
};

export default ChatDrawer;
