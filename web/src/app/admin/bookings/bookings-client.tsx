'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { formatPrice, formatTime, formatDate } from '@/lib/format';
import { StatusBadge } from '@/components/status-badge';

type Preset = 'all' | 'today' | 'week' | 'month' | 'custom';

const STATUSES = ['all', 'pending', 'confirmed', 'in_progress', 'completed', 'cancelled'];

function todayStr() {
  return new Date().toISOString().split('T')[0];
}

function presetRange(preset: Preset): { from: string; to: string } {
  const now = new Date();
  const pad = (n: number) => String(n).padStart(2, '0');
  const fmt = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;

  if (preset === 'today') {
    const t = fmt(now);
    return { from: t, to: t };
  }
  if (preset === 'week') {
    const day = now.getDay();
    const mon = new Date(now);
    mon.setDate(now.getDate() - ((day + 6) % 7));
    const sun = new Date(mon);
    sun.setDate(mon.getDate() + 6);
    return { from: fmt(mon), to: fmt(sun) };
  }
  if (preset === 'month') {
    const first = new Date(now.getFullYear(), now.getMonth(), 1);
    const last = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    return { from: fmt(first), to: fmt(last) };
  }
  // 'all' or 'custom' — empty means all
  return { from: '', to: '' };
}

export function BookingsClient({ token }: { token: string }) {
  const [preset, setPreset] = useState<Preset>('all');
  const [customFrom, setCustomFrom] = useState(todayStr());
  const [customTo, setCustomTo] = useState(todayStr());
  const [statusFilter, setStatusFilter] = useState('all');
  const [search, setSearch] = useState('');
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchBookings = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const range = preset === 'custom' ? { from: customFrom, to: customTo } : presetRange(preset);
      const params = new URLSearchParams();
      if (range.from) params.set('dateFrom', range.from);
      if (range.to) params.set('dateTo', range.to);
      const qs = params.toString();
      const res = await api<{ bookings: Booking[] }>(
        `/admin/bookings${qs ? `?${qs}` : ''}`,
        { token },
      );
      setBookings(res.bookings ?? []);
    } catch {
      setError('Failed to load bookings');
    } finally {
      setLoading(false);
    }
  }, [preset, customFrom, customTo, token]);

  useEffect(() => {
    fetchBookings();
  }, [fetchBookings]);

  const filtered = bookings.filter(b => {
    if (statusFilter !== 'all' && b.status !== statusFilter) return false;
    if (search.trim()) {
      const q = search.toLowerCase();
      const name = b.customer?.name?.toLowerCase() ?? '';
      const rego = b.vehicle?.rego?.toLowerCase() ?? '';
      const svc = b.services?.map(s => s.serviceName).join(' ').toLowerCase() ?? '';
      if (!name.includes(q) && !rego.includes(q) && !svc.includes(q)) return false;
    }
    return true;
  });

  const PRESETS: { id: Preset; label: string }[] = [
    { id: 'all', label: 'All' },
    { id: 'today', label: 'Today' },
    { id: 'week', label: 'This Week' },
    { id: 'month', label: 'This Month' },
    { id: 'custom', label: 'Custom' },
  ];

  return (
    <div>
      {/* Preset tabs */}
      <div className="flex flex-wrap items-center gap-2 mb-4">
        {PRESETS.map(p => (
          <button
            key={p.id}
            onClick={() => setPreset(p.id)}
            className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
              preset === p.id
                ? 'bg-brand-500 text-white'
                : 'bg-surface-input border border-border-subtle text-text-secondary hover:text-foreground hover:bg-surface-hover'
            }`}
          >
            {p.label}
          </button>
        ))}
      </div>

      {/* Custom date inputs */}
      {preset === 'custom' && (
        <div className="flex flex-wrap gap-3 items-end mb-4">
          <div>
            <label className="block text-xs font-medium text-text-muted mb-1">From</label>
            <input
              type="date"
              value={customFrom}
              onChange={e => setCustomFrom(e.target.value)}
              className="bg-surface-input border border-border-subtle rounded-md px-3 py-1.5 text-sm text-foreground"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-text-muted mb-1">To</label>
            <input
              type="date"
              value={customTo}
              onChange={e => setCustomTo(e.target.value)}
              className="bg-surface-input border border-border-subtle rounded-md px-3 py-1.5 text-sm text-foreground"
            />
          </div>
        </div>
      )}

      {/* Filters row */}
      <div className="flex flex-wrap gap-3 items-center mb-6">
        <input
          type="text"
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Search customer, rego, service..."
          className="bg-surface-input border border-border-subtle rounded-md px-3 py-1.5 text-sm text-foreground w-64"
        />
        <select
          value={statusFilter}
          onChange={e => setStatusFilter(e.target.value)}
          className="bg-surface-input border border-border-subtle rounded-md px-3 py-1.5 text-sm text-foreground"
        >
          {STATUSES.map(s => (
            <option key={s} value={s}>{s === 'all' ? 'All Statuses' : s.replace('_', ' ')}</option>
          ))}
        </select>
        {!loading && (
          <span className="text-sm text-text-muted ml-auto">
            {filtered.length} booking{filtered.length !== 1 ? 's' : ''}
          </span>
        )}
      </div>

      {error && <p className="text-red-400 text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-sm text-text-muted">Loading bookings...</p>
      ) : filtered.length === 0 ? (
        <p className="text-sm text-text-muted">No bookings found.</p>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-surface-input border-b border-border-subtle">
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
                <tr key={b.id} className="border-b border-border-subtle hover:bg-surface-hover">
                  <td className="py-3 px-4 text-sm">{formatDate(b.scheduledDate)}</td>
                  <td className="py-3 px-4 text-sm">{formatTime(b.scheduledTime)}</td>
                  <td className="py-3 px-4 text-sm font-medium">
                    <Link href={`/admin/bookings/${b.id}`} className="text-brand-400 hover:text-brand-500">
                      {b.customer?.name || 'Unknown'}
                    </Link>
                  </td>
                  <td className="py-3 px-4 text-sm text-text-secondary">
                    {b.vehicle
                      ? `${b.vehicle.make} ${b.vehicle.model}${b.vehicle.rego ? ` (${b.vehicle.rego})` : ''}`
                      : 'N/A'}
                  </td>
                  <td className="py-3 px-4 text-sm text-text-secondary">
                    {b.services?.map(s => s.serviceName).join(', ') || 'N/A'}
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
