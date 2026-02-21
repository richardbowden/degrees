import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { ServiceRecord, Booking } from '@/lib/types';
import { formatDate } from '@/lib/format';

export default async function HistoryPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const [historyRes, bookingsRes] = await Promise.all([
    api<{ records: ServiceRecord[] }>('/me/history', { token }),
    api<{ bookings: Booking[] }>('/me/bookings', { token }).catch(() => ({ bookings: [] as Booking[] })),
  ]);

  const records = historyRes.records;
  const bookingsMap = new Map(bookingsRes.bookings.map(b => [b.id, b]));

  const sorted = [...records].sort(
    (a, b) => b.completedDate.localeCompare(a.completedDate)
  );

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Service History</h1>

      {sorted.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500">No service history yet.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {sorted.map(record => {
            const booking = bookingsMap.get(record.bookingId);
            return (
              <Link
                key={record.id}
                href={`/account/history/${record.id}`}
                className="block border border-gray-200 rounded-lg p-4 hover:bg-gray-50"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-gray-900">
                      {formatDate(record.completedDate)}
                    </p>
                    {booking?.vehicle && (
                      <p className="text-sm text-gray-500 mt-1">
                        {booking.vehicle.year} {booking.vehicle.make} {booking.vehicle.model}
                        {booking.vehicle.rego && ` (${booking.vehicle.rego})`}
                      </p>
                    )}
                    {booking?.services && booking.services.length > 0 && (
                      <p className="text-sm text-gray-400 mt-0.5">
                        {booking.services.map(s => s.serviceName).join(', ')}
                      </p>
                    )}
                  </div>
                  <span className="text-gray-400">&rarr;</span>
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
