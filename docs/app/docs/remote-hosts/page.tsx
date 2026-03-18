'use client';

import CodeBlock from '@/components/CodeBlock';

export default function RemoteHostsPage() {
  return (
    <>
      <div className="mb-16 border-b border-dashed border-slate-700 pb-12">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">
          Remote Hosts (SSH)
        </h1>
        <p className="text-slate-400 text-lg md:text-xl font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
          Monitor your servers securely over SSH.
        </p>
      </div>

      <div className="space-y-24">
        <div id="ssh-setup" className="doc-section scroll-mt-32">
          <div className="flex items-center gap-3 mb-8 border-b border-dashed border-slate-700 pb-6">
            <span className="w-3 h-3 bg-slate-200" />
            <h2 id="ssh-setup" className="text-3xl font-black tracking-tighter text-slate-100 uppercase">SSH Setup</h2>
          </div>
          <p className="text-slate-400 leading-relaxed mb-6">
            bubbleMonitor can connect to remote servers via SSH to gather metrics without installing any agents on the target machine. It uses your existing SSH configuration and keys.
          </p>
          
          <h3 id="configuring-hosts" className="text-xl font-bold tracking-tight text-slate-200 mt-12 mb-6">Configuring Hosts</h3>
          <p className="text-slate-400 leading-relaxed mb-6">
            Add remote hosts to your <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">config.json</code> using the <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">remote_hosts</code> array. Each host requires a <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">name</code> and a <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">host</code> in the format <code className="font-mono text-emerald-400 bg-emerald-400/10 px-1.5 py-0.5">user@hostname:port</code>:
          </p>

          <CodeBlock code={`{
  "remote_hosts": [
    {
      "name": "web-prod-01",
      "host": "admin@192.168.1.10:22"
    },
    {
      "name": "db-prod-01", 
      "host": "admin@192.168.1.20"
    }
  ]
}`} lang="json" filename="config.json" showCopy={true} />

          <blockquote className="border-l-4 border-emerald-400 bg-emerald-400/10 p-4 my-6 text-emerald-200/80 font-mono text-sm">
            <strong>Security Tip:</strong> We strongly recommend using SSH keys with ssh-agent rather than hardcoding passwords in your configuration. Make sure your SSH key is added to the ssh-agent and the remote server has your public key in <code className="text-emerald-300">authorized_keys</code>.
          </blockquote>
        </div>
      </div>
    </>
  );
}
