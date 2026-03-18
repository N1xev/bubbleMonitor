import Link from 'next/link';
import { Terminal } from 'lucide-react';

export default function Footer() {
  return (
    <footer className="w-full bg-zinc-950 p-8 md:p-12 border-t border-dashed border-slate-700">
      <div className="max-w-5xl mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-8">
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-2 text-slate-200">
            <Terminal className="w-5 h-5" />
            <span className="font-black tracking-widest text-lg uppercase">bubbleMonitor</span>
          </div>
          <span className="text-xs font-mono text-slate-500 uppercase tracking-widest">A minimal TUI system monitor.</span>
        </div>
        
        <div className="flex flex-col md:items-end gap-2 text-xs font-mono uppercase tracking-widest text-slate-500">
          <div>
            AGPLv3 License &copy; {new Date().getFullYear()}
          </div>
          <div>
            Made by <Link href="https://github.com/N1xev" target="_blank" className="text-slate-300 hover:text-slate-100 hover:underline underline-offset-4 decoration-dashed transition-colors font-bold">Alaa Elsamouly</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
