import { api } from '@/lib/api';
import { formatPrice } from '@/lib/format';
import { AddToCart } from './add-to-cart';
import type { DetailingService } from '@/lib/types';

interface Props {
  params: Promise<{ slug: string }>;
}

export default async function ServiceDetailPage({ params }: Props) {
  const { slug } = await params;

  let service: DetailingService;
  try {
    const data = await api<{ service: DetailingService }>(`/catalogue/${slug}`);
    service = data.service;
  } catch {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        <a href="/services" className="text-sm text-text-muted hover:text-white mb-4 inline-block">
          &larr; Back to services
        </a>
        <p className="text-text-muted mt-4">Service not found or unavailable.</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <a href="/services" className="text-sm text-text-muted hover:text-white mb-4 inline-block">
        &larr; Back to services
      </a>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        <div>
          <h1 className="text-3xl font-bold mb-2">{service.name}</h1>
          <p className="text-sm text-text-muted mb-4">{service.categoryName}</p>
          <p className="text-text-secondary leading-relaxed mb-6">{service.description}</p>

          <div className="flex items-center gap-6 text-sm">
            <div>
              <span className="text-text-muted">Starting from</span>
              <p className="text-2xl font-bold text-brand-400">{formatPrice(service.basePrice)}</p>
            </div>
            <div>
              <span className="text-text-muted">Duration</span>
              <p className="text-lg font-semibold">{service.durationMinutes} mins</p>
            </div>
          </div>
        </div>

        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold mb-4">Add to Cart</h2>
          <AddToCart service={service} />
        </div>
      </div>
    </div>
  );
}
