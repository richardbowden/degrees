import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Components } from 'react-markdown';

const components: Components = {
  h2: ({ children }) => (
    <h2 className="text-xl font-bold text-white mt-8 mb-3 flex items-center gap-2">
      <span className="w-1 h-6 bg-brand-500 rounded-full inline-block" />
      {children}
    </h2>
  ),
  h3: ({ children }) => (
    <h3 className="text-lg font-semibold text-white mt-6 mb-2">{children}</h3>
  ),
  p: ({ children }) => (
    <p className="text-text-secondary leading-relaxed mb-4">{children}</p>
  ),
  ul: ({ children }) => (
    <ul className="space-y-2 mb-4">{children}</ul>
  ),
  li: ({ children }) => (
    <li className="flex items-start gap-2 text-text-secondary">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-4 h-4 text-brand-400 mt-1 shrink-0">
        <path strokeLinecap="round" strokeLinejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
      </svg>
      <span>{children}</span>
    </li>
  ),
  strong: ({ children }) => (
    <strong className="text-white font-semibold">{children}</strong>
  ),
  table: ({ children }) => (
    <div className="overflow-x-auto mb-4 rounded-lg border border-border-subtle">
      <table className="w-full text-sm">{children}</table>
    </div>
  ),
  thead: ({ children }) => (
    <thead className="bg-white/5 text-text-secondary">{children}</thead>
  ),
  th: ({ children }) => (
    <th className="px-4 py-2 text-left font-medium border-b border-border-subtle">{children}</th>
  ),
  td: ({ children }) => (
    <td className="px-4 py-2 text-text-secondary border-b border-border-subtle">{children}</td>
  ),
  blockquote: ({ children }) => (
    <blockquote className="border-l-2 border-brand-500 pl-4 italic text-text-muted mb-4">{children}</blockquote>
  ),
};

export function MarkdownContent({ content }: { content: string }) {
  return (
    <div className="markdown-content">
      <ReactMarkdown remarkPlugins={[remarkGfm]} components={components}>
        {content}
      </ReactMarkdown>
    </div>
  );
}
