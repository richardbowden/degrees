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
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const optionsTotal = (service.options ?? [])
    .filter(o => selectedOptions.includes(o.id))
    .reduce((sum, o) => sum + Number(o.price), 0);
  const total = (Number(service.basePrice) + optionsTotal) * quantity;

  async function handleAdd() {
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
      <OptionPicker
        options={service.options ?? []}
        selectedIds={selectedOptions}
        onChange={setSelectedOptions}
      />

      <div className="flex items-center gap-4">
        <label className="text-sm font-medium text-gray-700">Quantity</label>
        <div className="flex items-center border border-gray-300 rounded">
          <button
            type="button"
            onClick={() => setQuantity(Math.max(1, quantity - 1))}
            className="px-3 py-1 text-gray-600 hover:bg-gray-100"
          >
            -
          </button>
          <span className="px-4 py-1 text-sm font-medium">{quantity}</span>
          <button
            type="button"
            onClick={() => setQuantity(quantity + 1)}
            className="px-3 py-1 text-gray-600 hover:bg-gray-100"
          >
            +
          </button>
        </div>
      </div>

      <div className="border-t border-gray-200 pt-4">
        <div className="flex items-center justify-between mb-4">
          <span className="text-lg font-semibold">Total</span>
          <span className="text-lg font-bold">{formatPrice(total)}</span>
        </div>

        {error && (
          <p className="text-red-600 text-sm mb-3">{error}</p>
        )}

        {success ? (
          <p className="text-green-600 font-medium">Added to cart. Redirecting...</p>
        ) : (
          <button
            onClick={handleAdd}
            disabled={loading}
            className="w-full bg-gray-900 text-white py-3 rounded font-semibold hover:bg-gray-800 disabled:opacity-50"
          >
            {loading ? 'Adding...' : 'Add to Cart'}
          </button>
        )}
      </div>
    </div>
  );
}
