import type { Metadata, Viewport } from 'next';
import './globals.css';
import { Providers } from '@/components/Providers';

export const viewport: Viewport = {
  themeColor: '#0f172a',
  width: 'device-width',
  initialScale: 1,
};

export const metadata: Metadata = {
  metadataBase: new URL('https://algoholic.dev'),
  title: {
    default: 'Algoholic - Master DSA for FAANG Interviews',
    template: '%s | Algoholic',
  },
  description: 'Practice platform for algorithm and data structure problems with AI-powered learning paths, spaced repetition, and personalized recommendations.',
  keywords: ['DSA', 'algorithms', 'data structures', 'coding interview', 'FAANG', 'leetcode', 'practice'],
  authors: [{ name: 'Algoholic Team' }],
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: 'https://algoholic.dev',
    siteName: 'Algoholic',
    title: 'Algoholic - Master DSA for FAANG Interviews',
    description: 'Practice platform for algorithm and data structure problems with AI-powered learning paths',
    images: [
      {
        url: '/logo.svg',
        width: 512,
        height: 512,
        alt: 'Algoholic Logo',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Algoholic - Master DSA for FAANG Interviews',
    description: 'Practice platform for algorithm and data structure problems with AI-powered learning paths',
    images: ['/logo.svg'],
  },
  icons: {
    icon: '/icon.svg',
    shortcut: '/icon.svg',
    apple: '/apple-touch-icon.svg',
  },
  manifest: '/manifest.json',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="/icon.svg" type="image/svg+xml" />
        <link rel="apple-touch-icon" href="/apple-touch-icon.svg" />
      </head>
      <body className="antialiased">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
