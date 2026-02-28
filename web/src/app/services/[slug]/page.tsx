import { api } from '@/lib/api';
import { formatPrice } from '@/lib/format';
import { MarkdownContent } from '@/components/markdown-content';
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
        <a href="/services" className="text-sm text-text-muted hover:text-foreground mb-4 inline-block">
          &larr; Back to services
        </a>
        <p className="text-text-muted mt-4">Service not found or unavailable.</p>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <a href="/services" className="text-sm text-text-muted hover:text-foreground mb-4 inline-block">
        &larr; Back to services
      </a>

      <div className="grid grid-cols-1 lg:grid-cols-5 gap-12">
        <div className="lg:col-span-3">
          <h1 className="text-3xl font-bold mb-2">{service.name}</h1>
          <p className="text-sm text-text-muted mb-6">{service.categoryName}</p>

          <MarkdownContent content={service.description} />

          {service.priceTiers && service.priceTiers.length > 0 && (
            <div className="mt-8 rounded-lg bg-brand-500/5 border border-brand-500/20 overflow-hidden">
              <div className="px-4 py-2 border-b border-brand-500/10">
                <p className="text-xs font-medium text-text-muted">Pricing by Vehicle Size</p>
              </div>
              {service.priceTiers.map(tier => (
                <div
                  key={tier.vehicleCategoryId}
                  className="flex items-center justify-between px-4 py-3 border-t border-brand-500/10 first:border-t-0"
                >
                  <span className="text-sm text-text-secondary">{tier.categoryName}</span>
                  <span className="text-lg font-bold text-brand-400">{formatPrice(tier.price)}</span>
                </div>
              ))}
            </div>
          )}

          <div className="flex items-center gap-6 text-sm mt-8 pt-6 border-t border-border-subtle">
            <div>
              <span className="text-text-muted">{service.priceTiers && service.priceTiers.length > 0 ? 'From' : 'Starting from'}</span>
              <p className="text-2xl font-bold text-brand-400">{formatPrice(service.basePrice)}</p>
            </div>
            <div>
              <span className="text-text-muted">Duration</span>
              <p className="text-lg font-semibold">{service.durationMinutes} mins</p>
            </div>
          </div>
        </div>

        <div className="lg:col-span-2">
          <div className="glass-card p-6 sticky top-24">
            <h2 className="text-lg font-semibold mb-4">Add to Cart</h2>
            <AddToCart service={service} />
          </div>
        </div>
      </div>
    </div>
  );
}
