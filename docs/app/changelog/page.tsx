import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import CodeBlock from '@/components/CodeBlock';
import { Terminal, GitCommit, GitMerge, Star, Zap } from 'lucide-react';

export default function ChangelogPage() {
  const releases = [
    {
      version: "v0.1.4",
      date: "March 8, 2026",
      title: "The Theming & Braille Update",
      description: "This release introduces a completely overhauled rendering engine that supports braille characters for ultra-precise charting, alongside a robust custom theming system. We've also squashed a few memory leaks related to remote host monitoring.",
      type: "minor",
      changes: [
        { type: "feat", title: "Braille Chart Support", text: "Added braille chart support for CPU and Memory metrics, increasing visual fidelity by 4x in the terminal." },
        { type: "feat", title: "Custom Theming Engine", text: "New custom theming engine with 30+ built-in themes. You can now define primary, secondary, and accent colors in your config.json." },
        { type: "fix", title: "Remote Host Memory Leak", text: "Resolved memory leak when monitoring remote hosts over long periods." },
        { type: "perf", title: "Process List Optimization", text: "Optimized process list rendering, reducing CPU usage by 15% on average." },
      ],
      codeSnippet: `{
  "theme": "custom",
  "custom_theme": {
    "primary": "#7D56F4",
    "secondary": "#EE6FF8",
    "success": "#A1E3AD",
    "background": "#1C1C1C"
  }
}`
    },
    {
      version: "v0.1.3",
      date: "February 15, 2026",
      title: "Cross-Platform Fixes",
      description: "A patch release focused on improving stability across different operating systems, particularly addressing disk space reporting on macOS and network interfaces on Linux.",
      type: "patch",
      changes: [
        { type: "fix", title: "macOS Disk Space", text: "Fixed an issue where disk space was reported incorrectly on macOS due to APFS volume mapping." },
        { type: "fix", title: "Linux Network Interfaces", text: "Corrected network interface selection logic on Linux to ignore virtual interfaces by default." },
        { type: "chore", title: "Dependency Updates", text: "Updated dependencies and improved build scripts for faster compilation." },
      ]
    },
    {
      version: "v0.1.2",
      date: "January 22, 2026",
      title: "Remote Monitoring",
      description: "You can now monitor multiple servers from a single bubbleMonitor instance. We've also added a new Services tab for managing systemd and launchd services directly from the TUI.",
      type: "minor",
      changes: [
        { type: "feat", title: "SSH Remote Hosts", text: "Added support for monitoring remote hosts via SSH. Configure hosts in your config.json." },
        { type: "feat", title: "Services Tab", text: "Introduced new 'Services' tab for systemd/launchd management (start, stop, restart)." },
        { type: "perf", title: "Network Module Rewrite", text: "Rewrote the network monitoring module for lower latency and less CPU overhead." },
      ],
      codeSnippet: `{
  "remote_hosts": [
    {
      "name": "Production DB",
      "host": "10.0.0.5:22",
      "user": "admin",
      "key_path": "~/.ssh/id_rsa"
    }
  ]
}`
    },
    {
      version: "v0.1.0",
      date: "December 10, 2025",
      title: "Initial Release",
      description: "The first public release of bubbleMonitor. A beautiful, terminal-based system monitor built with Go and BubbleTea.",
      type: "major",
      changes: [
        { type: "feat", title: "Core Modules", text: "Implemented CPU, Memory, Disk, Network, Temperatures, and Processes modules." },
        { type: "feat", title: "Cross-Platform", text: "Full support for Windows, macOS, and Linux out of the box." },
        { type: "feat", title: "Process Management", text: "Filter, sort, kill, and suspend processes seamlessly." },
      ]
    }
  ];

  const getTypeColor = (type: string) => {
    switch(type) {
      case 'feat': return 'text-emerald-400 border-emerald-400/30 bg-emerald-400/10';
      case 'fix': return 'text-amber-400 border-amber-400/30 bg-amber-400/10';
      case 'perf': return 'text-purple-400 border-purple-400/30 bg-purple-400/10';
      case 'chore': return 'text-slate-400 border-slate-400/30 bg-slate-400/10';
      default: return 'text-slate-400 border-slate-400/30 bg-slate-400/10';
    }
  };

  return (
    <div className="min-h-screen flex flex-col font-sans text-slate-300 selection:bg-slate-300 selection:text-zinc-950 relative overflow-x-hidden">
      <Navbar />

      <main className="flex-1 w-full mx-auto">
        
        {/* Hero Section */}
        <section className="p-8 md:p-16 lg:p-24 border-b border-dashed border-slate-700 bg-zinc-950/90 relative z-20">
          {/* Decorative Crosshairs */}
          <div className="absolute top-4 left-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute top-4 right-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute bottom-4 left-4 text-slate-700 font-mono text-xs">+</div>
          <div className="absolute bottom-4 right-4 text-slate-700 font-mono text-xs">+</div>

          <div className="max-w-5xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-400 mb-8 bg-zinc-900">
              <Terminal className="w-3 h-3 text-slate-300" />
              Release Notes
            </div>
            
            <h1 className="text-5xl md:text-7xl lg:text-8xl font-black tracking-tighter text-slate-100 mb-8 leading-[0.9] uppercase">
              System <br />
              <span className="text-slate-500">Updates &</span> <br />
              Logs.
            </h1>
            
            <p className="text-slate-400 text-lg md:text-xl max-w-2xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
              Keep track of all the new features, performance improvements, and bug fixes in bubbleMonitor.
            </p>
          </div>
        </section>

        {/* Changelog Entries */}
        <section className="flex flex-col bg-zinc-950/90">
          {releases.map((release, index) => (
            <div key={release.version} className="grid grid-cols-1 lg:grid-cols-12 border-b border-dashed border-slate-700">
              
              {/* Left Column: Version & Date */}
              <div className="lg:col-span-4 p-8 md:p-12 border-b lg:border-b-0 lg:border-r border-dashed border-slate-700 bg-zinc-900/30 flex flex-col justify-between relative">
                <div className="absolute top-4 right-4 text-[10px] font-mono text-slate-600">REL_{release.version.replace(/\./g, '_')}</div>
                
                <div>
                  <div className="flex items-center gap-3 mb-4">
                    <span className="text-xs font-mono text-slate-500 uppercase tracking-widest">[ VERSION ]</span>
                    {release.type === 'major' && <Star className="w-4 h-4 text-amber-400" />}
                  </div>
                  <h2 className="text-5xl md:text-6xl font-black tracking-tighter text-slate-100 uppercase leading-none">
                    {release.version}
                  </h2>
                </div>
                
                <div className="mt-12 flex flex-col gap-3 items-start">
                  <div className="inline-flex items-center gap-2 px-3 py-1.5 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-300 bg-zinc-950 shadow-[2px_2px_0px_0px_rgba(51,65,85,1)]">
                    {release.date}
                  </div>
                  <span className="text-[10px] font-mono uppercase tracking-widest text-slate-500 px-1">
                    TYPE: {release.type}
                  </span>
                </div>
              </div>

              {/* Right Column: Details */}
              <div className="lg:col-span-8 p-8 md:p-12 bg-zinc-950/50">
                <h3 className="text-2xl md:text-3xl font-bold text-slate-200 uppercase tracking-wide mb-6">
                  {release.title}
                </h3>
                <p className="text-slate-400 leading-relaxed mb-10 text-lg font-light">
                  {release.description}
                </p>
                
                <div className="space-y-4">
                  {release.changes.map((change, i) => (
                    <div key={i} className="p-5 border border-dashed border-slate-700 bg-zinc-900/40 flex flex-col sm:flex-row gap-4 sm:gap-6 hover:bg-zinc-900/80 transition-colors group">
                      <div className="shrink-0 pt-0.5">
                        <span className={`px-2 py-1 text-[10px] font-mono uppercase tracking-widest border border-dashed ${getTypeColor(change.type)}`}>
                          {change.type}
                        </span>
                      </div>
                      <div>
                        <h4 className="text-base font-bold text-slate-200 mb-2 group-hover:text-white transition-colors">{change.title}</h4>
                        <p className="text-sm text-slate-400 leading-relaxed">{change.text}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {release.codeSnippet && (
                  <div className="mt-8">
                    <CodeBlock 
                      code={release.codeSnippet} 
                      lang="json" 
                      filename="CONFIG.JSON"
                      wrapperClassName="relative border border-dashed border-slate-700 bg-zinc-950"
                      className="!p-6 !text-xs !leading-loose"
                      showCopy={true}
                    />
                  </div>
                )}
              </div>
              
            </div>
          ))}
        </section>

      </main>

      <Footer />
    </div>
  );
}
