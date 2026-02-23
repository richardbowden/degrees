'use client';

import Link from 'next/link';
import { Cart } from '@/lib/types';
import { formatPrice } from '@/lib/format';

export function CartSummary({ cart }: { cart: Cart }) {
  if (!cart.items || cart.items.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-text-muted mb-4">Your cart is empty</p>
        <Link href="/services" className="text-brand-400 hover:underline">
          Browse services
        </Link>
      </div>
    );
  }

  return (
    <div>
      <div className="divide-y divide-border-subtle">
        {cart.items.map(item => (
          <div key={item.id} className="py-4 flex items-center justify-between">
            <div>
              <p className="font-medium">{item.serviceName}</p>
              <p className="text-sm text-text-muted">Qty: {item.quantity}</p>
              {item.optionIds.length > 0 && (
                <p className="text-xs text-text-muted">{item.optionIds.length} add-on(s)</p>
              )}
            </div>
            <span className="font-semibold text-brand-400">{formatPrice(item.servicePrice)}</span>
          </div>
        ))}
      </div>
      <div className="border-t border-border-subtle pt-4 mt-4 flex items-center justify-between">
        <span className="text-lg font-bold">Subtotal</span>
        <span className="text-lg font-bold text-brand-400">{formatPrice(cart.subtotal)}</span>
      </div>
      <Link
        href="/checkout"
        className="block mt-6 w-full btn-brand text-center py-3"
      >
        Proceed to Checkout
      </Link>
    </div>
  );
}
