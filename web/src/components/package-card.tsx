import Link from 'next/link';
import { DetailingService } from '@/lib/types';
import { formatPrice } from '@/lib/format';
import { MarkdownContent } from './markdown-content';

function PriceGrid({ service }: { service: DetailingService }) {
  const tiers = service.priceTiers ?? [];
  if (tiers.length === 0) return null;
  return (
    <div className="mb-6 rounded-lg bg-brand-500/5 border border-brand-500/20 overflow-hidden">
      {tiers.map((tier, i) => (
        <div
          key={tier.vehicleCategoryId}
          className={`flex items-center justify-between px-4 py-2.5 ${
            i > 0 ? 'border-t border-brand-500/10' : ''
          }`}
        >
          <span className="text-sm text-text-secondary">{tier.categoryName}</span>
          <span className="text-sm font-bold text-brand-400">{formatPrice(tier.price)}</span>
        </div>
      ))}
    </div>
  );
}

export function PackageCard({ service }: { service: DetailingService }) {
  const hasTiers = service.priceTiers && service.priceTiers.length > 0;

  return (
    <div className="glass-card p-8 md:p-10 flex flex-col">
      <h3 className="text-2xl font-bold text-brand-400 mb-2">{service.name}</h3>
      <p className="text-text-secondary leading-relaxed mb-6">{service.shortDesc}</p>

      <div className="flex-1 mb-6">
        <MarkdownContent content={service.description} />
      </div>

      {service.options && service.options.length > 0 && (
        <div className="border-l-2 border-brand-500 pl-4 mb-6 bg-brand-500/5 py-3 pr-4 rounded-r-lg space-y-1">
          {service.options.map(opt => (
            <p key={opt.id} className="text-sm text-text-secondary">
              <strong className="text-foreground">{opt.name}:</strong> {opt.description}{' '}
              <span className="text-brand-400 font-semibold">{formatPrice(opt.price)}</span>
            </p>
          ))}
        </div>
      )}

      <PriceGrid service={service} />

      <div className="flex items-center justify-between pt-6 border-t border-border-subtle">
        <div>
          <span className="text-text-muted text-sm">{hasTiers ? 'From' : 'Starting from'}</span>
          <p className="text-2xl font-bold text-brand-400">{formatPrice(service.basePrice)}</p>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-text-muted text-sm hidden sm:inline">{service.durationMinutes} mins</span>
          <Link
            href={`/services/${service.slug}`}
            className="btn-brand px-6 py-3 text-sm font-medium"
          >
            Book Now
          </Link>
        </div>
      </div>
    </div>
  );
}
