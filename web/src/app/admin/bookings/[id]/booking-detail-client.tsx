'use client';

import { useState } from 'react';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import type { Booking, ServiceRecord, ServiceNote, ProductUsed } from '@/lib/types';
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

  // Service record state
  const [serviceRecord, setServiceRecord] = useState<ServiceRecord | null>(null);
  const [recordLoading, setRecordLoading] = useState(false);
  const [recordError, setRecordError] = useState('');

  // Add note form
  const [noteContent, setNoteContent] = useState('');
  const [noteType, setNoteType] = useState('general');
  const [noteVisible, setNoteVisible] = useState(false);
  const [noteLoading, setNoteLoading] = useState(false);

  // Add product form
  const [productName, setProductName] = useState('');
  const [productNotes, setProductNotes] = useState('');
  const [productLoading, setProductLoading] = useState(false);

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

  async function handleCreateRecord() {
    setRecordLoading(true);
    setRecordError('');
    try {
      const res = await api<{ record: ServiceRecord }>('/admin/records', {
        method: 'POST',
        body: {
          bookingId: booking.id,
          customerId: booking.customerId,
          vehicleId: booking.vehicleId,
        },
        token,
      });
      setServiceRecord(res.record);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setRecordError(apiErr.detail || 'Failed to create service record');
    } finally {
      setRecordLoading(false);
    }
  }

  async function handleAddNote() {
    if (!serviceRecord || !noteContent.trim()) return;
    setNoteLoading(true);
    setRecordError('');
    try {
      const res = await api<{ note: ServiceNote }>(`/admin/records/${serviceRecord.id}/notes`, {
        method: 'POST',
        body: {
          noteType,
          content: noteContent,
          isVisibleToCustomer: noteVisible,
        },
        token,
      });
      setServiceRecord(prev => prev ? {
        ...prev,
        notes: [...(prev.notes ?? []), res.note],
      } : prev);
      setNoteContent('');
      setNoteVisible(false);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setRecordError(apiErr.detail || 'Failed to add note');
    } finally {
      setNoteLoading(false);
    }
  }

  async function handleAddProduct() {
    if (!serviceRecord || !productName.trim()) return;
    setProductLoading(true);
    setRecordError('');
    try {
      const res = await api<{ product: ProductUsed }>(`/admin/records/${serviceRecord.id}/products`, {
        method: 'POST',
        body: {
          productName,
          notes: productNotes,
        },
        token,
      });
      setServiceRecord(prev => prev ? {
        ...prev,
        products: [...(prev.products ?? []), res.product],
      } : prev);
      setProductName('');
      setProductNotes('');
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setRecordError(apiErr.detail || 'Failed to add product');
    } finally {
      setProductLoading(false);
    }
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link href="/admin/bookings" className="text-sm text-text-muted hover:text-white">
          &larr; All Bookings
        </Link>
        <h1 className="text-2xl font-bold text-white">Booking Detail</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Schedule & Status */}
          <div className="glass-card p-5">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-white">Schedule</h2>
              <div className="flex gap-2">
                <StatusBadge status={booking.status} />
                <StatusBadge status={booking.paymentStatus} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-text-muted">Date</p>
                <p className="font-medium">{formatDate(booking.scheduledDate)}</p>
              </div>
              <div>
                <p className="text-text-muted">Time</p>
                <p className="font-medium">{formatTime(booking.scheduledTime)}</p>
              </div>
              <div>
                <p className="text-text-muted">Duration</p>
                <p className="font-medium">{booking.estimatedDurationMins} mins</p>
              </div>
              <div>
                <p className="text-text-muted">Created</p>
                <p className="font-medium">{formatDate(booking.createdAt)}</p>
              </div>
            </div>
            {booking.notes && (
              <div className="mt-4 pt-4 border-t border-white/5">
                <p className="text-text-muted text-sm">Notes</p>
                <p className="text-sm mt-1 whitespace-pre-wrap">{booking.notes}</p>
              </div>
            )}
          </div>

          {/* Customer */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Customer</h2>
            {booking.customer ? (
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-text-muted">Name</p>
                  <p className="font-medium">
                    <Link href={`/admin/customers/${booking.customerId}`} className="text-brand-400 hover:text-brand-500">
                      {booking.customer.name}
                    </Link>
                  </p>
                </div>
                <div>
                  <p className="text-text-muted">Email</p>
                  <p className="font-medium">{booking.customer.email}</p>
                </div>
                <div>
                  <p className="text-text-muted">Phone</p>
                  <p className="font-medium">{booking.customer.phone || 'N/A'}</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-text-muted">No customer info</p>
            )}
          </div>

          {/* Vehicle */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Vehicle</h2>
            {booking.vehicle ? (
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-text-muted">Vehicle</p>
                  <p className="font-medium">{booking.vehicle.year} {booking.vehicle.make} {booking.vehicle.model}</p>
                </div>
                <div>
                  <p className="text-text-muted">Colour</p>
                  <p className="font-medium">{booking.vehicle.colour || 'N/A'}</p>
                </div>
                <div>
                  <p className="text-text-muted">Rego</p>
                  <p className="font-medium">{booking.vehicle.rego || 'N/A'}</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-text-muted">No vehicle info</p>
            )}
          </div>

          {/* Services */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Services</h2>
            {booking.services && booking.services.length > 0 ? (
              <div className="space-y-3">
                {booking.services.map((s, i) => (
                  <div key={i} className="flex items-center justify-between text-sm">
                    <div>
                      <p className="font-medium">{s.serviceName}</p>
                      {s.options && s.options.length > 0 && (
                        <p className="text-text-muted text-xs mt-0.5">Options: {s.options.join(', ')}</p>
                      )}
                    </div>
                    <p className="font-medium">{formatPrice(s.price)}</p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-text-muted">No services listed</p>
            )}
          </div>

          {/* Service Record (shown after booking is completed) */}
          {booking.status === 'completed' && (
            <div className="glass-card p-5">
              <h2 className="text-lg font-semibold text-white mb-3">Service Record</h2>
              {recordError && <p className="text-sm text-red-400 mb-3">{recordError}</p>}

              {!serviceRecord ? (
                <div>
                  <p className="text-sm text-text-muted mb-3">
                    Create a service record to log notes and products used.
                  </p>
                  <button
                    onClick={handleCreateRecord}
                    disabled={recordLoading}
                    className="btn-brand px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
                  >
                    {recordLoading ? 'Creating...' : 'Create Service Record'}
                  </button>
                </div>
              ) : (
                <div className="space-y-6">
                  {/* Notes list */}
                  {serviceRecord.notes && serviceRecord.notes.length > 0 && (
                    <div>
                      <h3 className="text-sm font-medium text-text-secondary mb-2">Notes</h3>
                      <div className="space-y-2">
                        {serviceRecord.notes.map(note => (
                          <div key={note.id} className="border border-white/5 rounded p-3 text-sm">
                            <div className="flex items-center gap-2 mb-1">
                              <span className="font-medium text-white capitalize">{note.noteType}</span>
                              {note.isVisibleToCustomer ? (
                                <span className="text-xs bg-blue-500/20 text-blue-400 px-1.5 py-0.5 rounded">Customer visible</span>
                              ) : (
                                <span className="text-xs bg-white/10 text-text-muted px-1.5 py-0.5 rounded">Internal</span>
                              )}
                            </div>
                            <p className="text-text-secondary whitespace-pre-wrap">{note.content}</p>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Add note form */}
                  <div className="border-t border-white/5 pt-4">
                    <h3 className="text-sm font-medium text-text-secondary mb-2">Add Note</h3>
                    <div className="space-y-2">
                      <div className="flex gap-2">
                        <select
                          value={noteType}
                          onChange={e => setNoteType(e.target.value)}
                          className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                        >
                          <option value="general">General</option>
                          <option value="condition">Condition</option>
                          <option value="recommendation">Recommendation</option>
                        </select>
                        <label className="flex items-center gap-1.5 text-sm text-text-secondary">
                          <input
                            type="checkbox"
                            checked={noteVisible}
                            onChange={e => setNoteVisible(e.target.checked)}
                            className="rounded border-border-subtle"
                          />
                          Visible to customer
                        </label>
                      </div>
                      <textarea
                        value={noteContent}
                        onChange={e => setNoteContent(e.target.value)}
                        rows={2}
                        placeholder="Enter note..."
                        className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-2 text-sm text-white"
                      />
                      <button
                        onClick={handleAddNote}
                        disabled={noteLoading || !noteContent.trim()}
                        className="btn-brand px-4 py-1.5 rounded-md text-sm font-medium disabled:opacity-50"
                      >
                        {noteLoading ? 'Adding...' : 'Add Note'}
                      </button>
                    </div>
                  </div>

                  {/* Products list */}
                  {serviceRecord.products && serviceRecord.products.length > 0 && (
                    <div className="border-t border-white/5 pt-4">
                      <h3 className="text-sm font-medium text-text-secondary mb-2">Products Used</h3>
                      <div className="space-y-2">
                        {serviceRecord.products.map(product => (
                          <div key={product.id} className="flex items-center justify-between text-sm border border-white/5 rounded p-3">
                            <span className="font-medium text-white">{product.productName}</span>
                            {product.notes && <span className="text-text-muted">{product.notes}</span>}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Add product form */}
                  <div className="border-t border-white/5 pt-4">
                    <h3 className="text-sm font-medium text-text-secondary mb-2">Log Product Used</h3>
                    <div className="space-y-2">
                      <div className="grid grid-cols-2 gap-2">
                        <input
                          type="text"
                          value={productName}
                          onChange={e => setProductName(e.target.value)}
                          placeholder="Product name (e.g. Bowden's Nanolicious)"
                          className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                        />
                        <input
                          type="text"
                          value={productNotes}
                          onChange={e => setProductNotes(e.target.value)}
                          placeholder="Notes (optional)"
                          className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                        />
                      </div>
                      <button
                        onClick={handleAddProduct}
                        disabled={productLoading || !productName.trim()}
                        className="btn-brand px-4 py-1.5 rounded-md text-sm font-medium disabled:opacity-50"
                      >
                        {productLoading ? 'Adding...' : 'Add Product'}
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Sidebar: Pricing & Actions */}
        <div className="space-y-6">
          {/* Pricing */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Pricing</h2>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-text-muted">Subtotal</span>
                <span>{formatPrice(booking.subtotal)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-text-muted">Deposit</span>
                <span>{formatPrice(booking.depositAmount)}</span>
              </div>
              <div className="flex justify-between font-semibold text-base pt-2 border-t border-white/5">
                <span>Total</span>
                <span>{formatPrice(booking.totalAmount)}</span>
              </div>
            </div>
          </div>

          {/* Status Update */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Update Status</h2>
            <div className="space-y-3">
              <select
                value={newStatus}
                onChange={e => setNewStatus(e.target.value)}
                className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
              >
                {BOOKING_STATUSES.map(s => (
                  <option key={s} value={s}>{s.replace('_', ' ')}</option>
                ))}
              </select>
              {statusError && <p className="text-sm text-red-400">{statusError}</p>}
              <button
                onClick={handleStatusUpdate}
                disabled={statusLoading || newStatus === booking.status}
                className="w-full btn-brand px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
              >
                {statusLoading ? 'Updating...' : 'Update Status'}
              </button>
            </div>
          </div>

          {/* Complete Booking */}
          {booking.status !== 'completed' && booking.status !== 'cancelled' && (
            <div className="glass-card p-5">
              <h2 className="text-lg font-semibold text-white mb-3">Complete Booking</h2>
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
                    className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-2 text-sm text-white"
                  />
                  {completeError && <p className="text-sm text-red-400">{completeError}</p>}
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
                      className="px-4 py-2 border border-border-subtle rounded-md text-sm text-text-secondary hover:bg-white/5"
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
