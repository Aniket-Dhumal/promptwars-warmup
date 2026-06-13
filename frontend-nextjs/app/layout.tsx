import React from 'react';
import type { Metadata, Viewport } from 'next';
import { Orbitron, JetBrains_Mono } from 'next/font/google';
import './globals.css';

const orbitron = Orbitron({
  subsets: ['latin'],
  variable: '--font-orbitron',
  weight: ['400', '500', '700', '900'],
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-mono',
  weight: ['400', '700'],
});

export const metadata: Metadata = {
  title: 'Culinary Digital Twin | PromptWars Mumbai Domination Framework',
  description: 'Low-latency high-concurrency metabolic visualization cockpit and automated billing settlements engine.',
  icons: {
    icon: '/favicon.ico', // standard fallback path
  }
};

export const viewport: Viewport = {
  themeColor: '#030712',
  width: 'device-width',
  initialScale: 1,
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${orbitron.variable} ${jetbrainsMono.variable} scroll-smooth`}>
      <body className="antialiased selection:bg-cyber-orange/30 selection:text-white">
        {children}
      </body>
    </html>
  );
}
