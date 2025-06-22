import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeRaw from 'rehype-raw';
import rehypeKatex from 'rehype-katex';
import rehypeHighlight from 'rehype-highlight';
import 'katex/dist/katex.min.css';
import 'highlight.js/styles/github.css';
import styles from './MessageContent.module.css';

interface MessageContentProps {
    content: string;
}

const MessageContent: React.FC<MessageContentProps> = ({ content }) => {
    // 先把 <think>…</think> 的段落拆出来
    const parts = content.split(/<think>([\s\S]*?)<\/think>/g);
    // parts => [normal, think, normal, think, normal…]
    return (
        <>
            {parts.map((part, idx) => {
                // 偶数 idx 是普通文本
                if (idx % 2 === 0) {
                    return (
                        <ReactMarkdown
                            key={idx}
                            remarkPlugins={[remarkGfm, remarkMath]}
                            rehypePlugins={[rehypeKatex, rehypeHighlight, rehypeRaw]}
                        >
                            {part}
                        </ReactMarkdown>
                    );
                }
                // 奇数 idx 是 think 段
                return (
                    <details key={idx} className={styles.thinkContainer}>
                        <summary>思考过程</summary>
                        <div className={styles.thinkContent}>
                            <ReactMarkdown remarkPlugins={[remarkGfm]}>
                                {part}
                            </ReactMarkdown>
                        </div>
                    </details>
                );
            })}
        </>
    );
};

export default MessageContent;
