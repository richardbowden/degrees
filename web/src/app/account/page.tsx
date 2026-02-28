import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking, CustomerProfile, User } from '@/lib/types';
import { StatusBadge } from '@/components/status-badge';
import { formatDate, formatTime, formatPrice } from '@/lib/format';

export default async function AccountDashboard() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const { profile } = await api<{ profile: CustomerProfile }>('/me/profile', { token });

  const [userRes, bookingsRes] = await Promise.all([
    api<{ user: User }>(`/user/${profile.userId}`, { token }),
    api<{ bookings: Booking[] }>('/me/bookings', { token }).catch(() => ({ bookings: [] as Booking[] })),
  ]);

  const user = userRes.user;
  const upcomingBooking = bookingsRes.bookings
    .filter(b => b.status !== 'cancelled' && b.status !== 'completed')
    .sort((a, b) => a.scheduledDate.localeCompare(b.scheduledDate) || a.scheduledTime.localeCompare(b.scheduledTime))
    [0];

  return (
    <div>
      <h1 className="text-2xl font-bold text-foreground mb-6">
        Welcome back, {user.firstName}
      </h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Next Booking</h2>
          {upcomingBooking ? (
            <div>
              <div className="flex items-center gap-3 mb-2">
                <span className="font-medium text-foreground">
                  {formatDate(upcomingBooking.scheduledDate)}
                </span>
                <span className="text-text-muted">
                  {formatTime(upcomingBooking.scheduledTime)}
                </span>
                <StatusBadge status={upcomingBooking.status} />
              </div>
              {upcomingBooking.vehicle && (
                <p className="text-sm text-text-secondary mb-2">
                  {upcomingBooking.vehicle.year} {upcomingBooking.vehicle.make} {upcomingBooking.vehicle.model}
                  {upcomingBooking.vehicle.rego && ` (${upcomingBooking.vehicle.rego})`}
                </p>
              )}
              {upcomingBooking.services && upcomingBooking.services.length > 0 && (
                <p className="text-sm text-text-secondary mb-3">
                  {upcomingBooking.services.map(s => s.serviceName).join(', ')}
                </p>
              )}
              <p className="text-sm font-medium text-foreground mb-3">
                Total: {formatPrice(upcomingBooking.totalAmount)}
              </p>
              <Link
                href={`/account/bookings/${upcomingBooking.id}`}
                className="text-sm font-medium text-brand-400 underline hover:text-brand-400"
              >
                View Details
              </Link>
            </div>
          ) : (
            <div>
              <p className="text-text-muted mb-3">No upcoming bookings</p>
              <Link
                href="/services"
                className="text-sm font-medium text-brand-400 underline hover:text-brand-400"
              >
                Browse Services
              </Link>
            </div>
          )}
        </div>

        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Quick Links</h2>
          <div className="space-y-3">
            <Link
              href="/account/bookings"
              className="block text-sm font-medium text-text-secondary hover:text-foreground"
            >
              View All Bookings
            </Link>
            <Link
              href="/account/profile"
              className="block text-sm font-medium text-text-secondary hover:text-foreground"
            >
              Edit Profile
            </Link>
            <Link
              href="/account/vehicles"
              className="block text-sm font-medium text-text-secondary hover:text-foreground"
            >
              My Vehicles
            </Link>
            <Link
              href="/account/history"
              className="block text-sm font-medium text-text-secondary hover:text-foreground"
            >
              Service History
            </Link>
            <Link
              href="/services"
              className="block text-sm font-medium text-text-secondary hover:text-foreground"
            >
              Browse Services
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
