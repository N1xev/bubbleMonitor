import type {Metadata} from 'next';
import { Inter, JetBrains_Mono } from 'next/font/google';
import './globals.css'; // Global styles

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-sans',
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-mono',
});

export const metadata: Metadata = {
  title: 'bubbleMonitor - Terminal System Monitor',
  description: 'A beautiful terminal-based system monitor built with Go and BubbleTea.',
};

export default function RootLayout({children}: {children: React.ReactNode}) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable} dark`}>
      <body className="bg-zinc-950 text-zinc-300 font-sans antialiased selection:bg-zinc-800 selection:text-zinc-100" suppressHydrationWarning>
        {children}
      </body>
    </html>
  );
}
