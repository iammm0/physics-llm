// ⚠️ 本文件为 React 组件示例，如需拆分 .tsx / .module.css 请自行调整。
// 已安装依赖：@uiw/react-json-view、copy-to-clipboard 等。

import React, { useMemo } from "react";
import ReactMarkdown from "react-markdown";
// -------- remark 插件（Markdown ⇢ mdast 阶段）------------------
import remarkGfm from "remark-gfm";                // GitHub 风格Markdown
import remarkMath from "remark-math";              // 支持 $LaTeX$ 公式
import remarkBreaks from "remark-breaks";          // 将换行符转换为 <br/>
import remarkFrontmatter from "remark-frontmatter"; // YAML front‑matter 解析
import remarkEmoji from "remark-emoji";            // :smile: → 😄
import remarkDirective from "remark-directive";    // :::note / :::warning 自定义块
// -------- rehype 插件（mdast ⇢ hast / HTML 阶段） ---------------
import rehypeRaw from "rehype-raw";                // 解析文档中的原生 HTML
import rehypeSanitize from "rehype-sanitize";      // XSS 安全防护
import rehypeKatex from "rehype-katex";            // 渲染 LaTeX
import rehypeHighlight from "rehype-highlight";    // highlight.js 语法高亮
import rehypeSlug from "rehype-slug";              // 给标题生成 id
import rehypeAutolinkHeadings from "rehype-autolink-headings"; // 标题锚点链接
import rehypeMermaid from "rehype-mermaid";        // 支持 ```mermaid``` 流程图

// --------- 其它工具 -------------------------------------------
import copyToClipboard from "copy-to-clipboard";
import JsonView from "@uiw/react-json-view"; // ← 新库，支持 React 18+

// --------- 样式 & 高亮主题 ------------------------------------
import "katex/dist/katex.min.css";                 // Katex 样式
import "highlight.js/styles/github.css";          // highlight.js 主题
import styles from "./MessageContent.module.css";

// ------------------- 组件 Props 定义 ---------------------------
interface MessageContentProps {
    /** LLM 返回的原始字符串 */
    content: string;
    /** 是否允许渲染原生 HTML（跳过 sanitize） */
    allowUnsafeHtml?: boolean;
}

/**
 * 判断字符串是否为有效 JSON。
 * 为避免频繁 try/catch，仅在首尾字符形似 JSON 时尝试解析。
 */
const isJson = (input: string): boolean => {
    if (!input) return false;
    const trimmed = input.trim();
    const looksLikeJson =
        (trimmed.startsWith("{") && trimmed.endsWith("}")) ||
        (trimmed.startsWith("[") && trimmed.endsWith("]"));
    if (!looksLikeJson) return false;
    try {
        JSON.parse(trimmed);
        return true;
    } catch {
        return false;
    }
};

const MessageContent: React.FC<MessageContentProps> = ({
                                                           content,
                                                           allowUnsafeHtml = false,
                                                       }) => {
    // 1️⃣ 将 <think>…</think> 拆成奇偶片段，奇数为折叠的思考过程
    const parts = useMemo(() => content.split(/<think>([\s\S]*?)<\/think>/g), [content]);

    // remark / rehype 插件数组只在首次渲染时创建，避免重复实例化
    const remarkPlugins = useMemo(
        () => [
            remarkGfm,
            remarkMath,
            remarkBreaks,
            remarkFrontmatter,
            remarkEmoji,
            remarkDirective,
        ],
        []
    );

    const rehypePlugins = useMemo(() => {
        const base: any[] = [
            rehypeKatex,
            rehypeHighlight,
            rehypeSlug,
            rehypeAutolinkHeadings,
            rehypeMermaid,
            rehypeRaw, // ⚠️ rehypeRaw 必须放在 sanitize 之前
        ];
        if (!allowUnsafeHtml) base.push(rehypeSanitize);
        return base;
    }, [allowUnsafeHtml]);

    /**
     * 自定义渲染器：给代码块添加「复制」按钮
     * - 针对 <pre> 而非 <code>，因为 rehype-highlight 会输出 <pre><code>
     */
    const markdownComponents = useMemo(
        () => ({
            pre({ node, children, ...rest }: any) {
                // 提取纯文本代码，用于复制
                const rawCode = node.children?.[0]?.value || "";
                return (
                    <div className={styles.codeWrapper}>
                        <button
                            className={styles.copyBtn}
                            onClick={() => copyToClipboard(rawCode)}
                        >
                            复制
                        </button>
                        <pre {...rest}>{children}</pre>
                    </div>
                );
            },
        }),
        []
    );

    // --------------------------- 渲染 ---------------------------
    return (
        <>
            {parts.map((part, idx) => {
                // 偶数索引：正常内容
                if (idx % 2 === 0) {
                    // 若为 JSON，则使用 JsonView 美化展示 (新 API: value)
                    if (isJson(part)) {
                        return (
                            <JsonView
                                key={idx}
                                value={JSON.parse(part)}
                                collapsed={2}
                                /** 取消 keyName 显示，保持纯粹树结构 */
                                keyName={null as unknown as string}
                            />
                        );
                    }
                    // 否则走 Markdown 渲染管线
                    return (
                        <ReactMarkdown
                            key={idx}
                            remarkPlugins={remarkPlugins as any}
                            rehypePlugins={rehypePlugins as any}
                            components={markdownComponents as any}
                        >
                            {part}
                        </ReactMarkdown>
                    );
                }
                // 奇数索引：<think> 折叠区
                return (
                    <details key={idx} className={styles.thinkContainer}>
                        <summary>思考过程</summary>
                        <div className={styles.thinkContent}>
                            <ReactMarkdown remarkPlugins={[remarkGfm]}>{part}</ReactMarkdown>
                        </div>
                    </details>
                );
            })}
        </>
    );
};

export default MessageContent;