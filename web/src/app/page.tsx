import { api } from '@/lib/api';
import { DetailingService } from '@/lib/types';
import { ServiceCard } from '@/components/service-card';

export default async function HomePage() {
  let services: DetailingService[] = [];
  try {
    const data = await api<{ services: DetailingService[] }>('/catalogue');
    services = data.services ?? [];
  } catch { /* show empty */ }

  return (
    <div>
      {/* Hero */}
      <section className="relative min-h-[90vh] flex items-center justify-center px-4 text-center">
        <div className="max-w-3xl mx-auto">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 border border-border-subtle text-sm text-text-secondary mb-8">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-4 h-4 text-brand-400">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" />
            </svg>
            Perth, Western Australia
          </div>
          <h1 className="text-5xl md:text-7xl font-bold mb-6">
            <span className="text-brand-gradient">Premium</span> Mobile Car Detailing
          </h1>
          <p className="text-lg md:text-xl text-text-secondary mb-10 max-w-2xl mx-auto">
            We come to you with Bowden&apos;s Own premium Australian products. Wildlife and marine life safe, environmentally focused detailing.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <a href="/services" className="btn-brand text-lg px-8 py-4">
              Book Now
            </a>
            <a href="#about" className="inline-flex items-center justify-center px-8 py-4 rounded-full border border-border-subtle text-white font-semibold hover:bg-white/5 transition-colors text-lg">
              Learn More
            </a>
          </div>
        </div>
      </section>

      {/* About */}
      <section id="about" className="max-w-6xl mx-auto px-4 py-20">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Why <span className="text-brand-gradient">40 Degrees?</span>
          </h2>
          <p className="text-text-secondary max-w-2xl mx-auto">
            We bring professional-grade detailing to your doorstep. No need to drop your car off &mdash; we come to your home or workplace across Perth.
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="glass-card p-8 text-center">
            <div className="w-14 h-14 rounded-full bg-brand-500/20 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-7 h-7 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 0 0-2.455 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Premium Products</h3>
            <p className="text-text-secondary text-sm">
              Exclusively using Bowden&apos;s Own, Australia&apos;s finest car care range. Safe for wildlife and marine life.
            </p>
          </div>
          <div className="glass-card p-8 text-center">
            <div className="w-14 h-14 rounded-full bg-brand-500/20 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-7 h-7 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 18.75a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 0 1-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 0 0-3.213-9.193 2.056 2.056 0 0 0-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 0 0-10.026 0 1.106 1.106 0 0 0-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">We Come to You</h3>
            <p className="text-text-secondary text-sm">
              Mobile detailing across Perth. At your home, office, or wherever is convenient. No drop-off needed.
            </p>
          </div>
          <div className="glass-card p-8 text-center">
            <div className="w-14 h-14 rounded-full bg-brand-500/20 flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-7 h-7 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M12.75 3.03v.568c0 .334.148.65.405.864a4.5 4.5 0 0 1 0 6.836.878.878 0 0 0-.405.864v.568M12.75 3.03A7.502 7.502 0 0 0 5.25 10.5c0 1.47.318 2.874.896 4.13M12.75 3.03a7.44 7.44 0 0 1 3.346.816M3.998 17.336A5.001 5.001 0 0 0 9 21.75h5.25M8.25 21.75a5.002 5.002 0 0 0 5.002-4.414M12.75 3.03a7.504 7.504 0 0 1 3.346.816m0 0a7.445 7.445 0 0 1 2.158 1.922M8.25 21.75a5.002 5.002 0 0 1-3.252-4.414M16.096 3.846a7.445 7.445 0 0 1 2.158 1.922m0 0A7.458 7.458 0 0 1 19.5 10.5c0 1.472-.319 2.878-.9 4.14m0 0a5.002 5.002 0 0 1-3.35 4.414" />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Eco Friendly</h3>
            <p className="text-text-secondary text-sm">
              Environmentally conscious practices. All products are biodegradable, safe for waterways and wildlife.
            </p>
          </div>
        </div>
      </section>

      {/* Services */}
      {services.length > 0 && (
        <section className="max-w-6xl mx-auto px-4 py-20">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">Our Services</h2>
            <p className="text-text-secondary">From a quick wash to a full paint correction &mdash; we&apos;ve got you covered.</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {services.map(s => <ServiceCard key={s.id} service={s} />)}
          </div>
          <div className="text-center mt-10">
            <a href="/services" className="btn-brand">View All Services</a>
          </div>
        </section>
      )}

      {/* Contact */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="glass-card p-12 text-center">
          <h2 className="text-3xl font-bold mb-4">Get in Touch</h2>
          <p className="text-text-secondary mb-8 max-w-xl mx-auto">
            Ready to give your car the treatment it deserves? Book online or get in touch directly.
          </p>
          <div className="flex flex-col sm:flex-row gap-6 justify-center items-center text-text-secondary">
            <a href="tel:0400000000" className="flex items-center gap-2 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 0 0 2.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 0 1-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 0 0-1.091-.852H4.5A2.25 2.25 0 0 0 2.25 4.5v2.25Z" />
              </svg>
              0400 000 000
            </a>
            <a href="mailto:hello@40degrees.au" className="flex items-center gap-2 hover:text-white transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75" />
              </svg>
              hello@40degrees.au
            </a>
          </div>
        </div>
      </section>
    </div>
  );
}
