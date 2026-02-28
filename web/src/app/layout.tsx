import type { Metadata } from 'next';
import './globals.css';
import { Nav } from '@/components/nav';
import { Footer } from '@/components/footer';
import { BokehBackground } from '@/components/bokeh-background';
import { ThemeProvider } from '@/components/theme-provider';

export const metadata: Metadata = {
  title: '40 Degrees Car Detailing',
  description: 'Premium mobile car detailing in Perth, Western Australia',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `try{var t=localStorage.getItem('theme')||(window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light');if(t==='dark')document.documentElement.classList.add('dark')}catch(e){}`,
          }}
        />
      </head>
      <body className="min-h-screen flex flex-col bg-surface text-foreground">
        <ThemeProvider>
          <BokehBackground />
          <Nav />
          <main className="relative z-10 flex-1">{children}</main>
          <Footer />
        </ThemeProvider>
      </body>
    </html>
  );
}
