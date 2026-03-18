'use client';

import { useState } from 'react';
import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import CodeBlock from '@/components/CodeBlock';
import { Palette, ChevronLeft, ChevronRight } from 'lucide-react';

export default function ThemesPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const themesPerPage = 4;

  const themes = [
    {
      name: "Dracula",
      description: "A dark theme for many editors, shells, and more.",
      colors: {
        primary: "#bd93f9",
        secondary: "#ff79c6",
        success: "#50fa7b",
        warning: "#f1fa8c",
        alert: "#ff5555",
        text: "#f8f8f2",
        muted: "#6272a4",
        border: "#44475a",
        background: "#282a36"
      }
    },
    {
      name: "Catppuccin Mocha",
      description: "Soothing pastel theme for the high-spirited!",
      colors: {
        primary: "#cba6f7",
        secondary: "#f5c2e7",
        success: "#a6e3a1",
        warning: "#f9e2af",
        alert: "#f38ba8",
        text: "#cdd6f4",
        muted: "#6c7086",
        border: "#45475a",
        background: "#1e1e2e"
      }
    },
    {
      name: "Nord",
      description: "An arctic, north-bluish color palette.",
      colors: {
        primary: "#88C0D0",
        secondary: "#81A1C1",
        success: "#A3BE8C",
        warning: "#EBCB8B",
        alert: "#BF616A",
        text: "#ECEFF4",
        muted: "#4C566A",
        border: "#3B4252",
        background: "#2E3440"
      }
    },
    {
      name: "Gruvbox Dark",
      description: "Retro groove color scheme for Vim.",
      colors: {
        primary: "#83a598",
        secondary: "#d3869b",
        success: "#b8bb26",
        warning: "#fabd2f",
        alert: "#fb4934",
        text: "#ebdbb2",
        muted: "#928374",
        border: "#504945",
        background: "#282828"
      }
    },
    {
      name: "Tokyo Night",
      description: "A clean, dark theme that celebrates the lights of Downtown Tokyo at night.",
      colors: {
        primary: "#7aa2f7",
        secondary: "#bb9af7",
        success: "#9ece6a",
        warning: "#e0af68",
        alert: "#f7768e",
        text: "#c0caf5",
        muted: "#565f89",
        border: "#414868",
        background: "#1a1b26"
      }
    },
    {
      name: "Monokai Pro",
      description: "Professional theme with color-coordinated UI.",
      colors: {
        primary: "#ab9df2",
        secondary: "#fc9867",
        success: "#a9dc76",
        warning: "#ffd866",
        alert: "#ff6188",
        text: "#fcfcfa",
        muted: "#727072",
        border: "#403e41",
        background: "#2d2a2e"
      }
    },
    {
      name: "Solarized Dark",
      description: "Precision colors for machines and people.",
      colors: {
        primary: "#268bd2",
        secondary: "#d33682",
        success: "#859900",
        warning: "#b58900",
        alert: "#dc322f",
        text: "#839496",
        muted: "#586e75",
        border: "#073642",
        background: "#002b36"
      }
    },
    {
      name: "One Dark",
      description: "Atom's iconic dark theme.",
      colors: {
        primary: "#61afef",
        secondary: "#c678dd",
        success: "#98c379",
        warning: "#e5c07b",
        alert: "#e06c75",
        text: "#abb2bf",
        muted: "#5c6370",
        border: "#3e4451",
        background: "#282c34"
      }
    },
    {
      name: "Synthwave '84",
      description: "Do you remember the 80s? This theme does.",
      colors: {
        primary: "#36f9f6",
        secondary: "#f92aad",
        success: "#72f1b8",
        warning: "#fede5d",
        alert: "#fe4450",
        text: "#f92aad",
        muted: "#848bbd",
        border: "#495495",
        background: "#262335"
      }
    },
    {
      name: "Rosé Pine",
      description: "All natural pine, faux fur and a bit of soho vibes.",
      colors: {
        primary: "#c4a7e7",
        secondary: "#ebbcba",
        success: "#9ccfd8",
        warning: "#f6c177",
        alert: "#eb6f92",
        text: "#e0def4",
        muted: "#6e6a86",
        border: "#44415a",
        background: "#191724"
      }
    },
    {
      name: "Cyberpunk",
      description: "High tech, low life.",
      colors: {
        primary: "#0abdc6",
        secondary: "#ea00d9",
        success: "#00ff00",
        warning: "#ffb800",
        alert: "#ff003c",
        text: "#d3d7cf",
        muted: "#711c91",
        border: "#133e7c",
        background: "#000b18"
      }
    },
    {
      name: "Ayu Dark",
      description: "A simple theme with bright colors.",
      colors: {
        primary: "#39bae6",
        secondary: "#f07178",
        success: "#c2d94c",
        warning: "#ffb454",
        alert: "#ff3333",
        text: "#b3b1ad",
        muted: "#e6e1cf",
        border: "#3e4b59",
        background: "#0a0e14"
      }
    }
  ];

  const totalPages = Math.ceil(themes.length / themesPerPage);
  const indexOfLastTheme = currentPage * themesPerPage;
  const indexOfFirstTheme = indexOfLastTheme - themesPerPage;
  const currentThemes = themes.slice(indexOfFirstTheme, indexOfLastTheme);

  const nextPage = () => {
    if (currentPage < totalPages) setCurrentPage(prev => prev + 1);
  };

  const prevPage = () => {
    if (currentPage > 1) setCurrentPage(prev => prev - 1);
  };

  return (
    <div className="min-h-screen flex flex-col font-sans text-slate-300 selection:bg-slate-300 selection:text-zinc-950 relative overflow-x-hidden">
      <Navbar />

      <main className="flex-1 w-full mx-auto">
        <section className="p-8 md:p-16 lg:p-24 border-b border-dashed border-slate-700 bg-zinc-950/90 relative z-20">
          <div className="absolute top-4 left-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute top-4 right-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute bottom-4 left-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute bottom-4 right-4 text-slate-700 font-mono text-xs">+</div>

          <div className="max-w-5xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-400 mb-8 bg-zinc-900">
              <Palette className="w-3 h-3 text-slate-300" />
              Theme Gallery
            </div>
            
            <h1 className="text-5xl md:text-7xl lg:text-8xl font-black tracking-tighter text-slate-100 mb-8 leading-[0.9] uppercase">
              Make it <br />
              <span className="text-slate-500">yours.</span>
            </h1>
            
            <p className="text-slate-400 text-lg md:text-xl max-w-2xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
              Browse community and built-in themes. Copy the JSON configuration directly into your config file to instantly change the look of bubbleMonitor.
            </p>
          </div>
        </section>

        <section className="p-8 md:p-16 lg:p-24 bg-zinc-950/50">
          <div className="max-w-5xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12 mb-12">
              {currentThemes.map((theme, idx) => (
                <div key={idx} className="border border-dashed border-slate-700 bg-zinc-900/30 flex flex-col">
                  <div className="p-6 border-b border-dashed border-slate-700">
                    <h2 className="text-2xl font-bold text-slate-100 uppercase tracking-tight mb-2">{theme.name}</h2>
                    <p className="text-sm text-slate-400">{theme.description}</p>
                  </div>
                  
                  <div className="p-6 flex-1">
                    <div className="grid grid-cols-3 gap-3 mb-6">
                      {Object.entries(theme.colors).map(([key, value]) => (
                        <div key={key} className="flex flex-col gap-1">
                          <div className="h-8 w-full border border-dashed border-slate-700" style={{ backgroundColor: value }}></div>
                          <span className="text-[10px] font-mono text-slate-500 uppercase">{key}</span>
                        </div>
                      ))}
                    </div>

                    <CodeBlock 
                      code={JSON.stringify({ theme: "custom", custom_theme: theme.colors }, null, 2)} 
                      lang="json" 
                      filename="CONFIG.JSON"
                      wrapperClassName="relative border border-dashed border-slate-700 bg-zinc-950"
                      className="!p-4 !text-xs"
                      showCopy={true}
                    />
                  </div>
                </div>
              ))}
            </div>

            {/* Pagination Controls */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-4 mt-16">
                <button 
                  onClick={prevPage}
                  disabled={currentPage === 1}
                  className="p-3 border border-dashed border-slate-700 bg-zinc-900 hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <ChevronLeft className="w-5 h-5 text-slate-300" />
                </button>
                
                <div className="font-mono text-sm text-slate-400">
                  PAGE <span className="text-slate-200">{currentPage}</span> OF <span className="text-slate-200">{totalPages}</span>
                </div>
                
                <button 
                  onClick={nextPage}
                  disabled={currentPage === totalPages}
                  className="p-3 border border-dashed border-slate-700 bg-zinc-900 hover:bg-zinc-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <ChevronRight className="w-5 h-5 text-slate-300" />
                </button>
              </div>
            )}
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
