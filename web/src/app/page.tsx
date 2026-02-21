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
      <section className="bg-gray-900 text-white py-20 px-4 text-center">
        <h1 className="text-4xl font-bold mb-4">40 Degrees Car Detailing</h1>
        <p className="text-lg text-gray-300 mb-8">Premium mobile detailing in Perth, Western Australia</p>
        <a href="/services" className="bg-white text-gray-900 px-6 py-3 rounded font-semibold hover:bg-gray-100">View Services</a>
      </section>
      <section className="max-w-6xl mx-auto py-12 px-4">
        <h2 className="text-2xl font-bold mb-6">Our Services</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {services.map(s => <ServiceCard key={s.id} service={s} />)}
        </div>
      </section>
    </div>
  );
}
