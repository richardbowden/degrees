import { api } from '@/lib/api';
import { formatPrice } from '@/lib/format';
import type { DetailingService } from '@/lib/types';

export default async function HomePage() {
  let services: DetailingService[] = [];
  try {
    const data = await api<{ services: DetailingService[] }>('/catalogue');
    services = data.services ?? [];
  } catch { /* show empty */ }

  const lowestPrice = services.length > 0
    ? Math.min(...services.map(s => Number(s.basePrice)))
    : 0;
  const categories = [...new Set(services.map(s => s.categoryName).filter(Boolean))];

  return (
    <div>
      {/* Hero */}
      <section className="relative min-h-[90vh] flex items-center justify-center px-4 text-center">
        <div className="max-w-3xl mx-auto">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-surface-input border border-border-subtle text-sm text-text-secondary mb-8">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-4 h-4 text-brand-400">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" />
            </svg>
            Yanchep, down to Perth and beyond
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
            <a href="#about" className="inline-flex items-center justify-center px-8 py-4 rounded-full border border-border-subtle text-foreground font-semibold hover:bg-surface-hover transition-colors text-lg">
              Learn More
            </a>
          </div>
        </div>
      </section>

      {/* About */}
      <section id="about" className="max-w-6xl mx-auto px-4 py-20">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Why <span className="text-brand-gradient">40 Degrees?</span>
          </h2>
        </div>

        {/* The feeling */}
        <div className="glass-card p-8 md:p-12 mb-8">
          <p className="text-xl md:text-2xl text-foreground leading-relaxed mb-6 font-light">
            There&apos;s nothing quite like getting into your vehicle after a proper detailing session. <span className="text-brand-400 font-normal">That new car sensation, the pride you feel, the shine that makes you smile.</span>
          </p>
          <p className="text-text-secondary leading-relaxed text-lg mb-6">
            Whether it&apos;s your daily companion for school runs, your trusted partner for weekend adventures, or your pride and joy for off-road exploration &mdash; we treat it like it&apos;s our own. We love the transformation: removing the grime, restoring the shine, and delivering results that make you fall in love with your car all over again.
          </p>
        </div>

        {/* Products */}
        <div className="glass-card p-8 md:p-12">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-full bg-brand-500/20 flex items-center justify-center shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 0 0-2.455 2.456Z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold">Carefully Selected Premium Products</h3>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <p className="text-brand-400 font-medium mb-2">Bowden&apos;s Own</p>
              <p className="text-text-secondary leading-relaxed">
                Our Aussie-crafted range handles the complete end-to-end wash and finish. Safe for wildlife and marine life, delivering a flawless result every time.
              </p>
            </div>
            <div>
              <p className="text-brand-400 font-medium mb-2">Koch Chemie</p>
              <p className="text-text-secondary leading-relaxed">
                Tackles the deep cleaning work: preserving ceramic coatings, removing contamination, and reducing water spots over time.
              </p>
            </div>
          </div>
          <p className="text-text-muted text-sm mt-6 pt-6 border-t border-border-subtle">
            Both are biodegradable and exceed strict environmental standards &mdash; because caring for your car shouldn&apos;t come at the cost of caring for our planet.
          </p>
        </div>
      </section>

      {/* Services teaser */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="glass-card p-8 md:p-12">
          <div className="max-w-2xl mx-auto text-center">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              Ready to <span className="text-brand-gradient">get started?</span>
            </h2>
            <p className="text-text-secondary text-lg mb-8">
              From a quick exterior wash to a full paint correction &mdash; browse our packages and pick what suits your car.
            </p>

            {categories.length > 0 && (
              <div className="flex flex-wrap gap-2 justify-center mb-8">
                {categories.map(cat => (
                  <span key={cat} className="px-4 py-1.5 rounded-full bg-surface-input border border-border-subtle text-sm text-text-secondary">
                    {cat}
                  </span>
                ))}
              </div>
            )}

            {lowestPrice > 0 && (
              <p className="text-text-muted mb-8">
                {services.length} packages available &middot; Starting from <span className="text-brand-400 font-semibold">{formatPrice(lowestPrice)}</span>
              </p>
            )}

            <a href="/services" className="btn-brand text-lg px-10 py-4 inline-block">
              View Our Services
            </a>
          </div>
        </div>
      </section>

      {/* Contact */}
      <section className="max-w-6xl mx-auto px-4 py-20">
        <div className="glass-card p-12 text-center">
          <h2 className="text-3xl font-bold mb-4">Get in Touch</h2>
          <p className="text-text-secondary mb-8 max-w-xl mx-auto">
            Ready to give your car the treatment it deserves? Book online or get in touch directly.
          </p>
          <div className="flex flex-col sm:flex-row gap-6 justify-center items-center text-text-secondary">
            <a href="tel:0448263659" className="flex items-center gap-2 hover:text-foreground transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 0 0 2.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 0 1-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 0 0-1.091-.852H4.5A2.25 2.25 0 0 0 2.25 4.5v2.25Z" />
              </svg>
              0448 263 659
            </a>
            <a href="mailto:detailing@40degrees.au" className="flex items-center gap-2 hover:text-foreground transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5 text-brand-400">
                <path strokeLinecap="round" strokeLinejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91a2.25 2.25 0 0 1-1.07-1.916V6.75" />
              </svg>
              detailing@40degrees.au
            </a>
          </div>
        </div>
      </section>
    </div>
  );
}
