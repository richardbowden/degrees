import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { api } from '@/lib/api';
import { CheckoutClient } from './checkout-client';
import type { Vehicle, Cart, DetailingService } from '@/lib/types';

export default async function CheckoutPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  const cartSession = cookieStore.get('cart_session')?.value;

  if (!token) {
    redirect('/login?redirect=/checkout');
  }

  try {
    const [vehiclesData, cartData, catalogueData] = await Promise.all([
      api<{ vehicles: Vehicle[] }>('/me/vehicles', { token }),
      api<{ cart: Cart }>('/cart', { token, cartSession }),
      api<{ services: DetailingService[] }>('/catalogue'),
    ]);

    const serviceDurations: Record<string, number> = {};
    for (const svc of catalogueData.services ?? []) {
      serviceDurations[svc.id] = svc.durationMinutes;
    }

    return (
      <CheckoutClient
        initialVehicles={vehiclesData.vehicles ?? []}
        initialCart={cartData.cart}
        serviceDurations={serviceDurations}
      />
    );
  } catch (err) {
    const apiErr = err as { status?: number };
    if (!apiErr.status || apiErr.status === 401 || apiErr.status === 403) {
      redirect('/login?redirect=/checkout');
    }
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-red-400">Failed to load checkout. Please try again.</p>
      </div>
    );
  }
}
