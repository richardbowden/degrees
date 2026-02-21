'use client';

import { useState, useEffect, useCallback } from 'react';
import { api, type ApiError } from '@/lib/api';
import type { ScheduleDay, Blackout } from '@/lib/types';

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export function ScheduleClient({ token }: { token: string }) {
  const [days, setDays] = useState<ScheduleDay[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [savingDay, setSavingDay] = useState<number | null>(null);
  const [dayErrors, setDayErrors] = useState<Record<number, string>>({});

  // Blackout form
  const [blackoutDate, setBlackoutDate] = useState('');
  const [blackoutReason, setBlackoutReason] = useState('');
  const [blackoutLoading, setBlackoutLoading] = useState(false);
  const [blackoutError, setBlackoutError] = useState('');
  const [blackoutSuccess, setBlackoutSuccess] = useState('');

  // Blackout list
  const [blackouts, setBlackouts] = useState<Blackout[]>([]);
  const [deletingBlackout, setDeletingBlackout] = useState<string | null>(null);

  const fetchConfig = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await api<{ days: ScheduleDay[] }>('/admin/schedule/config', { token });
      const sorted = (res.days ?? []).sort((a, b) => Number(a.dayOfWeek) - Number(b.dayOfWeek));
      setDays(sorted);
    } catch {
      setError('Failed to load schedule config');
    } finally {
      setLoading(false);
    }
  }, [token]);

  const fetchBlackouts = useCallback(async () => {
    try {
      const res = await api<{ blackouts: Blackout[] }>('/admin/schedule/blackout', { token });
      setBlackouts(res.blackouts ?? []);
    } catch {
      // Endpoint may not exist; ignore
    }
  }, [token]);

  useEffect(() => {
    fetchConfig();
    fetchBlackouts();
  }, [fetchConfig, fetchBlackouts]);

  function updateDay(index: number, field: keyof ScheduleDay, value: string | boolean | number) {
    setDays(prev => prev.map((d, i) => i === index ? { ...d, [field]: value } : d));
  }

  async function saveDay(day: ScheduleDay, index: number) {
    setSavingDay(day.dayOfWeek);
    setDayErrors(prev => ({ ...prev, [index]: '' }));
    try {
      const res = await api<{ day: ScheduleDay }>('/admin/schedule/config', {
        method: 'PUT',
        body: {
          dayOfWeek: day.dayOfWeek,
          openTime: day.openTime,
          closeTime: day.closeTime,
          isOpen: day.isOpen,
          bufferMinutes: day.bufferMinutes,
        },
        token,
      });
      setDays(prev => prev.map((d, i) => i === index ? res.day : d));
    } catch (err) {
      setDayErrors(prev => ({ ...prev, [index]: (err as ApiError)?.detail || 'Failed to save' }));
    } finally {
      setSavingDay(null);
    }
  }

  async function handleAddBlackout(e: React.FormEvent) {
    e.preventDefault();
    if (!blackoutDate) return;
    setBlackoutLoading(true);
    setBlackoutError('');
    setBlackoutSuccess('');
    try {
      const res = await api<{ blackout: Blackout }>('/admin/schedule/blackout', {
        method: 'POST',
        body: { date: blackoutDate, reason: blackoutReason },
        token,
      });
      setBlackouts(prev => [...prev, res.blackout]);
      setBlackoutDate('');
      setBlackoutReason('');
      setBlackoutSuccess('Blackout date added');
      setTimeout(() => setBlackoutSuccess(''), 3000);
    } catch (err) {
      setBlackoutError((err as ApiError)?.detail || 'Failed to add blackout date');
    } finally {
      setBlackoutLoading(false);
    }
  }

  async function handleDeleteBlackout(id: string) {
    setDeletingBlackout(id);
    try {
      await api(`/admin/schedule/blackout/${id}`, { method: 'DELETE', token });
      setBlackouts(prev => prev.filter(b => b.id !== id));
    } catch {
      alert('Failed to delete blackout date');
    } finally {
      setDeletingBlackout(null);
    }
  }

  if (loading) return <p className="text-sm text-gray-500">Loading schedule...</p>;
  if (error && days.length === 0) return <p className="text-red-600 text-sm">{error}</p>;

  return (
    <div className="space-y-8">
      {/* Business Hours */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Business Hours</h2>
        <div className="space-y-3">
          {days.map((day, index) => (
            <div key={day.id || day.dayOfWeek} className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
              <div className="flex flex-wrap items-center gap-4">
                <div className="w-28">
                  <p className="font-medium text-sm text-gray-900">{DAY_NAMES[day.dayOfWeek]}</p>
                </div>

                <label className="flex items-center gap-1.5 text-sm text-gray-700">
                  <input
                    type="checkbox"
                    checked={day.isOpen}
                    onChange={e => updateDay(index, 'isOpen', e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  Open
                </label>

                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500">From</label>
                  <input
                    type="time"
                    value={day.openTime}
                    onChange={e => updateDay(index, 'openTime', e.target.value)}
                    disabled={!day.isOpen}
                    className="border border-gray-300 rounded-md px-2 py-1 text-sm disabled:opacity-50"
                  />
                </div>

                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500">To</label>
                  <input
                    type="time"
                    value={day.closeTime}
                    onChange={e => updateDay(index, 'closeTime', e.target.value)}
                    disabled={!day.isOpen}
                    className="border border-gray-300 rounded-md px-2 py-1 text-sm disabled:opacity-50"
                  />
                </div>

                <div className="flex items-center gap-2">
                  <label className="text-xs text-gray-500">Buffer (mins)</label>
                  <input
                    type="number"
                    min="0"
                    value={day.bufferMinutes}
                    onChange={e => updateDay(index, 'bufferMinutes', parseInt(e.target.value || '0', 10))}
                    disabled={!day.isOpen}
                    className="w-20 border border-gray-300 rounded-md px-2 py-1 text-sm disabled:opacity-50"
                  />
                </div>

                <button
                  onClick={() => saveDay(day, index)}
                  disabled={savingDay === day.dayOfWeek}
                  className="bg-gray-900 text-white px-3 py-1 rounded-md text-sm hover:bg-gray-800 disabled:opacity-50"
                >
                  {savingDay === day.dayOfWeek ? 'Saving...' : 'Save'}
                </button>
              </div>
              {dayErrors[index] && <p className="text-sm text-red-600 mt-2">{dayErrors[index]}</p>}
            </div>
          ))}
        </div>
      </div>

      {/* Blackout Dates */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Blackout Dates</h2>

        {/* Existing blackouts */}
        {blackouts.length > 0 && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 mb-4 overflow-hidden">
            <table className="w-full text-left">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="py-2 px-4 text-xs font-medium text-gray-500 uppercase">Date</th>
                  <th className="py-2 px-4 text-xs font-medium text-gray-500 uppercase">Reason</th>
                  <th className="py-2 px-4 text-xs font-medium text-gray-500 uppercase"></th>
                </tr>
              </thead>
              <tbody>
                {blackouts.map(b => (
                  <tr key={b.id} className="border-b border-gray-100">
                    <td className="py-2 px-4 text-sm">{b.date}</td>
                    <td className="py-2 px-4 text-sm text-gray-600">{b.reason || '-'}</td>
                    <td className="py-2 px-4 text-sm">
                      <button
                        onClick={() => handleDeleteBlackout(b.id)}
                        disabled={deletingBlackout === b.id}
                        className="text-red-600 hover:text-red-800 text-xs disabled:opacity-50"
                      >
                        {deletingBlackout === b.id ? '...' : 'Remove'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
          <h3 className="text-sm font-semibold text-gray-900 mb-3">Add Blackout Date</h3>
          <form onSubmit={handleAddBlackout} className="flex flex-wrap items-end gap-3">
            <div>
              <label className="block text-xs font-medium text-gray-500 mb-1">Date</label>
              <input
                type="date"
                value={blackoutDate}
                onChange={e => setBlackoutDate(e.target.value)}
                className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                required
              />
            </div>
            <div className="flex-1 min-w-48">
              <label className="block text-xs font-medium text-gray-500 mb-1">Reason</label>
              <input
                type="text"
                value={blackoutReason}
                onChange={e => setBlackoutReason(e.target.value)}
                placeholder="e.g. Public holiday"
                className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
              />
            </div>
            <button
              type="submit"
              disabled={blackoutLoading || !blackoutDate}
              className="bg-gray-900 text-white px-4 py-1.5 rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-50"
            >
              {blackoutLoading ? 'Adding...' : 'Add Blackout'}
            </button>
          </form>
          {blackoutError && <p className="text-sm text-red-600 mt-2">{blackoutError}</p>}
          {blackoutSuccess && <p className="text-sm text-green-600 mt-2">{blackoutSuccess}</p>}
        </div>
      </div>
    </div>
  );
}
