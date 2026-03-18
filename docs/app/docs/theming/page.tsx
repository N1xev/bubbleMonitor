'use client';

import CodeBlock from '@/components/CodeBlock';

export default function ThemingPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Theming
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          Customize the look of bubbleMonitor with built-in themes or create your own.
        </p>
      </div>

      <div className="space-y-24">
        <div id="built-in-themes" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="built-in-themes" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Built-in Themes</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor comes with 32 built-in themes. To use one, simply set the <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">theme</code> option in your <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">config.json</code>:
          </p>

          <CodeBlock code={`{
  "theme": "nord"
}`} lang="json" filename="config.json" showCopy={true} />

          <h3 id="available-themes" className="text-xl font-bold tracking-tight text-slate-200 mt-12 mb-6">Available Themes</h3>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            {[
              'dark', 'light', 'nord', 'dracula', 'gruvbox', 'rosepine', 
              'everforest', 'nightowl', 'palenight', 'material', 'synthwave',
              'cobalt2', 'horizon', 'oceanic', 'palefire', 'github', 'moonlight',
              'shades', 'midnight', 'forest', 'autumn', 'cyberpunk', 'sunset',
              'ocean', 'coffee', 'solarized', 'monokai', 'catppuccin', 'tokyonight',
              'onedark', 'ayu', 'tty'
            ].map((theme) => (
              <div key={theme} className="bg-zinc-900/50 border border-dashed border-slate-700 px-4 py-2 rounded font-mono text-sm text-slate-300">
                {theme}
              </div>
            ))}
          </div>
        </div>

        <div id="custom-themes" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="custom-themes" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Custom Themes</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            If you want complete control over the colors, you can define a custom theme directly in your <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">config.json</code>. Set <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">theme</code> to <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">custom</code> and provide your color palette:
          </p>

          <CodeBlock code={`{
  "theme": "custom",
  "custom_theme": {
    "primary": "#3B82F6",
    "secondary": "#8B5CF6",
    "success": "#10B981",
    "warning": "#F59E0B",
    "alert": "#EF4444",
    "text": "#F9FAFB",
    "muted": "#9CA3AF",
    "border": "#4B5563",
    "background": "#111827"
  }
}`} lang="json" filename="config.json" showCopy={true} />

          <h3 id="custom-theme-options" className="text-xl font-bold tracking-tight text-slate-200 mt-12 mb-6">Custom Theme Options</h3>
          <div className="space-y-6">
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">primary <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Primary accent color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">secondary <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Secondary accent color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">success <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Success/healthy state color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">warning <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Warning state color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">alert <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Alert/critical state color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">text <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Primary text color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">muted <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Muted/secondary text color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">border <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Border color (hex code).</p>
            </div>
            <div className="border-l-2 border-dashed border-slate-700 pl-6">
              <h4 className="font-mono text-emerald-400 mb-2">background <span className="text-slate-500 text-sm">(string)</span></h4>
              <p className="text-slate-400 text-sm">Background color (hex code).</p>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
