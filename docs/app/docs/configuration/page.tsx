'use client';

import CodeBlock from '@/components/CodeBlock';

export default function ConfigurationPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Configuration
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          Customize bubbleMonitor to fit your workflow.
        </p>
      </div>

      <div className="space-y-24">
        <div id="config-file" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="the-config-json" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">The config.json</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            By default, bubbleMonitor looks for a configuration file at <code className="bg-slate-800 text-slate-200 px-1.5 py-0.5 rounded-sm font-mono text-sm">~/.config/bubblemonitor/config.json</code>. If it doesn&apos;t exist, a default one will be created on the first run.
          </p>
          
          <h3 id="example-configuration" className="text-xl font-bold tracking-tight text-slate-200 mt-12 mb-6">Example Configuration</h3>
          <CodeBlock 
            code={`{
  "tabs": ["Processes", "Metrics", "Disks", "Network", "System", "Services", "Connections", "Logs", "Remote"],
  "default_tab": "Processes",
  "refresh_rate": 1000,
  "chart_type": "braille",
  "view_type": "normal",
  "sort_by": "cpu",
  "sort_direction": "asc",
  "theme": "dark",
  "border_type": "rounded",
  "border_style": "dashed",
  "background_opaque": true,
  "process_cpu_normalized": true,
  "history_length": 60,
  "thresholds": {
    "cpu": 90,
    "mem": 90,
    "disk": 90,
    "temp": 85
  },
  "logging": {
    "enabled": false,
    "path": ""
  },
  "remote_hosts": []
}`} 
            lang="json" 
            filename="config.json"
            showCopy={true} 
          />

          <h3 id="configuration-options" className="text-xl font-bold tracking-tight text-slate-200 mt-12 mb-6">Configuration Options</h3>
          
          <div className="space-y-6">
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">tabs <span className="text-slate-500 text-sm">(string[])</span></h4>
              <p className="text-slate-400 text-sm">Which tabs to display. Available: <code className="text-slate-300">Processes</code>, <code className="text-slate-300">Metrics</code>, <code className="text-slate-300">Disks</code>, <code className="text-slate-300">Network</code>, <code className="text-slate-300">System</code>, <code className="text-slate-300">Services</code>, <code className="text-slate-300">Connections</code>, <code className="text-slate-300">Logs</code>, <code className="text-slate-300">Remote</code>.</p>
            </div>
            
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">default_tab <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">The tab focused on startup. Default: <code className="text-slate-300">Processes</code>.</p>
            </div>
            
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">refresh_rate <span className="text-slate-500 text-sm">(number)</span></h4>
              <p className="text-slate-400 text-sm">How often the UI updates in milliseconds. Options: <code className="text-slate-300">500</code>, <code className="text-slate-300">1000</code>, <code className="text-slate-300">2000</code>, <code className="text-slate-300">5000</code>. Default: <code className="text-slate-300">1000</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">chart_type <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Type of chart to display. Options: <code className="text-slate-300">line</code>, <code className="text-slate-300">bar</code>, <code className="text-slate-300">braille</code>. Default: <code className="text-slate-300">braille</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">view_type <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Process view mode. Options: <code className="text-slate-300">normal</code>, <code className="text-slate-300">tree</code>. Default: <code className="text-slate-300">normal</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">sort_by <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">How to sort processes. Options: <code className="text-slate-300">cpu</code>, <code className="text-slate-300">mem</code>, <code className="text-slate-300">pid</code>, <code className="text-slate-300">name</code>. Default: <code className="text-slate-300">cpu</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">sort_direction <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Sort direction. Options: <code className="text-slate-300">asc</code>, <code className="text-slate-300">desc</code>. Default: <code className="text-slate-300">asc</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">theme <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Sets the color scheme. Built-in options: <code className="text-slate-300">dark</code>, <code className="text-slate-300">light</code>, <code className="text-slate-300">nord</code>, <code className="text-slate-300">dracula</code>, <code className="text-slate-300">gruvbox</code>, <code className="text-slate-300">solarized</code>, <code className="text-slate-300">monokai</code>, <code className="text-slate-300">catppuccin</code>, <code className="text-slate-300">tokyonight</code>, <code className="text-slate-300">onedark</code>, <code className="text-slate-300">ayu</code>, <code className="text-slate-300">rosepine</code>, <code className="text-slate-300">everforest</code>, <code className="text-slate-300">nightowl</code>, <code className="text-slate-300">palenight</code>, <code className="text-slate-300">material</code>, <code className="text-slate-300">synthwave</code>, <code className="text-slate-300">cobalt2</code>, <code className="text-slate-300">horizon</code>, <code className="text-slate-300">oceanic</code>, <code className="text-slate-300">palefire</code>, <code className="text-slate-300">github</code>, <code className="text-slate-300">moonlight</code>, <code className="text-slate-300">shades</code>, <code className="text-slate-300">midnight</code>, <code className="text-slate-300">forest</code>, <code className="text-slate-300">autumn</code>, <code className="text-slate-300">cyberpunk</code>, <code className="text-slate-300">sunset</code>, <code className="text-slate-300">ocean</code>, <code className="text-slate-300">coffee</code>, <code className="text-slate-300">custom</code>. Default: <code className="text-slate-300">dark</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">border_type <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Border style. Options: <code className="text-slate-300">normal</code>, <code className="text-slate-300">rounded</code>. Default: <code className="text-slate-300">rounded</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">border_style <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Border line style. Options: <code className="text-slate-300">single</code>, <code className="text-slate-300">double</code>, <code className="text-slate-300">dashed</code>. Default: <code className="text-slate-300">dashed</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">background_opaque <span className="text-slate-500 text-sm">(bool)</span></h4>
              <p className="text-slate-400 text-sm">Whether the background is opaque or transparent. Default: <code className="text-slate-300">true</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">process_cpu_normalized <span className="text-slate-500 text-sm">(bool)</span></h4>
              <p className="text-slate-400 text-sm">Normalize CPU usage across all cores (divide by number of CPUs). Default: <code className="text-slate-300">true</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">history_length <span className="text-slate-500 text-sm">(number)</span></h4>
              <p className="text-slate-400 text-sm">Number of data points to keep in history for charts. Default: <code className="text-slate-300">60</code>.</p>
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">thresholds <span className="text-slate-500 text-sm">(object)</span></h4>
              <p className="text-slate-400 text-sm mb-2">Health threshold settings for alerts. Values are percentages (0-100).</p>
              <CodeBlock code={`{
  "cpu": 90,
  "mem": 90,
  "disk": 90,
  "temp": 85
}`} lang="json" showCopy={true} />
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">logging <span className="text-slate-500 text-sm">(object)</span></h4>
              <p className="text-slate-400 text-sm mb-2">Logging configuration.</p>
              <CodeBlock code={`{
  "enabled": false,
  "path": ""
}`} lang="json" showCopy={true} />
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">remote_hosts <span className="text-slate-500 text-sm">(object[])</span></h4>
              <p className="text-slate-400 text-sm mb-2">SSH remote hosts for monitoring. See Remote Hosts documentation.</p>
              <CodeBlock code={`{
  "remote_hosts": [
    {
      "name": "server-1",
      "host": "user@192.168.1.100:22"
    }
  ]
}`} lang="json" showCopy={true} />
            </div>

            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">custom_theme <span className="text-slate-500 text-sm">(object)</span></h4>
              <p className="text-slate-400 text-sm mb-2">Define a custom theme inline. See Custom Theming documentation.</p>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
