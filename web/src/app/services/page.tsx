import { api } from '@/lib/api';
import { PackageCard } from '@/components/package-card';
import type { ServiceCategory, DetailingService } from '@/lib/types';

interface Props {
  searchParams: Promise<{ category?: string }>;
}

export default async function ServicesPage({ searchParams }: Props) {
  const params = await searchParams;

  let categories: ServiceCategory[] = [];
  let allServices: DetailingService[] = [];

  try {
    const [categoriesRes, servicesRes] = await Promise.allSettled([
      api<{ categories: ServiceCategory[] }>('/catalogue/categories'),
      api<{ services: DetailingService[] }>('/catalogue'),
    ]);
    if (categoriesRes.status === 'fulfilled') categories = categoriesRes.value.categories ?? [];
    if (servicesRes.status === 'fulfilled') allServices = servicesRes.value.services ?? [];
  } catch {
    /* both failed, show empty */
  }

  const selectedCategory = params.category ?? null;
  const services = selectedCategory
    ? allServices.filter(s => {
        const cat = categories.find(c => c.slug === selectedCategory);
        return cat ? s.categoryId === cat.id : true;
      })
    : allServices;

  return (
    <div className="max-w-7xl mx-auto px-4 py-12">
      <div className="text-center mb-12">
        <h1 className="text-3xl md:text-4xl font-bold mb-4">Our Packages</h1>
        <p className="text-text-secondary text-lg max-w-2xl mx-auto">
          Choose the perfect package for your vehicle&apos;s needs
        </p>
      </div>

      {categories.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-10 justify-center">
          <a
            href="/services"
            className={`px-4 py-2 rounded-full text-sm font-medium border transition-colors ${
              !selectedCategory
                ? 'bg-brand-500 text-white border-brand-500'
                : 'border-border-subtle text-text-secondary hover:border-brand-400'
            }`}
          >
            All
          </a>
          {categories.map(cat => (
            <a
              key={cat.id}
              href={`/services?category=${cat.slug}`}
              className={`px-4 py-2 rounded-full text-sm font-medium border transition-colors ${
                selectedCategory === cat.slug
                  ? 'bg-brand-500 text-white border-brand-500'
                  : 'border-border-subtle text-text-secondary hover:border-brand-400'
              }`}
            >
              {cat.name}
            </a>
          ))}
        </div>
      )}

      {services.length === 0 ? (
        <p className="text-text-muted text-center">No services found.</p>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8">
          {services.map(service => (
            <PackageCard key={service.id} service={service} />
          ))}
        </div>
      )}

      <div className="mt-16 glass-card p-8">
        <h3 className="text-lg font-semibold text-white mb-4">Important Notes</h3>
        <ul className="space-y-3 text-sm text-text-secondary">
          <li className="flex items-start gap-2">
            <span className="text-brand-400 mt-0.5 shrink-0">&bull;</span>
            <span>Prices include travel up to ~50km, all products used and results</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-brand-400 mt-0.5 shrink-0">&bull;</span>
            <span>A travel charge of $30 may be added for &ldquo;Just Get It Clean&rdquo; package over 30km away</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-brand-400 mt-0.5 shrink-0">&bull;</span>
            <span>Duration is approximate and may vary based on vehicle condition</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-brand-400 mt-0.5 shrink-0">&bull;</span>
            <span>Work commences once walk-around and details have been agreed</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-brand-400 mt-0.5 shrink-0">&bull;</span>
            <span>Access to standard hose connection and electrical outlet required</span>
          </li>
        </ul>
      </div>
    </div>
  );
}
