import Link from 'next/link';
import { DetailingService } from '@/lib/types';
import { formatPrice } from '@/lib/format';

export function ServiceCard({ service }: { service: DetailingService }) {
  return (
    <Link
      href={`/services/${service.slug}`}
      className="block border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow"
    >
      <h3 className="text-lg font-semibold mb-2">{service.name}</h3>
      <p className="text-gray-600 text-sm mb-4">{service.shortDesc}</p>
      <div className="flex items-center justify-between text-sm">
        <span className="font-semibold text-gray-900">{formatPrice(service.basePrice)}</span>
        <span className="text-gray-500">{service.durationMinutes} mins</span>
      </div>
    </Link>
  );
}
