'use client';

import { useEffect, useState } from 'react';
import { codeToHtml } from 'shiki';
import { Check, Copy } from 'lucide-react';

export default function CodeBlock({ code, lang, filename, className = "", wrapperClassName = "my-6 border border-dashed border-slate-700 bg-zinc-950", showCopy = false, showHeader = true }: { code: string, lang: string, filename?: string, className?: string, wrapperClassName?: string, showCopy?: boolean, showHeader?: boolean }) {
  const [html, setHtml] = useState('');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    codeToHtml(code, {
      lang,
      theme: 'vitesse-dark',
    }).then(setHtml);
  }, [code, lang]);

  const handleCopy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const copyButton = showCopy && (
    <button
      onClick={handleCopy}
      className="absolute top-3 right-3 p-2 bg-zinc-900 hover:bg-zinc-800 text-slate-400 hover:text-slate-200 border border-dashed border-slate-600 opacity-0 group-hover:opacity-100 transition-all z-10"
      title="Copy to clipboard"
    >
      {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
    </button>
  );

  const terminalHeader = (
    <div className="flex items-center justify-between pl-4 pr-0 py-0 border-b border-dashed border-slate-700 bg-zinc-900/50">
      <div className="flex gap-2 py-2">
        <div className="w-3 h-3 rounded-full bg-slate-700"></div>
        <div className="w-3 h-3 rounded-full bg-slate-700"></div>
        <div className="w-3 h-3 rounded-full bg-slate-700"></div>
      </div>
      <div className="px-3 py-2 bg-slate-800 text-[10px] font-mono text-slate-300 border-l border-dashed border-slate-700 uppercase tracking-widest">
        {filename || lang}
      </div>
    </div>
  );

  if (!html) {
    return (
      <div className={`relative group ${wrapperClassName}`}>
        {showHeader && terminalHeader}
        <div className="relative">
          {copyButton}
          <pre className={`p-4 overflow-x-auto text-sm font-mono text-slate-300 leading-relaxed ${className}`}><code>{code}</code></pre>
        </div>
      </div>
    );
  }

  return (
    <div className={`relative group ${wrapperClassName}`}>
      {showHeader && terminalHeader}
      <div className="relative">
        {copyButton}
        <div className={`p-4 overflow-x-auto text-sm font-mono leading-relaxed [&>pre]:!bg-transparent [&>pre]:!m-0 [&>pre]:!p-0 ${className}`} dangerouslySetInnerHTML={{ __html: html }} />
      </div>
    </div>
  );
}
