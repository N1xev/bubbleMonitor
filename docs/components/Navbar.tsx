'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Github, Download, Menu, X } from 'lucide-react';

export default function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const navLinks = [
    { href: '/docs', label: 'Docs' },
    { href: '/themes', label: 'Themes' },
    { href: '/roadmap', label: 'Roadmap' },
    { href: '/changelog', label: 'Changelog' },
  ];

  return (
    <>
      <header className="sticky top-0 z-50 w-full border-b border-dashed border-slate-700 bg-zinc-950">
        <div className="mx-auto px-4 sm:px-8 py-3 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-4">
            <div className="w-10 h-10 bg-slate-200 text-zinc-950 flex items-center justify-center font-black font-mono text-lg shadow-[4px_4px_0px_0px_rgba(51,65,85,1)]">
              {'>_'}
            </div>
            <div className="flex flex-col justify-center">
              <span className="font-sans tracking-widest text-sm uppercase text-slate-100 font-black">bubbleMonitor</span>
              <span className="font-mono text-[10px] text-slate-500 uppercase tracking-widest">bub-v0.1.4</span>
            </div>
          </Link>
          
          <div className="hidden md:flex items-center gap-6 text-xs font-mono uppercase tracking-widest">
            <div className="flex items-center gap-2 px-3 py-1.5 border border-dashed border-slate-700 bg-zinc-900/50 text-slate-400">
              <span className="w-2 h-2 bg-emerald-500"></span>
              Build Status: <span className="text-emerald-400 font-bold">Passing</span>
            </div>
            <Link href="/docs" className="text-slate-400 hover:text-slate-100 transition-colors flex items-center gap-2 px-2">
              Docs
            </Link>
            <Link href="/themes" className="text-slate-400 hover:text-slate-100 transition-colors flex items-center gap-2 px-2">
              Themes
            </Link>
            <Link href="/roadmap" className="text-slate-400 hover:text-slate-100 transition-colors flex items-center gap-2 px-2">
              Roadmap
            </Link>
            <Link href="/changelog" className="text-slate-400 hover:text-slate-100 transition-colors flex items-center gap-2 px-2">
              Changelog
            </Link>
            <Link href="https://github.com/N1xev/bubbleMonitor" target="_blank" className="text-slate-400 hover:text-slate-100 transition-colors flex items-center gap-2 px-2">
              <Github className="w-4 h-4" /> Source
            </Link>
            <Link href="/#download" className="px-5 py-2 bg-slate-200 text-zinc-950 hover:bg-white transition-all flex items-center gap-2 font-bold shadow-[4px_4px_0px_0px_rgba(51,65,85,1)] hover:translate-y-[2px] hover:translate-x-[2px] hover:shadow-[2px_2px_0px_0px_rgba(51,65,85,1)] active:translate-y-[4px] active:translate-x-[4px] active:shadow-none">
              <Download className="w-4 h-4" /> Download
            </Link>
          </div>

          <button 
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="md:hidden p-2 text-slate-400 hover:text-white transition-colors"
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </header>

      {isMobileMenuOpen && (
        <>
          <div 
            className="fixed inset-0 z-40 bg-black/50 md:hidden"
            onClick={() => setIsMobileMenuOpen(false)}
          />
          <div className="fixed top-[57px] right-0 z-50 w-64 bg-zinc-950 border-l border-dashed border-slate-700 shadow-xl md:hidden">
            <nav className="p-4 space-y-2 font-mono text-xs uppercase tracking-widest">
              {navLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-400 hover:text-white hover:bg-zinc-900 transition-colors"
                >
                  {link.label}
                </Link>
              ))}
              <Link
                href="https://github.com/N1xev/bubbleMonitor"
                target="_blank"
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center gap-3 px-4 py-3 rounded-lg text-slate-400 hover:text-white hover:bg-zinc-900 transition-colors"
              >
                <Github className="w-4 h-4" /> Source
              </Link>
              <Link
                href="/#download"
                onClick={() => setIsMobileMenuOpen(false)}
                className="flex items-center gap-3 px-4 py-3 rounded-lg bg-slate-200 text-zinc-950 font-bold"
              >
                <Download className="w-4 h-4" /> Download
              </Link>
            </nav>
          </div>
        </>
      )}
    </>
  );
}
