'use client';

import { motion } from 'motion/react';
import { 
  Terminal, 
  Cpu, 
  HardDrive, 
  Activity, 
  Thermometer, 
  List, 
  Palette, 
  Github, 
  Download,
  Command,
  ArrowRight,
  Monitor,
  Zap,
  ShieldAlert,
  Copy,
  Check
} from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { codeToHtml } from 'shiki';
import Navbar from '@/components/Navbar';
import Footer from '@/components/Footer';
import CodeBlock from '@/components/CodeBlock';

export default function LandingPage() {
  return (
    <div className="min-h-screen flex flex-col font-sans text-slate-300 selection:bg-slate-300 selection:text-zinc-950 relative overflow-x-hidden">
      
      {/* Top Navigation Bar */}
      <Navbar />

      <main className="flex-1 w-full mx-auto">
        
        {/* Hero Section - Split Layout */}
        <section className="grid grid-cols-1 lg:grid-cols-12">
          
          {/* Left: Typography & CTA */}
          <div className="lg:col-span-7 p-8 md:p-16 lg:p-24 border-b lg:border-b-0 lg:border-r border-dashed border-slate-700 flex flex-col justify-center relative bg-zinc-950/90 z-20">
            {/* Decorative Crosshairs */}
            <div className="absolute top-4 left-4 text-slate-700 font-mono text-xs">+</div>
            <div className="absolute top-4 right-4 text-slate-700 font-mono text-xs">+</div>
            <div className="absolute bottom-4 left-4 text-slate-700 font-mono text-xs">+</div>
            <div className="absolute bottom-4 right-4 text-slate-700 font-mono text-xs">+</div>

            <motion.div 
              initial={{ opacity: 0, x: -20 }} 
              animate={{ opacity: 1, x: 0 }} 
              transition={{ duration: 0.7, ease: "easeOut" }}
            >
              <div className="inline-flex items-center gap-2 px-3 py-1 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-400 mb-8 bg-zinc-900">
                <Zap className="w-3 h-3 text-slate-300" />
                High-Performance TUI
              </div>
              
              <h1 className="text-5xl md:text-7xl lg:text-8xl font-black tracking-tighter text-slate-100 mb-8 leading-[0.9] uppercase">
                Shows you <br />
                <span className="text-slate-500">only what</span> <br />
                you want <br />
                <span className="text-slate-500">to see.</span>
              </h1>
              
              <p className="text-slate-400 text-lg md:text-xl max-w-xl mb-12 font-light leading-relaxed border-l-2 border-dashed border-slate-700 pl-6">
                A beautiful, terminal-based system monitor built with Go and BubbleTea. 
                Track your system metrics in real-time with a minimal, highly customizable interface.
              </p>
              
              <div className="flex flex-col sm:flex-row items-stretch gap-6">
                <Link 
                  href="#download" 
                  className="px-8 py-4 bg-slate-200 text-zinc-950 text-sm font-mono uppercase tracking-widest hover:bg-white transition-all flex items-center justify-center gap-3 font-bold shadow-[6px_6px_0px_0px_rgba(51,65,85,1)] hover:translate-y-[2px] hover:translate-x-[2px] hover:shadow-[4px_4px_0px_0px_rgba(51,65,85,1)] active:translate-y-[6px] active:translate-x-[6px] active:shadow-none"
                >
                  <Download className="w-5 h-5" /> Install Now
                </Link>
                <Link 
                  href="https://github.com/N1xev/bubbleMonitor" 
                  target="_blank" 
                  className="px-8 py-4 bg-zinc-900 border border-dashed border-slate-600 text-slate-200 text-sm font-mono uppercase tracking-widest hover:bg-zinc-800 transition-all flex items-center justify-center gap-3 font-bold shadow-[6px_6px_0px_0px_rgba(51,65,85,1)] hover:translate-y-[2px] hover:translate-x-[2px] hover:shadow-[4px_4px_0px_0px_rgba(51,65,85,1)] active:translate-y-[6px] active:translate-x-[6px] active:shadow-none"
                >
                  <Github className="w-5 h-5" /> View Source
                </Link>
              </div>
            </motion.div>
          </div>

          {/* Right: Terminal Mockup & Stats */}
          <div className="lg:col-span-5 bg-zinc-900/50 relative flex flex-col z-10">
            <div className="flex-1 p-8 md:p-12 flex items-center justify-center relative">
              {/* Decorative Grid Lines */}
              <div className="absolute inset-0 bg-[linear-gradient(to_right,#334155_1px,transparent_1px),linear-gradient(to_bottom,#334155_1px,transparent_1px)] bg-[size:2rem_2rem] opacity-20" />
              
              <motion.div 
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.8, delay: 0.2, ease: "easeOut" }}
                className="w-full relative z-10"
              >
                {/* Mockup Container */}
                <div className="border border-dashed border-slate-600 bg-zinc-950 p-2 shadow-[20px_20px_0px_0px_rgba(30,41,59,0.5)] group cursor-crosshair">
                  <div className="border border-dashed border-slate-700 p-1 flex justify-between items-center mb-2 bg-zinc-900">
                    <div className="flex gap-2 px-2">
                      <div className="w-2 h-2 bg-slate-700" />
                      <div className="w-2 h-2 bg-slate-700" />
                      <div className="w-2 h-2 bg-slate-700" />
                    </div>
                    <span className="text-[10px] font-mono text-slate-500 uppercase tracking-widest px-2">tty1</span>
                  </div>
                  <div className="relative aspect-[16/10] w-full bg-zinc-950 border border-dashed border-slate-800 overflow-hidden">
                    <Image 
                      src="https://github.com/user-attachments/assets/8929a57d-5160-4ef8-9169-69e1e42af11f"
                      alt="bubbleMonitor Screenshot"
                      fill
                      className="object-cover opacity-80 grayscale group-hover:grayscale-0 group-hover:scale-105 transition-all duration-700"
                      unoptimized
                      referrerPolicy="no-referrer"
                    />
                    {/* Scanline Overlay */}
                    <div className="absolute inset-0 bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.25)_50%)] bg-[size:100%_4px] pointer-events-none opacity-50" />
                  </div>
                </div>
              </motion.div>
            </div>
            
            {/* Quick Stats Bar */}
            <div className="grid grid-cols-3 border-t border-dashed border-slate-700 bg-zinc-950/90">
              {[
                { label: "MEMORY", val: "< 15MB" },
                { label: "BINARY", val: "~ 5MB" },
                { label: "THEMES", val: "30+" }
              ].map((stat, i) => (
                <div key={i} className="p-4 border-r border-dashed border-slate-700 last:border-r-0 flex flex-col items-center justify-center text-center">
                  <span className="text-[10px] font-mono text-slate-500 uppercase tracking-widest mb-1">{stat.label}</span>
                  <span className="text-lg font-bold font-mono text-slate-200">{stat.val}</span>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Marquee Tape */}
        <div className="border-t border-b border-dashed border-slate-700 bg-slate-200 text-zinc-950 py-3 overflow-hidden flex whitespace-nowrap relative z-20 shadow-[0_0_40px_rgba(0,0,0,0.5)]">
          <div className="animate-[marquee_40s_linear_infinite] flex items-center gap-8 font-mono text-sm font-bold uppercase tracking-widest w-max">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="flex items-center gap-8">
                <span>CPU usage (per-core!)</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Memory consumption</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Disk space</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Network activity</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Battery status</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>System temperatures</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Filter, sort, kill, suspend</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>30+ gorgeous themes</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>Windows, Linux, macOS</span><span className="w-1.5 h-1.5 bg-zinc-950" />
                <span>bub</span><span className="w-1.5 h-1.5 bg-zinc-950" />
              </div>
            ))}
          </div>
        </div>

        {/* Features Bento Grid */}
        <section className="bg-zinc-950/90">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
            {/* Section Header */}
            <div className="p-8 md:p-12 border-b md:border-b-0 md:border-r border-dashed border-slate-700 flex flex-col justify-between bg-zinc-900/30">
              <div>
                <span className="text-xs font-mono text-slate-500 uppercase tracking-widest block mb-4">[ MODULES ]</span>
                <h2 className="text-4xl font-black tracking-tighter text-slate-100 uppercase leading-none mb-6">Core<br/>Features</h2>
                <p className="text-slate-400 font-light leading-relaxed">
                  Everything you need to keep your system running smoothly, packed into a single binary.
                </p>
              </div>
              <Monitor className="w-12 h-12 text-slate-700 mt-12" strokeWidth={1} />
            </div>

            {/* Feature Cards */}
            {[
              { id: "01", icon: Cpu, title: "CPU & Memory", desc: "Track usage per-core and memory consumption in real-time." },
              { id: "02", icon: HardDrive, title: "Disk Space", desc: "Keep an eye on storage across all mounted drives." },
              { id: "03", icon: Activity, title: "Network", desc: "Monitor upload and download speeds instantly." },
              { id: "04", icon: Thermometer, title: "Temperatures", desc: "Check system temperatures to prevent overheating." },
              { id: "05", icon: List, title: "Processes", desc: "Filter, sort, kill, or suspend processes seamlessly." },
            ].map((feature, i) => (
              <motion.div 
                key={i} 
                initial={{ opacity: 0 }}
                whileInView={{ opacity: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: i * 0.1 }}
                className="p-8 border-b border-r border-dashed border-slate-700 hover:bg-slate-900/50 transition-colors group relative flex flex-col justify-between min-h-[240px]"
              >
                <div className="absolute top-4 right-4 text-[10px] font-mono text-slate-600">MOD_{feature.id}</div>
                <div className="w-12 h-12 border border-dashed border-slate-600 flex items-center justify-center mb-8 group-hover:border-slate-400 group-hover:bg-slate-800 transition-all">
                  <feature.icon className="w-5 h-5 text-slate-300" strokeWidth={1.5} />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-slate-200 mb-2 uppercase tracking-wide">{feature.title}</h3>
                  <p className="text-sm text-slate-400 leading-relaxed">{feature.desc}</p>
                </div>
              </motion.div>
            ))}
          </div>
        </section>

        {/* Config & Download Split */}
        <section id="download" className="grid grid-cols-1 lg:grid-cols-2 border-b border-dashed border-slate-700 bg-zinc-950/90">
          
          {/* Left: Installation */}
          <div className="p-8 md:p-16 border-b lg:border-b-0 lg:border-r border-dashed border-slate-700">
            <div className="flex items-center gap-3 mb-12">
              <span className="w-3 h-3 bg-slate-200" />
              <h2 className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Installation</h2>
            </div>
            
            <div className="space-y-8">
              {[
                { os: "Linux", cmd: "curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.4/bub-linux-amd64-v0.1.4 -o bub\nchmod +x bub" },
                { os: "macOS (Apple Silicon)", cmd: "curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.4/bub-darwin-arm64-v0.1.4 -o bub\nchmod +x bub" },
                { os: "Windows", cmd: "curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.4/bub-windows-amd64-v0.1.4.exe -o bub.exe" },
              ].map((platform, i) => (
                <div key={i} className="border border-dashed border-slate-700 bg-zinc-950">
                  <div className="pl-4 pr-0 py-0 border-b border-dashed border-slate-700 flex justify-between items-center bg-zinc-900/50">
                    <span className="text-xs font-mono uppercase tracking-widest text-slate-400 py-2">{platform.os}</span>
                    <div className="px-3 py-2 bg-slate-800 text-[10px] font-mono text-slate-300 border-l border-dashed border-slate-700 uppercase tracking-widest">
                      CMD
                    </div>
                  </div>
                  <CodeBlock 
                    code={platform.cmd} 
                    lang="bash" 
                    showCopy={true} 
                    showHeader={false}
                    wrapperClassName="relative group"
                  />
                </div>
              ))}
            </div>
          </div>

          {/* Right: Configuration & Shortcuts */}
          <div className="p-8 md:p-16 bg-zinc-900/20">
            <div className="flex items-center gap-3 mb-12">
              <span className="w-3 h-3 border border-slate-200" />
              <h2 className="text-3xl font-black tracking-tighter text-slate-100 uppercase">Configuration</h2>
            </div>
            
            <p className="text-sm text-slate-400 mb-8 leading-relaxed font-mono">
              &gt; Config path: <code className="text-slate-200 bg-slate-800 px-2 py-1">~/.config/bubble-monitor/config.json</code>
            </p>
            
            <CodeBlock 
              code={`{
  "thresholds": {
    "CPU": 50,
    "Disk": 90,
    "Memory": 70,
    "Temperature": 60
  },
  "history_length": 300,
  "chart_type": "braille",
  "view_type": "normal",
  "sort_by": "cpu",
  "theme": "custom",
  "refresh_rate": 500,
  "border_type": "normal",
  "border_style": "dashed",
  "background_opaque": true,
  "process_cpu_normalized": true,
  "sort_direction": "asc",
  "tabs": [
    "Processes",
    "Metrics",
    "Disks",
    "Network",
    "Services",
    "Connections",
    "Logs",
    "Remote",
    "System"
  ],
  "default_tab": "Processes",
  "logging": {
    "enabled": false,
    "path": "Projects/Golang/bubbleMonitor/bubLog.txt"
  },
  "remote_hosts": [
    {
      "name": "SamLab 1",
      "host": "192.168.1.10:8080"
    }
  ],
  "health_weights": {
    "cpu_critical": 30,
    "cpu_high": 10,
    "mem_critical": 30,
    "mem_high": 10,
    "disk_critical": 20,
    "temp_critical": 30,
    "temp_high": 10
  },
  "custom_theme": {
    "primary": "#7D56F4",
    "secondary": "#EE6FF8",
    "success": "#A1E3AD",
    "warning": "#F5A962",
    "alert": "#F25D94",
    "text": "#F0F0F0",
    "muted": "#A0A0A0",
    "border": "#4A4A4A",
    "background": "#1C1C1C"
  }
}`}
              lang="json"
              filename="CONFIG.JSON"
              wrapperClassName="border border-dashed border-slate-700 bg-zinc-950 mb-16 relative group"
              className="h-[280px] overflow-y-auto !p-6 !text-xs !leading-loose"
              showCopy={true}
            />

            <div className="flex items-center gap-3 mb-8">
              <Command className="w-5 h-5 text-slate-400" />
              <h3 className="text-xl font-bold tracking-tight text-slate-200 uppercase">Keybinds</h3>
            </div>
            
            <div className="border border-dashed border-slate-700 bg-zinc-950">
              {[
                { key: "Tab / 1-6", desc: "Navigate tabs", fullDesc: "Switch between different monitoring tabs" },
                { key: "P", desc: "Pause/resume", fullDesc: "Pause or resume real-time monitoring" },
                { key: "S / f", desc: "Sort / Filter", fullDesc: "Sort or filter the processes list" },
                { key: "K", desc: "Kill process", fullDesc: "Kill the currently selected process" },
                { key: "z / x", desc: "Suspend/resume", fullDesc: "Suspend or resume the selected process" },
              ].map((shortcut, i) => (
                <div key={i} className="group relative flex items-center justify-between p-4 border-b border-dashed border-slate-700 last:border-0 text-sm cursor-help hover:bg-slate-900 transition-colors">
                  <span className="text-slate-400 font-mono uppercase tracking-wider text-xs">{shortcut.desc}</span>
                  <kbd className="font-mono text-xs text-slate-200 bg-slate-800 px-2 py-1 border border-slate-600">{shortcut.key}</kbd>
                  
                  {/* Tooltip */}
                  <div className="absolute bottom-full right-0 mb-2 w-max max-w-xs px-4 py-3 bg-slate-200 text-zinc-950 text-xs font-mono opacity-0 group-hover:opacity-100 transition-all duration-200 pointer-events-none z-10 shadow-[8px_8px_0px_0px_rgba(30,41,59,0.5)] translate-y-2 group-hover:translate-y-0 border border-zinc-950">
                    {shortcut.fullDesc}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Comparison Section */}
        <section className="p-8 md:p-16 lg:p-24 border-b border-dashed border-slate-700 bg-zinc-950 relative">
          <div className="max-w-5xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-400 mb-12 bg-zinc-900">
              <Monitor className="w-3 h-3 text-slate-300" />
              vs. The Alternatives
            </div>
            
            <h2 className="text-4xl md:text-5xl font-black tracking-tighter text-slate-100 uppercase mb-12">
              Why bubbleMonitor?
            </h2>

            <div className="overflow-x-auto border border-dashed border-slate-700 bg-zinc-900/30">
              <table className="w-full text-left font-mono text-sm border-collapse">
                <thead>
                  <tr className="border-b border-dashed border-slate-700 bg-zinc-900/80">
                    <th className="p-4 text-slate-500 font-normal uppercase tracking-widest border-r border-dashed border-slate-700">Feature</th>
                    <th className="p-4 text-emerald-400 font-bold uppercase tracking-widest bg-emerald-400/10 border-r border-dashed border-slate-700">bubbleMonitor</th>
                    <th className="p-4 text-slate-500 font-normal uppercase tracking-widest border-r border-dashed border-slate-700">htop</th>
                    <th className="p-4 text-slate-500 font-normal uppercase tracking-widest">btop</th>
                  </tr>
                </thead>
                <tbody className="text-slate-300">
                  <tr className="border-b border-dashed border-slate-700 hover:bg-zinc-900/50 transition-colors">
                    <td className="p-4 border-r border-dashed border-slate-700">Language</td>
                    <td className="p-4 font-bold text-slate-100 bg-emerald-400/5 border-r border-dashed border-slate-700">Go (BubbleTea)</td>
                    <td className="p-4 text-slate-500 border-r border-dashed border-slate-700">C</td>
                    <td className="p-4 text-slate-500">C++</td>
                  </tr>
                  <tr className="border-b border-dashed border-slate-700 hover:bg-zinc-900/50 transition-colors">
                    <td className="p-4 border-r border-dashed border-slate-700">Memory Footprint</td>
                    <td className="p-4 font-bold text-slate-100 bg-emerald-400/5 border-r border-dashed border-slate-700">~15MB</td>
                    <td className="p-4 text-slate-500 border-r border-dashed border-slate-700">~5MB</td>
                    <td className="p-4 text-slate-500">~25MB</td>
                  </tr>
                  <tr className="border-b border-dashed border-slate-700 hover:bg-zinc-900/50 transition-colors">
                    <td className="p-4 border-r border-dashed border-slate-700">Custom Theming (JSON)</td>
                    <td className="p-4 font-bold text-emerald-400 bg-emerald-400/5 border-r border-dashed border-slate-700">Yes</td>
                    <td className="p-4 text-slate-500 border-r border-dashed border-slate-700">No</td>
                    <td className="p-4 text-slate-500">Yes (Custom format)</td>
                  </tr>
                  <tr className="border-b border-dashed border-slate-700 hover:bg-zinc-900/50 transition-colors">
                    <td className="p-4 border-r border-dashed border-slate-700">Remote SSH Monitoring</td>
                    <td className="p-4 font-bold text-emerald-400 bg-emerald-400/5 border-r border-dashed border-slate-700">Yes (Built-in)</td>
                    <td className="p-4 text-slate-500 border-r border-dashed border-slate-700">No</td>
                    <td className="p-4 text-slate-500">No</td>
                  </tr>
                  <tr className="hover:bg-zinc-900/50 transition-colors">
                    <td className="p-4 border-r border-dashed border-slate-700">Braille Charts</td>
                    <td className="p-4 font-bold text-emerald-400 bg-emerald-400/5 border-r border-dashed border-slate-700">Yes</td>
                    <td className="p-4 text-slate-500 border-r border-dashed border-slate-700">No</td>
                    <td className="p-4 text-slate-500">Yes</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* FAQ Section */}
        <section className="p-8 md:p-16 lg:p-24 border-b border-dashed border-slate-700 bg-zinc-900/20 relative">
          <div className="max-w-4xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 border border-dashed border-slate-700 text-xs font-mono uppercase tracking-widest text-slate-400 mb-12 bg-zinc-900">
              <Check className="w-3 h-3 text-slate-300" />
              FAQ
            </div>
            
            <h2 className="text-4xl md:text-5xl font-black tracking-tighter text-slate-100 uppercase mb-12">
              Frequently Asked Questions
            </h2>

            <div className="space-y-6">
              {[
                {
                  q: "Does bubbleMonitor support Windows?",
                  a: "Yes! bubbleMonitor is fully cross-platform and runs on Windows, macOS, and Linux out of the box."
                },
                {
                  q: "How do I create a custom theme?",
                  a: "You can create a custom theme by adding a JSON file to your ~/.config/bubblemonitor/themes/ directory. Check the documentation for the exact schema."
                },
                {
                  q: "Can I monitor remote servers?",
                  a: "Absolutely. You can add remote hosts to your config.json file and bubbleMonitor will connect to them securely via SSH."
                },
                {
                  q: "Is it resource intensive?",
                  a: "Not at all. bubbleMonitor is written in Go and uses the BubbleTea framework, meaning it typically consumes less than 15MB of memory and minimal CPU."
                }
              ].map((faq, i) => (
                <div key={i} className="border border-dashed border-slate-700 bg-zinc-950 p-6 md:p-8 hover:bg-zinc-900/50 transition-colors">
                  <h3 className="text-xl font-bold text-slate-200 mb-4 flex items-start gap-3">
                    <span className="text-emerald-400 font-mono mt-1">Q.</span>
                    {faq.q}
                  </h3>
                  <p className="text-slate-400 leading-relaxed font-light pl-7">
                    {faq.a}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Testimonial / Log */}
        <section className="border-b border-dashed border-slate-700 bg-zinc-950/90 overflow-hidden relative py-16 md:py-24">
          <div className="absolute top-4 left-4 text-[10px] font-mono text-slate-600 z-10">LOG_ENTRIES</div>
          <div className="absolute bottom-4 right-4 text-[10px] font-mono text-slate-600 z-10">VERIFIED_USERS</div>
          
          <div className="flex whitespace-nowrap animate-[marquee_60s_linear_infinite] w-max hover:[animation-play-state:paused]">
            {[
              {
                quote: "bubbleMonitor is exactly what I needed. It shows me only what I want to see, without any of the clutter. The minimal TUI is an absolute joy to use.",
                name: "Alex Developer",
                role: "Software Engineer",
                seed: "user1"
              },
              {
                quote: "Finally, a system monitor that doesn't consume half my CPU just to render. The braille charts are a fantastic touch.",
                name: "Sarah Jenkins",
                role: "DevOps Specialist",
                seed: "user2"
              },
              {
                quote: "I've replaced htop and btop across all my servers. The custom theming means it perfectly matches my dotfiles.",
                name: "Marcus Chen",
                role: "Linux Enthusiast",
                seed: "user3"
              },
              {
                quote: "The ability to filter and kill processes directly from the TUI without dropping to a shell is a massive time saver.",
                name: "Elena Rodriguez",
                role: "Backend Dev",
                seed: "user4"
              }
            ].concat([
              {
                quote: "bubbleMonitor is exactly what I needed. It shows me only what I want to see, without any of the clutter. The minimal TUI is an absolute joy to use.",
                name: "Alex Developer",
                role: "Software Engineer",
                seed: "user1"
              },
              {
                quote: "Finally, a system monitor that doesn't consume half my CPU just to render. The braille charts are a fantastic touch.",
                name: "Sarah Jenkins",
                role: "DevOps Specialist",
                seed: "user2"
              },
              {
                quote: "I've replaced htop and btop across all my servers. The custom theming means it perfectly matches my dotfiles.",
                name: "Marcus Chen",
                role: "Linux Enthusiast",
                seed: "user3"
              },
              {
                quote: "The ability to filter and kill processes directly from the TUI without dropping to a shell is a massive time saver.",
                name: "Elena Rodriguez",
                role: "Backend Dev",
                seed: "user4"
              }
            ]).map((t, i, arr) => (
              <div key={i} className="w-[85vw] md:w-[700px] flex-shrink-0 flex flex-col items-center justify-center text-center px-8 md:px-16 whitespace-normal relative">
                {i !== arr.length - 1 && (
                  <div className="absolute right-0 top-1/2 -translate-y-1/2 w-px h-1/2 border-r border-dashed border-slate-700" />
                )}
                <ShieldAlert className="w-8 h-8 text-slate-600 mb-8" strokeWidth={1} />
                
                <blockquote className="text-xl md:text-3xl font-light text-slate-200 mb-12 leading-tight tracking-tight">
                  &quot;{t.quote}&quot;
                </blockquote>
                
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 border border-dashed border-slate-600 p-1 bg-zinc-900">
                    <div className="w-full h-full relative grayscale">
                      <Image
                        src={`https://picsum.photos/seed/${t.seed}/100/100`}
                        alt="User Avatar"
                        fill
                        className="object-cover"
                        unoptimized
                        referrerPolicy="no-referrer"
                      />
                    </div>
                  </div>
                  <div className="text-left flex flex-col">
                    <span className="text-sm font-bold text-slate-200 uppercase tracking-widest">{t.name}</span>
                    <span className="text-[10px] font-mono uppercase tracking-widest text-slate-500">{t.role}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
