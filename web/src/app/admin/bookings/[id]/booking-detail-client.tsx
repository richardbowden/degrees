'use client';

import { useState } from 'react';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Booking } from '@/lib/types';
import { formatPrice, formatDate, formatTime } from '@/lib/format';
import { StatusBadge } from '@/components/status-badge';

const BOOKING_STATUSES = ['pending', 'confirmed', 'in_progress', 'completed', 'cancelled'];

export function BookingDetailClient({ booking: initial, token }: { booking: Booking; token: string }) {
  const [booking, setBooking] = useState(initial);
  const [newStatus, setNewStatus] = useState(booking.status);
  const [statusLoading, setStatusLoading] = useState(false);
  const [statusError, setStatusError] = useState('');
  const [completeNotes, setCompleteNotes] = useState('');
  const [showComplete, setShowComplete] = useState(false);
  const [completeLoading, setCompleteLoading] = useState(false);
  const [completeError, setCompleteError] = useState('');

  async function handleStatusUpdate() {
    if (newStatus === booking.status) return;
    setStatusLoading(true);
    setStatusError('');
    try {
      const res = await api<{ booking: Booking }>(`/admin/bookings/${booking.id}/status`, {
        method: 'PUT',
        body: { status: newStatus },
        token,
      });
      setBooking(res.booking);
    } catch {
      setStatusError('Failed to update status');
    } finally {
      setStatusLoading(false);
    }
  }

  async function handleComplete() {
    setCompleteLoading(true);
    setCompleteError('');
    try {
      const res = await api<{ booking: Booking }>(`/admin/bookings/${booking.id}/complete`, {
        method: 'POST',
        body: { notes: completeNotes },
        token,
      });
      setBooking(res.booking);
      setShowComplete(false);
      setCompleteNotes('');
    } catch {
      setCompleteError('Failed to complete booking');
    } finally {
      setCompleteLoading(false);
    }
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link href="/admin/bookings" className="text-sm text-gray-500 hover:text-gray-700">
          &larr; All Bookings
        </Link>
        <h1 className="text-2xl font-bold text-gray-900">Booking Detail</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Schedule & Status */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Schedule</h2>
              <div className="flex gap-2">
                <StatusBadge status={booking.status} />
                <StatusBadge status={booking.paymentStatus} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-500">Date</p>
                <p className="font-medium">{formatDate(booking.scheduledDate)}</p>
              </div>
              <div>
                <p className="text-gray-500">Time</p>
                <p className="font-medium">{formatTime(booking.scheduledTime)}</p>
              </div>
              <div>
                <p className="text-gray-500">Duration</p>
                <p className="font-medium">{booking.estimatedDurationMins} mins</p>
              </div>
              <div>
                <p className="text-gray-500">Created</p>
                <p className="font-medium">{formatDate(booking.createdAt)}</p>
              </div>
            </div>
            {booking.notes && (
              <div className="mt-4 pt-4 border-t border-gray-100">
                <p className="text-gray-500 text-sm">Notes</p>
                <p className="text-sm mt-1 whitespace-pre-wrap">{booking.notes}</p>
              </div>
            )}
          </div>

          {/* Customer */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Customer</h2>
            {booking.customer ? (
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-gray-500">Name</p>
                  <p className="font-medium">
                    <Link href={`/admin/customers/${booking.customerId}`} className="text-blue-600 hover:text-blue-800">
                      {booking.customer.name}
                    </Link>
                  </p>
                </div>
                <div>
                  <p className="text-gray-500">Email</p>
                  <p className="font-medium">{booking.customer.email}</p>
                </div>
                <div>
                  <p className="text-gray-500">Phone</p>
                  <p className="font-medium">{booking.customer.phone || 'N/A'}</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-500">No customer info</p>
            )}
          </div>

          {/* Vehicle */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Vehicle</h2>
            {booking.vehicle ? (
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-gray-500">Vehicle</p>
                  <p className="font-medium">{booking.vehicle.year} {booking.vehicle.make} {booking.vehicle.model}</p>
                </div>
                <div>
                  <p className="text-gray-500">Colour</p>
                  <p className="font-medium">{booking.vehicle.colour || 'N/A'}</p>
                </div>
                <div>
                  <p className="text-gray-500">Rego</p>
                  <p className="font-medium">{booking.vehicle.rego || 'N/A'}</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-gray-500">No vehicle info</p>
            )}
          </div>

          {/* Services */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Services</h2>
            {booking.services && booking.services.length > 0 ? (
              <div className="space-y-3">
                {booking.services.map((s, i) => (
                  <div key={i} className="flex items-center justify-between text-sm">
                    <div>
                      <p className="font-medium">{s.serviceName}</p>
                      {s.options && s.options.length > 0 && (
                        <p className="text-gray-500 text-xs mt-0.5">Options: {s.options.join(', ')}</p>
                      )}
                    </div>
                    <p className="font-medium">{formatPrice(s.price)}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-gray-500">No services listed</p>
            )}
          </div>
        </div>

        {/* Sidebar: Pricing & Actions */}
        <div className="space-y-6">
          {/* Pricing */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Pricing</h2>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500">Subtotal</span>
                <span>{formatPrice(booking.subtotal)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">Deposit</span>
                <span>{formatPrice(booking.depositAmount)}</span>
              </div>
              <div className="flex justify-between font-semibold text-base pt-2 border-t border-gray-100">
                <span>Total</span>
                <span>{formatPrice(booking.totalAmount)}</span>
              </div>
            </div>
          </div>

          {/* Status Update */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
            <h2 className="text-lg font-semibold text-gray-900 mb-3">Update Status</h2>
            <div className="space-y-3">
              <select
                value={newStatus}
                onChange={e => setNewStatus(e.target.value)}
                className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
              >
                {BOOKING_STATUSES.map(s => (
                  <option key={s} value={s}>{s.replace('_', ' ')}</option>
                ))}
              </select>
              {statusError && <p className="text-sm text-red-600">{statusError}</p>}
              <button
                onClick={handleStatusUpdate}
                disabled={statusLoading || newStatus === booking.status}
                className="w-full bg-gray-900 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-50"
              >
                {statusLoading ? 'Updating...' : 'Update Status'}
              </button>
            </div>
          </div>

          {/* Complete Booking */}
          {booking.status !== 'completed' && booking.status !== 'cancelled' && (
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5">
              <h2 className="text-lg font-semibold text-gray-900 mb-3">Complete Booking</h2>
              {!showComplete ? (
                <button
                  onClick={() => setShowComplete(true)}
                  className="w-full bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-green-700"
                >
                  Mark as Complete
                </button>
              ) : (
                <div className="space-y-3">
                  <textarea
                    value={completeNotes}
                    onChange={e => setCompleteNotes(e.target.value)}
                    rows={3}
                    placeholder="Completion notes (optional)"
                    className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                  />
                  {completeError && <p className="text-sm text-red-600">{completeError}</p>}
                  <div className="flex gap-2">
                    <button
                      onClick={handleComplete}
                      disabled={completeLoading}
                      className="flex-1 bg-green-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-green-700 disabled:opacity-50"
                    >
                      {completeLoading ? 'Completing...' : 'Confirm'}
                    </button>
                    <button
                      onClick={() => setShowComplete(false)}
                      className="px-4 py-2 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
