import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import { Terminal, CheckSquare, Square, GitPullRequest } from 'lucide-react';

export default function RoadmapPage() {
  const roadmap = [
    {
      version: "v1.0.0",
      status: "planned",
      title: "The Stable Release",
      tasks: [
        { done: false, text: "Full plugin system for community modules" },
        { done: false, text: "Web dashboard for remote viewing" },
        { done: false, text: "Comprehensive unit test coverage (>90%)" },
        { done: false, text: "Official Docker image and Helm chart" }
      ]
    },
    {
      version: "v0.2.0",
      status: "in-progress",
      title: "The Network Update",
      tasks: [
        { done: true, text: "Rewrite network module for lower latency" },
        { done: true, text: "Add support for monitoring remote hosts via SSH" },
        { done: false, text: "Per-process network bandwidth monitoring (eBPF)" },
        { done: false, text: "Active connection tracking and port mapping" }
      ]
    },
    {
      version: "v0.1.5",
      status: "completed",
      title: "The Theming Update",
      tasks: [
        { done: true, text: "Braille chart support for CPU and Memory metrics" },
        { done: true, text: "Custom theming engine with JSON config" },
        { done: true, text: "Resolve memory leak in remote host monitoring" },
        { done: true, text: "Optimize process list rendering" }
      ]
    }
  ];

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
              <GitPullRequest className="w-3 h-3 text-slate-300" />
              Project Roadmap
            </div>
            
            <h1 className="text-5xl md:text-7xl lg:text-8xl font-black tracking-tighter text-slate-100 mb-8 leading-[0.9] uppercase">
              What&apos;s <br />
              <span className="text-slate-500">Next.</span>
            </h1>
            
            <p className="text-slate-400 text-lg md:text-xl max-w-2xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
              A raw look at our upcoming features, planned improvements, and current progress.
            </p>
          </div>
        </section>

        <section className="p-8 md:p-16 lg:p-24 bg-zinc-950/50">
          <div className="max-w-4xl mx-auto space-y-12">
            {roadmap.map((phase, idx) => (
              <div key={idx} className="border border-dashed border-slate-700 bg-zinc-900/30 p-8 md:p-12 relative">
                <div className="absolute top-4 right-4 text-[10px] font-mono text-slate-600 uppercase">
                  STATUS: {phase.status.replace('-', '_')}
                </div>
                
                <div className="mb-8 border-b border-dashed border-slate-700 pb-6">
                  <h2 className="text-4xl md:text-5xl font-black tracking-tighter text-slate-100 uppercase mb-2">
                    {phase.version}
                  </h2>
                  <p className="text-lg text-slate-400 font-mono tracking-tight">{phase.title}</p>
                </div>
                
                <ul className="space-y-4 font-mono text-sm">
                  {phase.tasks.map((task, i) => (
                    <li key={i} className="flex items-start gap-4 group">
                      <span className="mt-0.5 shrink-0 text-slate-500 group-hover:text-slate-300 transition-colors">
                        {task.done ? <CheckSquare className="w-5 h-5 text-emerald-400" /> : <Square className="w-5 h-5" />}
                      </span>
                      <span className={`${task.done ? 'text-slate-500 line-through decoration-slate-600' : 'text-slate-200'} leading-relaxed`}>
                        {task.text}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
