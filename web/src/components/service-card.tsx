import Link from 'next/link';
import { DetailingService } from '@/lib/types';
import { formatPrice } from '@/lib/format';
import { extractFeaturePreview, extractSections } from '@/lib/markdown-utils';

export function ServiceCard({ service }: { service: DetailingService }) {
  const sections = extractSections(service.description);
  const features = extractFeaturePreview(service.description, 4);
  const addOnCount = (service.options ?? []).length;

  return (
    <Link
      href={`/services/${service.slug}`}
      className="block glass-card card-hover p-6 flex flex-col"
    >
      <h3 className="text-lg font-semibold mb-2">{service.name}</h3>
      <p className="text-text-secondary text-sm mb-4">{service.shortDesc}</p>

      {sections.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mb-4">
          {sections.map(s => (
            <span key={s} className="px-2 py-0.5 text-xs rounded-full bg-surface-input border border-border-subtle text-text-muted">
              {s}
            </span>
          ))}
        </div>
      )}

      {features.length > 0 && (
        <ul className="space-y-1.5 mb-4 text-sm text-text-secondary">
          {features.map((f, i) => (
            <li key={i} className="flex items-start gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-3.5 h-3.5 text-brand-400 mt-0.5 shrink-0">
                <path strokeLinecap="round" strokeLinejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
              </svg>
              <span>{f}</span>
            </li>
          ))}
        </ul>
      )}

      <div className="flex items-center justify-between text-sm mt-auto pt-4 border-t border-border-subtle">
        <span className="font-semibold text-brand-400">{formatPrice(service.basePrice)}</span>
        <div className="flex items-center gap-3">
          {addOnCount > 0 && (
            <span className="text-text-muted text-xs">{addOnCount} add-on{addOnCount !== 1 ? 's' : ''}</span>
          )}
          <span className="text-text-muted">{service.durationMinutes} mins</span>
        </div>
      </div>
    </Link>
  );
}
