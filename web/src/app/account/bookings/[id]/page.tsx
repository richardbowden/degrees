import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { StatusBadge } from '@/components/status-badge';
import { formatDate, formatTime, formatPrice } from '@/lib/format';
import { CancelButton } from './cancel-button';

export default async function BookingDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const { booking } = await api<{ booking: Booking }>(`/me/bookings/${id}`, { token });

  const canCancel = booking.status === 'pending' || booking.status === 'confirmed';

  return (
    <div>
      <Link
        href="/account/bookings"
        className="text-sm text-text-muted hover:text-text-secondary mb-4 inline-block"
      >
        &larr; Back to Bookings
      </Link>

      <div className="flex items-center gap-4 mb-6">
        <h1 className="text-2xl font-bold text-foreground">Booking</h1>
        <StatusBadge status={booking.status} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Schedule & Vehicle */}
        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Details</h2>
          <dl className="space-y-3 text-sm">
            <div className="flex justify-between">
              <dt className="text-text-muted">Date</dt>
              <dd className="text-foreground font-medium">{formatDate(booking.scheduledDate)}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-text-muted">Time</dt>
              <dd className="text-foreground">{formatTime(booking.scheduledTime)}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-text-muted">Duration</dt>
              <dd className="text-foreground">{booking.estimatedDurationMins} minutes</dd>
            </div>
            {booking.vehicle && (
              <div className="flex justify-between">
                <dt className="text-text-muted">Vehicle</dt>
                <dd className="text-foreground">
                  {booking.vehicle.year} {booking.vehicle.make} {booking.vehicle.model}
                  {booking.vehicle.rego && ` (${booking.vehicle.rego})`}
                </dd>
              </div>
            )}
            <div className="flex justify-between">
              <dt className="text-text-muted">Payment Status</dt>
              <dd><StatusBadge status={booking.paymentStatus} /></dd>
            </div>
          </dl>
        </div>

        {/* Payment */}
        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Payment</h2>
          <dl className="space-y-3 text-sm">
            <div className="flex justify-between">
              <dt className="text-text-muted">Subtotal</dt>
              <dd className="text-foreground">{formatPrice(booking.subtotal)}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-text-muted">Deposit</dt>
              <dd className="text-foreground">{formatPrice(booking.depositAmount)}</dd>
            </div>
            <div className="flex justify-between border-t border-border-subtle pt-3">
              <dt className="text-foreground font-medium">Total</dt>
              <dd className="text-foreground font-bold">{formatPrice(booking.totalAmount)}</dd>
            </div>
          </dl>
        </div>

        {/* Services */}
        {booking.services && booking.services.length > 0 && (
          <div className="glass-card p-6 lg:col-span-2">
            <h2 className="text-lg font-semibold text-foreground mb-4">Services</h2>
            <div className="space-y-3">
              {booking.services.map(service => (
                <div key={service.serviceId} className="flex items-center justify-between text-sm">
                  <div>
                    <p className="text-foreground font-medium">{service.serviceName}</p>
                    {service.options && service.options.length > 0 && (
                      <p className="text-text-muted text-xs mt-0.5">
                        Options: {service.options.join(', ')}
                      </p>
                    )}
                  </div>
                  <span className="text-foreground">{formatPrice(service.price)}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Notes */}
        {booking.notes && (
          <div className="glass-card p-6 lg:col-span-2">
            <h2 className="text-lg font-semibold text-foreground mb-4">Notes</h2>
            <p className="text-sm text-text-secondary whitespace-pre-wrap">{booking.notes}</p>
          </div>
        )}
      </div>

      {canCancel && (
        <div className="mt-6">
          <CancelButton bookingId={booking.id} />
        </div>
      )}
    </div>
  );
}
