'use client';

import CodeBlock from '@/components/CodeBlock';

export default function QuickStartPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Quick Start
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          Get up and running with bubbleMonitor in seconds.
        </p>
      </div>

      <div className="space-y-24">
        <div id="basic-usage" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="basic-usage" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Basic Usage</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            Once installed, simply run the <code className="bg-slate-800 text-slate-200 px-1.5 py-0.5 rounded-sm font-mono text-sm">bub</code> command in your terminal.
          </p>
          
          <CodeBlock code="bub" lang="bash" showCopy={true} />
          
          <p className="text-slate-400 leading-relaxed mt-6">
            On first run, bubbleMonitor creates a default config file at <code className="bg-slate-800 text-slate-200 px-1.5 py-0.5 rounded-sm font-mono text-sm">~/.config/bubblemonitor/config.json</code> with sensible defaults.
          </p>
        </div>

        <div id="navigation" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="navigation" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Navigation</h2>
          </div>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Key</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">tab</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">l</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">right</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Switch to next tab</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">shift+tab</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">h</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">left</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Switch to previous tab</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">1</kbd> - <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">9</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Jump directly to tab (1=Processes, 2=Metrics, etc.)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="processes" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="processes" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Process Management</h2>
          </div>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Key</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">j</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">↓</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Move selection down</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">k</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">↑</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Move selection up</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">g</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Go to top</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">G</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Go to bottom</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">pgup</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">ctrl+u</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Page up</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">pgdown</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">ctrl+d</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Page down</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">home</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Go to top</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">end</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Go to bottom</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">s</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Cycle sort by (cpu &rarr; mem &rarr; pid &rarr; name &rarr; cpu)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">S</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle sort direction (asc/desc)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">f</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Enter filter mode (type to filter)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">c</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Clear process filter</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">T</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle tree view (normal/tree)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">enter</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Expand/collapse tree node</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">o</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle open files inspector</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">b</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Bookmark process</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">K</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Kill selected process</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">z</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Suspend process</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">x</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Resume process</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">+</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Decrease process priority (nice -1)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">-</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Increase process priority (nice +1)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">n</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle CPU normalization (raw/normalized)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="global" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="global" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Global Shortcuts</h2>
          </div>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Key</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">q</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">ctrl+c</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Quit the application</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">p</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Pause/resume monitoring</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">r</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Refresh all metrics immediately</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">.</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle settings overlay</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">?</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Toggle help overlay</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">e</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Save snapshot to file</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">C</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Cycle chart type (line &rarr; bar &rarr; braille &rarr; line)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">H</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Cycle history length (60 &rarr; 300 &rarr; 900 &rarr; 3600 &rarr; 60)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">space</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Expand/collapse tree node (Processes tab)</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">esc</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Close current overlay / Cancel / Exit filter mode</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="system-tab" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="system-tab" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">System Tab</h2>
          </div>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Key</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">]</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">{"}"}</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Next scrollable block</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">[</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">{'{'}</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Previous scrollable block</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="settings" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="settings" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Settings Overlay</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            Press <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">.</kbd> to open the settings overlay. Use these keys to navigate and change settings:
          </p>

          <div className="overflow-x-auto border border-dashed border-slate-700">
            <table className="w-full text-left text-sm text-slate-400">
              <thead className="bg-zinc-900/50 text-xs uppercase font-mono tracking-widest text-slate-300 border-b border-dashed border-slate-700">
                <tr>
                  <th className="px-6 py-4">Key</th>
                  <th className="px-6 py-4 border-l border-dashed border-slate-700">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-dashed divide-slate-700">
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">esc</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">.</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Close settings and save</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">k</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">↑</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Previous setting</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">j</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">↓</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Next setting</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">+</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">=</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">l</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">right</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Increase value</td>
                </tr>
                <tr className="hover:bg-zinc-900/30 transition-colors">
                  <td className="px-6 py-4 font-mono"><kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">-</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">_</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">h</kbd> / <kbd className="bg-zinc-800 px-2 py-1 rounded border border-slate-700">left</kbd></td>
                  <td className="px-6 py-4 border-l border-dashed border-slate-700">Decrease value</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div id="environment-variables" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="environment-variables" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Environment Variables</h2>
          </div>
          <ul className="list-disc list-inside space-y-3 text-slate-400 marker:text-emerald-400">
            <li><code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5 rounded">BUBBLEMONITOR_CONFIG</code>: Path to config file (overrides default).</li>
            <li><code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5 rounded">DEBUG</code>: Set to <code className="text-slate-300">1</code> to enable debug mode.</li>
          </ul>
        </div>
      </div>
    </>
  );
}
