'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { formatPrice, formatTime, formatDate } from '@/lib/format';
import { StatusBadge } from '@/components/status-badge';

function weekRange(): { from: string; to: string } {
  const now = new Date();
  const day = now.getDay();
  const monday = new Date(now);
  monday.setDate(now.getDate() - ((day + 6) % 7));
  const sunday = new Date(monday);
  sunday.setDate(monday.getDate() + 6);
  return {
    from: monday.toISOString().split('T')[0],
    to: sunday.toISOString().split('T')[0],
  };
}

const STATUSES = ['all', 'pending', 'confirmed', 'in_progress', 'completed', 'cancelled'];

export function BookingsClient({ token }: { token: string }) {
  const range = weekRange();
  const [dateFrom, setDateFrom] = useState(range.from);
  const [dateTo, setDateTo] = useState(range.to);
  const [statusFilter, setStatusFilter] = useState('all');
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchBookings = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await api<{ bookings: Booking[] }>(
        `/admin/bookings?dateFrom=${dateFrom}&dateTo=${dateTo}`,
        { token },
      );
      setBookings(res.bookings ?? []);
    } catch {
      setError('Failed to load bookings');
    } finally {
      setLoading(false);
    }
  }, [dateFrom, dateTo, token]);

  useEffect(() => {
    fetchBookings();
  }, [fetchBookings]);

  const filtered = statusFilter === 'all'
    ? bookings
    : bookings.filter(b => b.status === statusFilter);

  return (
    <div>
      <div className="flex flex-wrap gap-4 items-end mb-6">
        <div>
          <label className="block text-xs font-medium text-text-muted mb-1">From</label>
          <input
            type="date"
            value={dateFrom}
            onChange={e => setDateFrom(e.target.value)}
            className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-text-muted mb-1">To</label>
          <input
            type="date"
            value={dateTo}
            onChange={e => setDateTo(e.target.value)}
            className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-text-muted mb-1">Status</label>
          <select
            value={statusFilter}
            onChange={e => setStatusFilter(e.target.value)}
            className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
          >
            {STATUSES.map(s => (
              <option key={s} value={s}>{s === 'all' ? 'All Statuses' : s.replace('_', ' ')}</option>
            ))}
          </select>
        </div>
      </div>

      {error && <p className="text-red-400 text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-sm text-text-muted">Loading bookings...</p>
      ) : filtered.length === 0 ? (
        <p className="text-sm text-text-muted">No bookings found for this period.</p>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-white/5 border-b border-border-subtle">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Date</th>
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
              {filtered.map(b => (
                <tr key={b.id} className="border-b border-white/5 hover:bg-white/5">
                  <td className="py-3 px-4 text-sm">{formatDate(b.scheduledDate)}</td>
                  <td className="py-3 px-4 text-sm">{formatTime(b.scheduledTime)}</td>
                  <td className="py-3 px-4 text-sm font-medium">
                    <Link href={`/admin/bookings/${b.id}`} className="text-brand-400 hover:text-brand-500">
                      {b.customer?.name ?? 'Unknown'}
                    </Link>
                  </td>
                  <td className="py-3 px-4 text-sm text-text-secondary">
                    {b.vehicle ? `${b.vehicle.year} ${b.vehicle.make} ${b.vehicle.model}` : 'N/A'}
                  </td>
                  <td className="py-3 px-4 text-sm text-text-secondary">
                    {b.services?.map(s => s.serviceName).join(', ') ?? 'N/A'}
                  </td>
                  <td className="py-3 px-4 text-sm">{formatPrice(b.totalAmount)}</td>
                  <td className="py-3 px-4"><StatusBadge status={b.status} /></td>
                  <td className="py-3 px-4"><StatusBadge status={b.paymentStatus} /></td>
                  <td className="py-3 px-4 text-sm">
                    <Link href={`/admin/bookings/${b.id}`} className="text-brand-400 hover:underline">
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
