import Link from 'next/link';
import { DetailingService } from '@/lib/types';
import { formatPrice } from '@/lib/format';

export function ServiceCard({ service }: { service: DetailingService }) {
  return (
    <Link
      href={`/services/${service.slug}`}
      className="block glass-card card-hover p-6"
    >
      <h3 className="text-lg font-semibold mb-2">{service.name}</h3>
      <p className="text-text-secondary text-sm mb-4">{service.shortDesc}</p>
      <div className="flex items-center justify-between text-sm">
        <span className="font-semibold text-brand-400">{formatPrice(service.basePrice)}</span>
        <span className="text-text-muted">{service.durationMinutes} mins</span>
      </div>
    </Link>
  );
}
