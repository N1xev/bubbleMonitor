'use client';

import { useEffect, useState } from 'react';
import { usePathname } from 'next/navigation';
import { List } from 'lucide-react';

interface Heading {
  id: string;
  text: string;
  level: number;
}

export default function TableOfContents() {
  const pathname = usePathname();
  const [headings, setHeadings] = useState<Heading[]>([]);
  const [activeId, setActiveId] = useState<string>('');
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    setHeadings([]);
    setActiveId('');
    setIsOpen(false);
    
    const timer = setTimeout(() => {
      const elements = Array.from(document.querySelectorAll('main h2, main h3'))
        .filter(element => element.id);
      
      const headingData = elements.map(element => ({
        id: element.id,
        text: element.textContent || '',
        level: Number(element.tagName.substring(1))
      }));

      setHeadings(headingData);

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              setActiveId(entry.target.id);
            }
          });
        },
        { rootMargin: '0px 0px -80% 0px' }
      );

      elements.forEach((element) => observer.observe(element));

      return () => observer.disconnect();
    }, 100);

    return () => clearTimeout(timer);
  }, [pathname]);

  if (headings.length === 0) return null;

  return (
    <>
      <div className="fixed bottom-6 right-6 z-50 md:hidden">
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="w-12 h-12 rounded-full bg-zinc-900 border border-dashed border-slate-700 flex items-center justify-center text-slate-400 hover:text-emerald-400 hover:border-emerald-400/50 transition-all shadow-lg"
        >
          <List className="w-5 h-5" />
        </button>
      </div>

      {isOpen && (
        <div className="fixed bottom-20 right-6 z-50 md:hidden bg-zinc-900 border border-dashed border-slate-700 rounded-lg p-4 max-h-[60vh] overflow-y-auto shadow-xl w-64">
          <h4 className="text-[10px] font-mono uppercase tracking-widest text-slate-500 font-bold mb-3">On this page</h4>
          <ul className="space-y-2 text-xs font-mono">
            {headings.map((heading) => (
              <li 
                key={heading.id} 
                className={`${heading.level === 3 ? 'ml-3' : ''}`}
              >
                <a 
                  href={`#${heading.id}`}
                  onClick={() => setIsOpen(false)}
                  className={`block transition-colors truncate ${
                    activeId === heading.id ? 'text-emerald-400' : 'text-slate-400 hover:text-white'
                  }`}
                >
                  {heading.text}
                </a>
              </li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
}
