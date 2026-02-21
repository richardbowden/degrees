'use client';

import Link from 'next/link';
import { CartIcon } from '@/components/cart-icon';

export function Nav() {
  return (
    <nav className="border-b border-gray-200 bg-white">
      <div className="max-w-6xl mx-auto px-4 flex items-center justify-between h-16">
        <Link href="/" className="text-xl font-bold text-gray-900">
          40 Degrees
        </Link>
        <div className="flex items-center gap-6">
          <Link href="/services" className="text-gray-700 hover:text-gray-900">
            Services
          </Link>
          <CartIcon />
          <Link href="/login" className="text-gray-700 hover:text-gray-900">
            Login
          </Link>
        </div>
      </div>
    </nav>
  );
}
