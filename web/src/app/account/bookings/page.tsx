import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { StatusBadge } from '@/components/status-badge';
import { formatDate, formatTime, formatPrice } from '@/lib/format';

export default async function BookingsPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const { bookings } = await api<{ bookings: Booking[] }>('/me/bookings', { token });

  const sorted = [...bookings].sort(
    (a, b) => b.scheduledDate.localeCompare(a.scheduledDate) || b.scheduledTime.localeCompare(a.scheduledTime)
  );

  return (
    <div>
      <h1 className="text-2xl font-bold text-foreground mb-6">Bookings</h1>

      {sorted.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-text-muted mb-4">You don't have any bookings yet.</p>
          <Link
            href="/services"
            className="btn-brand inline-block px-4 py-2 text-sm font-medium rounded-md"
          >
            Browse Services
          </Link>
        </div>
      ) : (
        <div className="border border-border-subtle rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-surface-input">
              <tr>
                <th className="text-left px-4 py-3 font-medium text-text-secondary">Date</th>
                <th className="text-left px-4 py-3 font-medium text-text-secondary">Time</th>
                <th className="text-left px-4 py-3 font-medium text-text-secondary">Services</th>
                <th className="text-left px-4 py-3 font-medium text-text-secondary">Status</th>
                <th className="text-right px-4 py-3 font-medium text-text-secondary">Total</th>
                <th className="px-4 py-3"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border-subtle">
              {sorted.map(booking => (
                <tr key={booking.id} className="hover:bg-surface-hover">
                  <td className="px-4 py-3 text-foreground">{formatDate(booking.scheduledDate)}</td>
                  <td className="px-4 py-3 text-text-secondary">{formatTime(booking.scheduledTime)}</td>
                  <td className="px-4 py-3 text-text-secondary">
                    {booking.services?.map(s => s.serviceName).join(', ') || '-'}
                  </td>
                  <td className="px-4 py-3">
                    <StatusBadge status={booking.status} />
                  </td>
                  <td className="px-4 py-3 text-right text-foreground font-medium">
                    {formatPrice(booking.totalAmount)}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <Link
                      href={`/account/bookings/${booking.id}`}
                      className="text-text-secondary hover:text-foreground font-medium"
                    >
                      View
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
