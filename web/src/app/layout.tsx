import type { Metadata } from 'next';
import './globals.css';
import { Nav } from '@/components/nav';
import { Footer } from '@/components/footer';
import { BokehBackground } from '@/components/bokeh-background';

export const metadata: Metadata = {
  title: '40 Degrees Car Detailing',
  description: 'Premium mobile car detailing in Perth, Western Australia',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen flex flex-col bg-surface text-white">
        <BokehBackground />
        <Nav />
        <main className="relative z-10 flex-1">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
