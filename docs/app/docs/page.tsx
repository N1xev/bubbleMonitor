'use client';

import CodeBlock from '@/components/CodeBlock';
import Link from 'next/link';

export default function IntroductionPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Introduction
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          A beautiful terminal-based system monitor built with Go and BubbleTea.
        </p>
      </div>

      <div className="space-y-24">
        <div id="what-is" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="what-is-bubblemonitor" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">What is bubbleMonitor?</h2>
          </div>
          
          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor is a <strong className="text-slate-200">real-time terminal-based system monitor</strong> that shows you exactly what you want to see. Built with Go and BubbleTea, it provides a slick TUI (Terminal User Interface) for tracking your system metrics without the overhead of GUI applications.
          </p>

          <blockquote className="border-l-4 border-emerald-400 bg-emerald-400/10 p-4 my-6 text-emerald-200/80 font-mono text-sm">
            <strong>Motto:</strong> &quot;shows you only what you want to see!&quot;
          </blockquote>
        </div>

        <div id="features" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="features" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Features</h2>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <h3 id="system-monitoring" className="text-lg font-bold text-slate-200">System Monitoring</h3>
              <ul className="space-y-2 text-slate-400">
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>CPU usage (per-core detailed breakdown)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Memory consumption and usage</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Disk space and I/O statistics</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Network activity and bandwidth</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Battery status and charge level</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>System temperatures and hardware info</span>
                </li>
              </ul>
            </div>

            <div className="space-y-4">
              <h3 id="process-management" className="text-lg font-bold text-slate-200">Process Management</h3>
              <ul className="space-y-2 text-slate-400">
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>View all running processes</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Filter and search processes</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Sort by CPU, memory, PID, or name</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Kill, suspend, and resume processes</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>View open files for each process</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Tree view for process hierarchy</span>
                </li>
              </ul>
            </div>

            <div className="space-y-4">
              <h3 id="customization" className="text-lg font-bold text-slate-200">Customization</h3>
              <ul className="space-y-2 text-slate-400">
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>32 built-in themes (Dracula, Nord, Gruvbox, Tokyo Night, and more)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Custom color palette support</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Multiple chart styles (bar, line, area, sparkline, braille)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Customizable border styles</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Configurable refresh rates</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Alert thresholds for CPU, memory, disk, temperature</span>
                </li>
              </ul>
            </div>

            <div className="space-y-4">
              <h3 id="advanced-features" className="text-lg font-bold text-slate-200">Advanced Features</h3>
              <ul className="space-y-2 text-slate-400">
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>GPU monitoring (NVIDIA and AMD)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Remote host monitoring via SSH</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Service monitoring (systemd, Windows services)</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Network connections and port listing</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Log file monitoring</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-emerald-400 mt-1">▸</span>
                  <span>Alerts for threshold violations</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <div id="platforms" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="cross-platform" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Cross-Platform</h2>
          </div>

          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor works beautifully on Windows, Linux, and macOS. Most features are available across all platforms, with some platform-specific differences.
          </p>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Feature</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Linux</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">macOS</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Windows</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4">CPU, Memory, Disk, Network</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Full</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Full</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Full</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4">GPU Monitoring</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">NVIDIA + AMD</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-slate-500">Limited</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">NVIDIA</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4">Load Averages</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Full</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-slate-500">N/A</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-slate-500">N/A</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4">Temperature</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Full</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-slate-500">Limited</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-amber-400">Admin required</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4">Services</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">systemd</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-slate-500">LaunchD</td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700 text-emerald-400">Windows Services</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="quick-example" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="quick-example" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Quick Example</h2>
          </div>

          <p className="text-slate-400 leading-relaxed mb-6">
            Once installed, simply run the <code className="bg-slate-800 text-slate-200 px-1.5 py-0.5 rounded-sm font-mono text-sm">bub</code> command in your terminal:
          </p>

          <CodeBlock code="bub" lang="bash" showCopy={true} />

          <p className="text-slate-400 leading-relaxed mt-6">
            On first run, bubbleMonitor creates a default config file at <code className="bg-slate-800 text-slate-200 px-1.5 py-0.5 rounded-sm font-mono text-sm">~/.config/bubblemonitor/config.json</code> with sensible defaults.
          </p>
        </div>

        <div id="next-steps" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="next-steps" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Next Steps</h2>
          </div>

          <p className="text-slate-400 leading-relaxed mb-6">
            Ready to get started? Here&apos;s what to do next:
          </p>

          <div className="grid md:grid-cols-3 gap-4">
            <Link href="/docs/installation" className="group p-6 border border-dashed border-slate-700 bg-zinc-900/30 hover:bg-zinc-900 transition-colors">
              <h3 className="text-lg font-bold text-slate-200 group-hover:text-emerald-400 mb-2">Install</h3>
              <p className="text-slate-400 text-sm">Get bubbleMonitor running on your system</p>
            </Link>

            <Link href="/docs/quick-start" className="group p-6 border border-dashed border-slate-700 bg-zinc-900/30 hover:bg-zinc-900 transition-colors">
              <h3 className="text-lg font-bold text-slate-200 group-hover:text-emerald-400 mb-2">Quick Start</h3>
              <p className="text-slate-400 text-sm">Learn the basics and keyboard shortcuts</p>
            </Link>

            <Link href="/docs/configuration" className="group p-6 border border-dashed border-slate-700 bg-zinc-900/30 hover:bg-zinc-900 transition-colors">
              <h3 className="text-lg font-bold text-slate-200 group-hover:text-emerald-400 mb-2">Configure</h3>
              <p className="text-slate-400 text-sm">Customize themes, charts, and alerts</p>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
