import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { BookingDetailClient } from './booking-detail-client';

export default async function AdminBookingDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  let booking: Booking | null = null;
  let error = '';

  try {
    const res = await api<{ booking: Booking }>(`/admin/bookings/${id}`, { token });
    booking = res.booking;
  } catch {
    error = 'Failed to load booking';
  }

  if (error || !booking) {
    return (
      <div>
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Booking Detail</h1>
        <p className="text-red-600 text-sm">{error || 'Booking not found'}</p>
      </div>
    );
  }

  return <BookingDetailClient booking={booking} token={token} />;
}
