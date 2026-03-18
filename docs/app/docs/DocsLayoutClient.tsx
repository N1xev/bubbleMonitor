'use client';

import { useState, useEffect } from 'react';
import { usePathname } from 'next/navigation';
import { Menu, X, ChevronRight } from 'lucide-react';
import Link from 'next/link';
import TableOfContents from '@/components/TableOfContents';

const navLinks = [
  { href: '/docs', label: 'Introduction' },
  { href: '/docs/installation', label: 'Installation' },
  { href: '/docs/quick-start', label: 'Quick Start' },
  { href: '/docs/configuration', label: 'The config.json' },
  { href: '/docs/theming', label: 'Custom Theming' },
  { href: '/docs/remote-hosts', label: 'Remote Hosts (SSH)' },
];

const pageTitles: Record<string, string> = {
  '/docs': 'Introduction',
  '/docs/installation': 'Installation',
  '/docs/quick-start': 'Quick Start',
  '/docs/configuration': 'The config.json',
  '/docs/theming': 'Custom Theming',
  '/docs/remote-hosts': 'Remote Hosts (SSH)',
};

export default function DocsLayoutClient({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [isPageTransitioning, setIsPageTransitioning] = useState(false);
  const pathname = usePathname();

  useEffect(() => {
    setIsPageTransitioning(true);
    const timer = setTimeout(() => {
      setIsPageTransitioning(false);
      setIsSidebarOpen(false);
    }, 150);
    return () => clearTimeout(timer);
  }, [pathname]);

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  const currentIndex = navLinks.findIndex(link => link.href === pathname);
  const prevLink = currentIndex > 0 ? navLinks[currentIndex - 1] : null;
  const nextLink = currentIndex < navLinks.length - 1 ? navLinks[currentIndex + 1] : null;

  const currentTitle = pageTitles[pathname] || 'Documentation';

  return (
    <div className="min-h-screen flex flex-col font-sans text-slate-300 selection:bg-slate-300 selection:text-zinc-950 bg-zinc-950 relative">
      <header className="sticky top-0 z-50 w-full border-b border-dashed border-slate-700 bg-zinc-950">
        <div className="flex items-center justify-between px-4 lg:px-8 h-14">
          <div className="flex items-center gap-4">
            <button 
              onClick={toggleSidebar}
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
            <div className="hidden lg:flex items-center gap-2">
              <ChevronRight className="w-4 h-4 text-slate-600" />
              <span className="text-xs font-mono uppercase tracking-widest text-slate-400">Documentation</span>
              <ChevronRight className="w-4 h-4 text-slate-600" />
              <span className="text-xs font-mono uppercase tracking-widest text-emerald-400">{currentTitle}</span>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <Link href="/themes" className="hidden sm:block text-xs font-mono uppercase tracking-widest text-slate-400 hover:text-slate-100 transition-colors">
              Themes
            </Link>
            <Link href="/changelog" className="hidden sm:block text-xs font-mono uppercase tracking-widest text-slate-400 hover:text-slate-100 transition-colors">
              Changelog
            </Link>
            <Link href="/docs" className="text-xs font-mono uppercase tracking-widest text-emerald-400 hover:text-emerald-300 transition-colors">
              Docs
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
          <nav className="flex-1 overflow-y-auto p-6 space-y-8 font-mono text-xs uppercase tracking-widest">
            <div>
              <h3 className="text-slate-500 mb-4 font-bold">Getting Started</h3>
              <ul className="space-y-3 border-l border-dashed border-slate-700 pl-4">
                {navLinks.slice(0, 3).map((link) => (
                  <li key={link.href}>
                    <Link 
                      href={link.href} 
                      className={`block transition-colors ${
                        pathname === link.href 
                          ? 'text-emerald-400' 
                          : 'text-slate-300 hover:text-white'
                      }`}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            
            <div>
              <h3 className="text-slate-500 mb-4 font-bold">Configuration</h3>
              <ul className="space-y-3 border-l border-dashed border-slate-700 pl-4">
                {navLinks.slice(3).map((link) => (
                  <li key={link.href}>
                    <Link 
                      href={link.href} 
                      className={`block transition-colors ${
                        pathname === link.href 
                          ? 'text-emerald-400' 
                          : 'text-slate-300 hover:text-white'
                      }`}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
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
          <nav className="flex-1 overflow-y-auto p-6 space-y-8 font-mono text-xs uppercase tracking-widest">
            <div>
              <h3 className="text-slate-500 mb-4 font-bold">Getting Started</h3>
              <ul className="space-y-3 border-l border-dashed border-slate-700 pl-4">
                {navLinks.slice(0, 3).map((link) => (
                  <li key={link.href}>
                    <Link 
                      href={link.href} 
                      onClick={() => setIsSidebarOpen(false)}
                      className={`block transition-colors ${
                        pathname === link.href 
                          ? 'text-emerald-400' 
                          : 'text-slate-300 hover:text-white'
                      }`}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            
            <div>
              <h3 className="text-slate-500 mb-4 font-bold">Configuration</h3>
              <ul className="space-y-3 border-l border-dashed border-slate-700 pl-4">
                {navLinks.slice(3).map((link) => (
                  <li key={link.href}>
                    <Link 
                      href={link.href} 
                      onClick={() => setIsSidebarOpen(false)}
                      className={`block transition-colors ${
                        pathname === link.href 
                          ? 'text-emerald-400' 
                          : 'text-slate-300 hover:text-white'
                      }`}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </nav>
        </aside>

        <main className={`flex-1 min-w-0 flex-1 transition-all duration-150 ${isPageTransitioning ? 'opacity-0 translate-x-2' : 'opacity-100 translate-x-0'}`}>
          <div className="px-6 py-8 lg:px-12 lg:py-12 max-w-4xl mx-auto w-full">
            {children}
            
            <div className="mt-16 pt-8 border-t border-dashed border-slate-700 grid grid-cols-1 sm:grid-cols-2 gap-4">
              {prevLink ? (
                <Link href={prevLink.href} className="group flex flex-col text-left p-4 border border-dashed border-slate-700 bg-zinc-900/30 hover:bg-zinc-900 transition-all hover:border-emerald-400/50">
                  <span className="text-[10px] font-mono uppercase tracking-widest text-slate-500 mb-1">Previous</span>
                  <span className="text-base font-bold text-slate-300 group-hover:text-emerald-400 transition-colors">{prevLink.label}</span>
                </Link>
              ) : <div />}
              
              {nextLink ? (
                <Link href={nextLink.href} className="group flex flex-col text-right p-4 border border-dashed border-slate-700 bg-zinc-900/30 hover:bg-zinc-900 transition-all hover:border-emerald-400/50">
                  <span className="text-[10px] font-mono uppercase tracking-widest text-slate-500 mb-1">Next</span>
                  <span className="text-base font-bold text-slate-300 group-hover:text-emerald-400 transition-colors">{nextLink.label}</span>
                </Link>
              ) : <div />}
            </div>
          </div>
          
          <TableOfContents />
        </main>
      </div>
    </div>
  );
}
