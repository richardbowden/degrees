import { api } from '@/lib/api';
import { ServiceCard } from '@/components/service-card';
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
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-2">Our Services</h1>
      <p className="text-gray-600 mb-8">
        Premium mobile car detailing across Perth, using Bowden&apos;s Own products.
      </p>

      {categories.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-8">
          <a
            href="/services"
            className={`px-4 py-2 rounded-full text-sm font-medium border transition-colors ${
              !selectedCategory
                ? 'bg-gray-900 text-white border-gray-900'
                : 'border-gray-300 text-gray-700 hover:border-gray-500'
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
                  ? 'bg-gray-900 text-white border-gray-900'
                  : 'border-gray-300 text-gray-700 hover:border-gray-500'
              }`}
            >
              {cat.name}
            </a>
          ))}
        </div>
      )}

      {services.length === 0 ? (
        <p className="text-gray-500">No services found.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {services.map(service => (
            <ServiceCard key={service.id} service={service} />
          ))}
        </div>
      )}
    </div>
  );
}
