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
          <label className="block text-xs font-medium text-gray-500 mb-1">From</label>
          <input
            type="date"
            value={dateFrom}
            onChange={e => setDateFrom(e.target.value)}
            className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">To</label>
          <input
            type="date"
            value={dateTo}
            onChange={e => setDateTo(e.target.value)}
            className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">Status</label>
          <select
            value={statusFilter}
            onChange={e => setStatusFilter(e.target.value)}
            className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
          >
            {STATUSES.map(s => (
              <option key={s} value={s}>{s === 'all' ? 'All Statuses' : s.replace('_', ' ')}</option>
            ))}
          </select>
        </div>
      </div>

      {error && <p className="text-red-600 text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-sm text-gray-500">Loading bookings...</p>
      ) : filtered.length === 0 ? (
        <p className="text-sm text-gray-500">No bookings found for this period.</p>
      ) : (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Date</th>
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
              {filtered.map(b => (
                <tr key={b.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-3 px-4 text-sm">{formatDate(b.scheduledDate)}</td>
                  <td className="py-3 px-4 text-sm">{formatTime(b.scheduledTime)}</td>
                  <td className="py-3 px-4 text-sm font-medium">
                    <Link href={`/admin/bookings/${b.id}`} className="text-blue-600 hover:text-blue-800">
                      {b.customer?.name ?? 'Unknown'}
                    </Link>
                  </td>
                  <td className="py-3 px-4 text-sm text-gray-600">
                    {b.vehicle ? `${b.vehicle.year} ${b.vehicle.make} ${b.vehicle.model}` : 'N/A'}
                  </td>
                  <td className="py-3 px-4 text-sm text-gray-600">
                    {b.services?.map(s => s.serviceName).join(', ') ?? 'N/A'}
                  </td>
                  <td className="py-3 px-4 text-sm">{formatPrice(b.totalAmount)}</td>
                  <td className="py-3 px-4"><StatusBadge status={b.status} /></td>
                  <td className="py-3 px-4"><StatusBadge status={b.paymentStatus} /></td>
                  <td className="py-3 px-4 text-sm">
                    <Link href={`/admin/bookings/${b.id}`} className="text-blue-600 hover:underline">
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
