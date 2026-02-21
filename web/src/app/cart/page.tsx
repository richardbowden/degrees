'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import { formatPrice } from '@/lib/format';
import type { Cart } from '@/lib/types';

export default function CartPage() {
  const [cart, setCart] = useState<Cart | null>(null);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchCart = useCallback(async () => {
    try {
      const data = await api<{ cart: Cart }>('/cart');
      setCart(data.cart);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      if (apiErr.status === 404) {
        setCart(null);
      } else {
        setError(apiErr.detail || 'Failed to load cart');
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCart();
  }, [fetchCart]);

  async function removeItem(itemId: string) {
    setActionLoading(itemId);
    setError(null);
    try {
      const data = await api<{ cart: Cart }>(`/cart/items/${itemId}`, { method: 'DELETE' });
      setCart(data.cart);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to remove item');
    } finally {
      setActionLoading(null);
    }
  }

  async function updateQuantity(itemId: string, quantity: number) {
    if (quantity < 1) return;
    setActionLoading(itemId);
    setError(null);
    try {
      const data = await api<{ cart: Cart }>(`/cart/items/${itemId}`, {
        method: 'PUT',
        body: { quantity },
      });
      setCart(data.cart);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to update quantity');
    } finally {
      setActionLoading(null);
    }
  }

  if (loading) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <p className="text-gray-500">Loading cart...</p>
      </div>
    );
  }

  const items = cart?.items ?? [];

  if (items.length === 0) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">Your Cart</h1>
        <p className="text-gray-500 mb-6">Your cart is empty.</p>
        <Link href="/services" className="text-gray-900 font-medium hover:underline">
          Browse services
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">Your Cart</h1>

      {error && (
        <p className="text-red-600 text-sm mb-4">{error}</p>
      )}

      <div className="divide-y divide-gray-200">
        {items.map(item => (
          <div key={item.id} className="py-4 flex items-start justify-between gap-4">
            <div className="flex-1">
              <p className="font-medium">{item.serviceName}</p>
              {item.optionIds.length > 0 && (
                <p className="text-xs text-gray-400 mt-1">{item.optionIds.length} add-on(s)</p>
              )}
              <div className="flex items-center gap-2 mt-2">
                <button
                  type="button"
                  onClick={() => updateQuantity(item.id, Number(item.quantity) - 1)}
                  disabled={actionLoading === item.id || Number(item.quantity) <= 1}
                  className="w-7 h-7 border border-gray-300 rounded text-sm hover:bg-gray-100 disabled:opacity-50"
                >
                  -
                </button>
                <span className="text-sm font-medium w-6 text-center">{item.quantity}</span>
                <button
                  type="button"
                  onClick={() => updateQuantity(item.id, Number(item.quantity) + 1)}
                  disabled={actionLoading === item.id}
                  className="w-7 h-7 border border-gray-300 rounded text-sm hover:bg-gray-100 disabled:opacity-50"
                >
                  +
                </button>
              </div>
            </div>
            <div className="text-right">
              <p className="font-semibold">{formatPrice(Number(item.servicePrice) * Number(item.quantity))}</p>
              <button
                type="button"
                onClick={() => removeItem(item.id)}
                disabled={actionLoading === item.id}
                className="text-xs text-red-600 hover:underline mt-1 disabled:opacity-50"
              >
                {actionLoading === item.id ? 'Removing...' : 'Remove'}
              </button>
            </div>
          </div>
        ))}
      </div>

      <div className="border-t border-gray-300 pt-4 mt-4 flex items-center justify-between">
        <span className="text-lg font-bold">Subtotal</span>
        <span className="text-lg font-bold">{formatPrice(cart!.subtotal)}</span>
      </div>

      <Link
        href="/checkout"
        className="block mt-6 w-full bg-gray-900 text-white text-center py-3 rounded font-semibold hover:bg-gray-800"
      >
        Proceed to Checkout
      </Link>
    </div>
  );
}
