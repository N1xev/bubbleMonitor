'use client';

import { useState } from 'react';
import { usePathname } from 'next/navigation';
import { Menu, X, ChevronRight, BookOpen, Palette, Map, FileText } from 'lucide-react';
import Link from 'next/link';

const rootNavLinks = [
  { href: '/', label: 'Home', icon: null },
  { href: '/docs', label: 'Docs', icon: BookOpen },
  { href: '/themes', label: 'Themes', icon: Palette },
  { href: '/changelog', label: 'Changelog', icon: FileText },
  { href: '/roadmap', label: 'Roadmap', icon: Map },
];

export default function RootLayoutClient({ children, title }: { children: React.ReactNode; title?: string }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const pathname = usePathname();

  return (
    <div className="min-h-screen flex flex-col font-sans text-slate-300 selection:bg-slate-300 selection:text-zinc-950 bg-zinc-950 relative">
      <header className="sticky top-0 z-50 w-full border-b border-dashed border-slate-700 bg-zinc-950">
        <div className="flex items-center justify-between px-4 lg:px-8 h-14">
          <div className="flex items-center gap-4">
            <button 
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="lg:hidden p-2 text-slate-400 hover:text-white transition-colors"
            >
              {isSidebarOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
            <Link href="/" className="flex items-center gap-3">
              <div className="w-8 h-8 bg-slate-200 text-zinc-950 flex items-center justify-center font-black font-mono text-sm shadow-[3px_3px_0px_0px_rgba(51,65,85,1)]">
                {'>'}_
              </div>
              <div className="hidden sm:flex flex-col justify-center">
                <span className="font-sans tracking-widest text-xs uppercase text-slate-100 font-black">bubbleMonitor</span>
              </div>
            </Link>
            {title && (
              <>
                <ChevronRight className="w-4 h-4 text-slate-600 hidden lg:block" />
                <span className="hidden lg:block text-xs font-mono uppercase tracking-widest text-emerald-400">{title}</span>
              </>
            )}
          </div>
          <div className="flex items-center gap-4">
            <Link href="/docs" className="text-xs font-mono uppercase tracking-widest text-slate-400 hover:text-slate-100 transition-colors">
              Docs
            </Link>
            <Link href="/themes" className="hidden sm:block text-xs font-mono uppercase tracking-widest text-slate-400 hover:text-slate-100 transition-colors">
              Themes
            </Link>
            <Link href="/changelog" className="hidden sm:block text-xs font-mono uppercase tracking-widest text-slate-400 hover:text-slate-100 transition-colors">
              Changelog
            </Link>
          </div>
        </div>
      </header>

      <div className="flex-1 w-full max-w-screen-2xl mx-auto flex relative">
        <aside className={`
          hidden lg:flex flex-col
          w-64 h-[calc(100vh-3.5rem)] sticky top-14
          border-r border-dashed border-slate-700
          bg-zinc-950
          p-6
        `}>
          <nav className="space-y-2 font-mono text-xs uppercase tracking-widest">
            {rootNavLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  pathname === link.href
                    ? 'bg-emerald-400/10 text-emerald-400 border border-emerald-400/30'
                    : 'text-slate-400 hover:text-white hover:bg-zinc-900'
                }`}
              >
                {link.icon && <link.icon className="w-4 h-4" />}
                {link.label}
              </Link>
            ))}
          </nav>
        </aside>

        {isSidebarOpen && (
          <div 
            className="fixed inset-0 z-40 bg-black/50 lg:hidden"
            onClick={() => setIsSidebarOpen(false)}
          />
        )}

        <aside className={`
          fixed lg:hidden top-14 left-0 z-40
          w-72 h-[calc(100vh-3.5rem)]
          bg-zinc-950 border-r border-dashed border-slate-700
          transition-transform duration-300 ease-out
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        `}>
          <nav className="p-6 space-y-2 font-mono text-xs uppercase tracking-widest">
            {rootNavLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                onClick={() => setIsSidebarOpen(false)}
                className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  pathname === link.href
                    ? 'bg-emerald-400/10 text-emerald-400 border border-emerald-400/30'
                    : 'text-slate-400 hover:text-white hover:bg-zinc-900'
                }`}
              >
                {link.icon && <link.icon className="w-4 h-4" />}
                {link.label}
              </Link>
            ))}
          </nav>
        </aside>

        <main className="flex-1 min-w-0">
          {children}
        </main>
      </div>
    </div>
  );
}
