'use client';

import CodeBlock from '@/components/CodeBlock';

export default function DocsInstallationPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Installation
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          Get bubbleMonitor running on your system.
        </p>
      </div>

      <div className="space-y-24">
        <div id="installation" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="installing-bubblmonitor" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Installing bubbleMonitor</h2>
          </div>
          
          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor is distributed as a single static binary. Choose your preferred installation method below.
          </p>

          <div className="space-y-8">
            <div>
              <h3 id="go-install" className="text-lg font-bold tracking-tight text-slate-200 mb-2">Go Install (Recommended)</h3>
              <p className="text-slate-400 mb-4">If you have Go installed, this is the quickest way to get the latest version.</p>
              <CodeBlock code="go install github.com/N1xev/bubbleMonitor@latest" lang="bash" showCopy={true} />
            </div>

            <div>
              <h3 id="direct-download" className="text-lg font-bold tracking-tight text-slate-200 mb-2">Direct Download</h3>
              <p className="text-slate-400 mb-4">Download the latest release for your platform from GitHub.</p>
              
              <div className="space-y-4">
                <div>
                  <p className="text-slate-300 text-sm mb-2">Linux (x86_64):</p>
                  <CodeBlock code={`curl -L https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_linux_x86_64.tar.gz -o bub.tar.gz
tar -xzf bub.tar.gz
chmod +x bub
./bub`} lang="bash" showCopy={true} />
                </div>
                
                <div>
                  <p className="text-slate-300 text-sm mb-2">Linux (ARM64):</p>
                  <CodeBlock code={`curl -L https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_linux_aarch64.tar.gz -o bub.tar.gz
tar -xzf bub.tar.gz
chmod +x bub
./bub`} lang="bash" showCopy={true} />
                </div>
                
                <div>
                  <p className="text-slate-300 text-sm mb-2">macOS (Intel):</p>
                  <CodeBlock code={`curl -L https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_macos_x86_64.tar.gz -o bub.tar.gz
tar -xzf bub.tar.gz
chmod +x bub
./bub`} lang="bash" showCopy={true} />
                </div>
                
                <div>
                  <p className="text-slate-300 text-sm mb-2">macOS (Apple Silicon):</p>
                  <CodeBlock code={`curl -L https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_macos_aarch64.tar.gz -o bub.tar.gz
tar -xzf bub.tar.gz
chmod +x bub
./bub`} lang="bash" showCopy={true} />
                </div>
                
                <div>
                  <p className="text-slate-300 text-sm mb-2">Windows (x86_64):</p>
                  <CodeBlock code={`# PowerShell
Invoke-WebRequest -Uri https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_windows_x86_64.zip -OutFile bub.zip
Expand-Archive -Path bub.zip -DestinationPath .\\bub
.\\bub.exe`} lang="bash" showCopy={true} />
                </div>
                
                <div>
                  <p className="text-slate-300 text-sm mb-2">Windows (ARM64):</p>
                  <CodeBlock code={`# PowerShell
Invoke-WebRequest -Uri https://github.com/N1xev/bubbleMonitor/releases/latest/download/bub_windows_aarch64.zip -OutFile bub.zip
Expand-Archive -Path bub.zip -DestinationPath .\\bub
.\\bub.exe`} lang="bash" showCopy={true} />
                </div>
              </div>
            </div>
            
            <div>
              <h3 id="docker" className="text-lg font-bold tracking-tight text-slate-200 mb-2">Docker</h3>
              <p className="text-slate-400 mb-4">Run bubbleMonitor in a container with host network mode.</p>
              <CodeBlock code={`docker run -it --network host N1xev/bubblemonitor:latest`} lang="bash" showCopy={true} />
            </div>

            <div>
              <h3 id="build-from-source" className="text-lg font-bold tracking-tight text-slate-200 mb-2">Build from Source</h3>
              <p className="text-slate-400 mb-4">Clone the repository and build it yourself.</p>
              <CodeBlock code={`git clone https://github.com/N1xev/bubbleMonitor.git
cd bubbleMonitor
go build -o bub ./cmd/bub`} lang="bash" showCopy={true} />
            </div>
          </div>
        </div>

        <div id="gpu-support" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="gpu-support" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">GPU Support</h2>
          </div>
          
          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor supports NVIDIA and AMD GPUs. GPU libraries are not included by default to ensure the app runs on all systems without crashes.
          </p>

          <div className="space-y-6">
            <div>
              <h4 className="font-mono text-emerald-400 mb-2">With GPU Support (NVIDIA + AMD)</h4>
              <p className="text-slate-400 text-sm mb-2">Install with GPU library support:</p>
              <CodeBlock code="go install github.com/N1xev/bubbleMonitor@latest -tags with_gpus" lang="bash" showCopy={true} />
              <p className="text-slate-500 text-sm mt-2">Requires: NVIDIA Driver (libnvidia-ml.so) OR AMD ROCm</p>
            </div>
            
            <div>
              <h4 className="font-mono text-emerald-400 mb-2">Default Build (No GPU Libraries)</h4>
              <p className="text-slate-400 text-sm">The default build works on all systems. GPU detection uses slower CLI fallback (nvidia-smi).</p>
            </div>
          </div>
        </div>

        <div id="verification" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="verify-installation" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Verify Installation</h2>
          </div>
          
          <CodeBlock code="bub --version" lang="bash" showCopy={true} />
        </div>
      </div>
    </>
  );
}
