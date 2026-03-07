import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import { formatPrice, formatDate, formatTime } from '@/lib/format';
import type { Booking } from '@/lib/types';

interface Props {
  searchParams: Promise<{ booking_id?: string }>;
}

export default async function DepositPage({ searchParams }: Props) {
  const params = await searchParams;
  const bookingId = params.booking_id ?? '';

  if (!bookingId) {
    redirect('/account/bookings');
  }

  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  if (!token) {
    redirect('/login');
  }

  let booking: Booking | null = null;
  try {
    const res = await api<{ booking: Booking }>(`/me/bookings/${bookingId}`, { token });
    booking = res.booking;
  } catch {
    // If we can't fetch booking details, still show a basic confirmation
  }

  const depositAmount = booking ? Number(booking.depositAmount) : 0;
  const totalAmount = booking ? Number(booking.totalAmount) : 0;

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <div className="w-16 h-16 bg-brand-500/20 rounded-full flex items-center justify-center mx-auto mb-6">
        <span className="text-brand-400 text-2xl">&#10003;</span>
      </div>

      <h1 className="text-2xl font-bold text-center mb-2">Booking Confirmed!</h1>
      <p className="text-text-secondary text-center mb-8">
        We&apos;ll see you on the day. Payment is collected when we arrive.
      </p>

      {booking && (
        <div className="border border-border-subtle rounded-lg divide-y divide-border-subtle mb-8">
          {/* Services */}
          {(booking.services ?? []).length > 0 && (
            <div className="p-4">
              <p className="text-xs text-text-muted uppercase tracking-wide mb-2">Services</p>
              <div className="space-y-1">
                {booking.services.map((svc, i) => (
                  <div key={i} className="flex justify-between text-sm">
                    <span className="text-foreground">{svc.serviceName}</span>
                    <span className="text-text-secondary">{formatPrice(svc.priceAtBooking)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Date & Time */}
          {(booking.scheduledDate || booking.scheduledTime) && (
            <div className="p-4 grid grid-cols-2 gap-4">
              {booking.scheduledDate && (
                <div>
                  <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Date</p>
                  <p className="text-sm text-foreground">{formatDate(booking.scheduledDate)}</p>
                </div>
              )}
              {booking.scheduledTime && (
                <div>
                  <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Time</p>
                  <p className="text-sm text-foreground">{formatTime(booking.scheduledTime)}</p>
                </div>
              )}
            </div>
          )}

          {/* Vehicle */}
          {booking.vehicle && (booking.vehicle.make || booking.vehicle.model) && (
            <div className="p-4">
              <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Vehicle</p>
              <p className="text-sm text-foreground">
                {[booking.vehicle.make, booking.vehicle.model, booking.vehicle.rego && `(${booking.vehicle.rego})`]
                  .filter(Boolean)
                  .join(' ')}
              </p>
            </div>
          )}

          {/* Payment */}
          <div className="p-4 space-y-1">
            <div className="flex justify-between text-sm">
              <span className="text-text-secondary">Total</span>
              <span>{formatPrice(totalAmount)}</span>
            </div>
            <div className="flex justify-between text-sm font-medium text-brand-400">
              <span>Deposit due on arrival (30%)</span>
              <span>{formatPrice(depositAmount)}</span>
            </div>
          </div>
        </div>
      )}

      {!booking && (
        <div className="border border-border-subtle rounded p-4 mb-8 text-sm text-text-secondary text-center">
          <p className="font-mono">{bookingId}</p>
        </div>
      )}

      <div className="flex flex-col gap-3">
        <Link href="/account/bookings" className="block w-full btn-brand text-center py-3">
          View My Bookings
        </Link>
        <Link href="/services" className="text-center text-sm text-text-secondary hover:underline">
          Browse more services
        </Link>
      </div>
    </div>
  );
}
