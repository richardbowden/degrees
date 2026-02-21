import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { AdminBookingRow } from '@/components/admin-booking-row';

function todayString(): string {
  const d = new Date();
  return d.toISOString().split('T')[0];
}

export default async function AdminDashboard() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  const today = todayString();

  let bookings: Booking[] = [];
  let error = '';

  try {
    const res = await api<{ bookings: Booking[] }>(
      `/admin/bookings?dateFrom=${today}&dateTo=${today}`,
      { token },
    );
    bookings = res.bookings ?? [];
  } catch {
    error = 'Failed to load bookings';
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-1">Dashboard</h1>
      <p className="text-sm text-gray-500 mb-6">Today&apos;s bookings &mdash; {today}</p>

      {error && <p className="text-red-600 text-sm mb-4">{error}</p>}

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 mb-6 p-5">
        <p className="text-3xl font-bold text-gray-900">{bookings.length}</p>
        <p className="text-sm text-gray-500 mt-1">
          {bookings.length === 1 ? 'booking' : 'bookings'} today
        </p>
      </div>

      {bookings.length > 0 ? (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Time</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Customer</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Vehicle</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Services</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Total</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Payment</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase"></th>
              </tr>
            </thead>
            <tbody>
              {bookings.map(b => <AdminBookingRow key={b.id} booking={b} />)}
            </tbody>
          </table>
        </div>
      ) : (
        !error && <p className="text-sm text-gray-500">No bookings scheduled for today.</p>
      )}
    </div>
  );
}
