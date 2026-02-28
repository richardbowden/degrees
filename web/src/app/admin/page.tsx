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
      <h1 className="text-2xl font-bold text-foreground mb-1">Dashboard</h1>
      <p className="text-sm text-text-muted mb-6">Today&apos;s bookings &mdash; {today}</p>

      {error && <p className="text-red-400 text-sm mb-4">{error}</p>}

      <div className="glass-card mb-6 p-5">
        <p className="text-3xl font-bold text-foreground">{bookings.length}</p>
        <p className="text-sm text-text-muted mt-1">
          {bookings.length === 1 ? 'booking' : 'bookings'} today
        </p>
      </div>

      {bookings.length > 0 ? (
        <div className="glass-card overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-surface-input border-b border-border-subtle">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Time</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Customer</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Vehicle</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Services</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Total</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Status</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Payment</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase"></th>
              </tr>
            </thead>
            <tbody>
              {bookings.map(b => <AdminBookingRow key={b.id} booking={b} />)}
            </tbody>
          </table>
        </div>
      ) : (
        !error && <p className="text-sm text-text-muted">No bookings scheduled for today.</p>
      )}
    </div>
  );
}
