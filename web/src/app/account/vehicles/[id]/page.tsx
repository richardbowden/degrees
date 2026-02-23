import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Vehicle, ServiceRecord, Booking } from '@/lib/types';
import { formatDate, formatTime } from '@/lib/format';
import { StatusBadge } from '@/components/status-badge';
import { VehicleActions } from './vehicle-actions';

interface VehicleDetailPageProps {
  params: Promise<{ id: string }>;
}

export default async function VehicleDetailPage({ params }: VehicleDetailPageProps) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const [vehicleRes, historyRes, bookingsRes] = await Promise.all([
    api<{ vehicle: Vehicle }>(`/me/vehicles/${id}`, { token }),
    api<{ records: ServiceRecord[] }>('/me/history', { token }).catch(() => ({ records: [] as ServiceRecord[] })),
    api<{ bookings: Booking[] }>('/me/bookings', { token }).catch(() => ({ bookings: [] as Booking[] })),
  ]);

  const vehicle = vehicleRes.vehicle;
  const bookingsMap = new Map(bookingsRes.bookings.map(b => [b.id, b]));

  const completedHistory = historyRes.records
    .filter(r => r.vehicleId === id)
    .sort((a, b) => b.completedDate.localeCompare(a.completedDate));

  const upcomingBookings = bookingsRes.bookings
    .filter(b => b.vehicleId === id && b.status !== 'cancelled' && b.status !== 'completed')
    .sort((a, b) => a.scheduledDate.localeCompare(b.scheduledDate) || a.scheduledTime.localeCompare(b.scheduledTime));

  return (
    <div>
      <Link
        href="/account/vehicles"
        className="text-sm text-text-muted hover:text-white mb-4 inline-block"
      >
        &larr; Back to Vehicles
      </Link>

      <div className="border border-border-subtle rounded-lg p-6 mb-8">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h1 className="text-2xl font-bold text-white">
              {vehicle.year} {vehicle.make} {vehicle.model}
            </h1>
            {vehicle.isPrimary && (
              <span className="inline-block mt-1 text-xs bg-white/10 text-text-secondary px-2 py-0.5 rounded-full">
                Primary
              </span>
            )}
          </div>
        </div>
        <dl className="grid grid-cols-2 sm:grid-cols-3 gap-4 text-sm">
          <div>
            <dt className="text-text-muted">Colour</dt>
            <dd className="text-white">{vehicle.colour}</dd>
          </div>
          {vehicle.rego && (
            <div>
              <dt className="text-text-muted">Registration</dt>
              <dd className="text-white">{vehicle.rego}</dd>
            </div>
          )}
          {vehicle.paintType && (
            <div>
              <dt className="text-text-muted">Paint Type</dt>
              <dd className="text-white">{vehicle.paintType}</dd>
            </div>
          )}
        </dl>
        {vehicle.conditionNotes && (
          <div className="mt-4 text-sm">
            <p className="text-text-muted">Condition Notes</p>
            <p className="text-text-secondary mt-1">{vehicle.conditionNotes}</p>
          </div>
        )}
        <div className="mt-6">
          <VehicleActions vehicle={vehicle} />
        </div>
      </div>

      {upcomingBookings.length > 0 && (
        <section className="mb-8">
          <h2 className="text-lg font-semibold text-white mb-4">Upcoming Bookings</h2>
          <div className="space-y-3">
            {upcomingBookings.map(booking => (
              <Link
                key={booking.id}
                href={`/account/bookings/${booking.id}`}
                className="block border border-border-subtle rounded-lg p-4 hover:bg-white/5"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <div className="flex items-center gap-3">
                      <span className="font-medium text-white">
                        {formatDate(booking.scheduledDate)}
                      </span>
                      <span className="text-text-muted text-sm">
                        {formatTime(booking.scheduledTime)}
                      </span>
                      <StatusBadge status={booking.status} />
                    </div>
                    {booking.services && booking.services.length > 0 && (
                      <p className="text-sm text-text-muted mt-1">
                        {booking.services.map(s => s.serviceName).join(', ')}
                      </p>
                    )}
                  </div>
                  <span className="text-text-muted">&rarr;</span>
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      <section>
        <h2 className="text-lg font-semibold text-white mb-4">Detailing History</h2>
        {completedHistory.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-text-muted mb-3">No detailing history yet.</p>
            <Link
              href="/services"
              className="text-sm font-medium text-brand-400 underline hover:text-brand-400"
            >
              Book a service
            </Link>
          </div>
        ) : (
          <div className="space-y-3">
            {completedHistory.map(record => {
              const booking = bookingsMap.get(record.bookingId);
              return (
                <Link
                  key={record.id}
                  href={`/account/history/${record.id}`}
                  className="block border border-border-subtle rounded-lg p-4 hover:bg-white/5"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-white">
                        {formatDate(record.completedDate)}
                      </p>
                      {booking?.services && booking.services.length > 0 && (
                        <p className="text-sm text-text-muted mt-1">
                          {booking.services.map(s => s.serviceName).join(', ')}
                        </p>
                      )}
                    </div>
                    <span className="text-text-muted">&rarr;</span>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </section>
    </div>
  );
}
