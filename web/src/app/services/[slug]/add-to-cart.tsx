'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { OptionPicker } from '@/components/option-picker';
import { formatPrice } from '@/lib/format';
import type { DetailingService, Cart } from '@/lib/types';

export function AddToCart({ service }: { service: DetailingService }) {
  const router = useRouter();
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
  const [quantity, setQuantity] = useState(1);
  const [selectedTierId, setSelectedTierId] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const hasTiers = (service.priceTiers ?? []).length > 0;
  const selectedTier = hasTiers
    ? (service.priceTiers ?? []).find(t => t.vehicleCategoryId === selectedTierId) ?? null
    : null;

  const basePrice = selectedTier ? Number(selectedTier.price) : Number(service.basePrice);
  const optionsTotal = (service.options ?? [])
    .filter(o => selectedOptions.includes(o.id))
    .reduce((sum, o) => sum + Number(o.price), 0);
  const total = (basePrice + optionsTotal) * quantity;

  async function handleAdd() {
    if (hasTiers && !selectedTierId) {
      setError('Please select your vehicle size to continue.');
      return;
    }
    setLoading(true);
    setError(null);
    try {
      await api<{ cart: Cart }>('/cart/items', {
        method: 'POST',
        body: {
          serviceId: service.id,
          optionIds: selectedOptions,
          quantity,
        },
      });
      setSuccess(true);
      setTimeout(() => router.push('/cart'), 1500);
    } catch (err: unknown) {
      const apiErr = err as { detail?: string };
      setError(apiErr.detail || 'Failed to add to cart');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      {hasTiers && (
        <div>
          <label className="block text-sm font-medium text-text-secondary mb-2">
            Vehicle Size <span className="text-red-400">*</span>
          </label>
          <div className="space-y-2">
            {(service.priceTiers ?? []).map(tier => (
              <button
                key={tier.vehicleCategoryId}
                type="button"
                onClick={() => { setSelectedTierId(tier.vehicleCategoryId); setError(null); }}
                className={`w-full flex items-center justify-between px-4 py-3 rounded border text-sm transition-colors ${
                  selectedTierId === tier.vehicleCategoryId
                    ? 'border-brand-500 bg-brand-500/10 text-foreground'
                    : 'border-border-subtle hover:border-brand-500/50 text-text-secondary hover:text-foreground'
                }`}
              >
                <span className="font-medium">{tier.categoryName}</span>
                <span className={`font-bold ${selectedTierId === tier.vehicleCategoryId ? 'text-brand-400' : ''}`}>
                  {formatPrice(tier.price)}
                </span>
              </button>
            ))}
          </div>
        </div>
      )}

      <OptionPicker
        options={service.options ?? []}
        selectedIds={selectedOptions}
        onChange={setSelectedOptions}
      />

      <div className="flex items-center gap-4">
        <label className="text-sm font-medium text-text-secondary">Quantity</label>
        <div className="flex items-center border border-border-subtle rounded">
          <button
            type="button"
            onClick={() => setQuantity(Math.max(1, quantity - 1))}
            className="px-3 py-1 text-text-secondary hover:bg-surface-hover"
          >
            -
          </button>
          <span className="px-4 py-1 text-sm font-medium">{quantity}</span>
          <button
            type="button"
            onClick={() => setQuantity(quantity + 1)}
            className="px-3 py-1 text-text-secondary hover:bg-surface-hover"
          >
            +
          </button>
        </div>
      </div>

      <div className="border-t border-border-subtle pt-4">
        <div className="flex items-center justify-between mb-4">
          <span className="text-lg font-semibold">Total</span>
          <div className="text-right">
            <span className="text-lg font-bold text-brand-400">{formatPrice(total)}</span>
            {hasTiers && !selectedTier && (
              <p className="text-xs text-text-muted">Select size above</p>
            )}
          </div>
        </div>

        {error && (
          <p className="text-red-400 text-sm mb-3">{error}</p>
        )}

        {success ? (
          <p className="text-green-400 font-medium">Added to cart. Redirecting...</p>
        ) : (
          <button
            onClick={handleAdd}
            disabled={loading}
            className="w-full btn-brand py-3"
          >
            {loading ? 'Adding...' : 'Add to Cart'}
          </button>
        )}
      </div>
    </div>
  );
}
