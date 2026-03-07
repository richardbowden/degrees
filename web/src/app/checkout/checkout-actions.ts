'use server';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { api } from '@/lib/api';
import type { Vehicle, Booking, VehicleFormData } from '@/lib/types';

export async function addVehicleAction(data: VehicleFormData): Promise<{ vehicle?: Vehicle; error?: string }> {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  if (!token) return { error: 'Not authenticated' };

  try {
    const res = await api<{ vehicle: Vehicle }>('/me/vehicles', {
      method: 'POST',
      body: data,
      token,
    });
    return { vehicle: res.vehicle };
  } catch (err: unknown) {
    const apiErr = err as { status?: number; detail?: string };
    if (apiErr.status === 401 || apiErr.status === 403) redirect('/login?redirect=/checkout');
    return { error: apiErr.detail || 'Failed to add vehicle' };
  }
}

export async function submitBookingAction(data: {
  vehicleId: string;
  scheduledDate: string;
  scheduledTime: string;
  notes: string;
  cartSession: string;
}): Promise<{ bookingId?: string; error?: string }> {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  if (!token) return { error: 'Not authenticated' };

  const { cartSession, ...bookingData } = data;

  // Claim the guest cart for this authenticated user before booking.
  // cartSession comes from the client (read from document.cookie) so it always
  // identifies the correct guest cart regardless of server-side cookie state.
  try {
    await api('/cart', { token, cartSession });
  } catch {
    // Non-fatal — if this fails the booking will return a proper error.
  }

  try {
    const res = await api<{ booking: Booking }>('/checkout', {
      method: 'POST',
      body: bookingData,
      token,
      cartSession,
    });
    return { bookingId: res.booking.id };
  } catch (err: unknown) {
    const apiErr = err as { status?: number; detail?: string };
    if (apiErr.status === 401 || apiErr.status === 403) redirect('/login?redirect=/checkout');
    return { error: apiErr.detail || 'Failed to create booking' };
  }
}
