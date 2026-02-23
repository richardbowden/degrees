'use client';

import { useState } from 'react';
import Link from 'next/link';
import type { CustomerProfile, Vehicle, ServiceRecord, ServiceNote, ProductUsed } from '@/lib/types';
import { formatDate } from '@/lib/format';
import { NoteCard } from '@/components/note-card';
import { AddNoteForm } from '@/components/add-note-form';
import { ProductUsedForm } from '@/components/product-used-form';

interface Props {
  profile: CustomerProfile;
  vehicles: Vehicle[];
  initialRecords: ServiceRecord[];
  token: string;
}

export function CustomerDetailClient({ profile, vehicles, initialRecords, token }: Props) {
  const [records, setRecords] = useState(initialRecords);
  const [expandedRecord, setExpandedRecord] = useState<string | null>(null);

  function handleNoteAdded(recordId: string, note: ServiceNote) {
    setRecords(prev =>
      prev.map(r =>
        r.id === recordId ? { ...r, notes: [...(r.notes ?? []), note] } : r,
      ),
    );
  }

  function handleProductAdded(recordId: string, product: ProductUsed) {
    setRecords(prev =>
      prev.map(r =>
        r.id === recordId ? { ...r, products: [...(r.products ?? []), product] } : r,
      ),
    );
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link href="/admin/customers" className="text-sm text-text-muted hover:text-white">
          &larr; All Customers
        </Link>
        <h1 className="text-2xl font-bold text-white">Customer Detail</h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile */}
        <div className="lg:col-span-1 space-y-6">
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Profile</h2>
            <div className="space-y-3 text-sm">
              <div>
                <p className="text-text-muted">Phone</p>
                <p className="font-medium">{profile.phone || 'N/A'}</p>
              </div>
              <div>
                <p className="text-text-muted">Address</p>
                <p className="font-medium">{profile.address || 'N/A'}</p>
              </div>
              <div>
                <p className="text-text-muted">Suburb</p>
                <p className="font-medium">{profile.suburb || 'N/A'}</p>
              </div>
              <div>
                <p className="text-text-muted">Postcode</p>
                <p className="font-medium">{profile.postcode || 'N/A'}</p>
              </div>
              {profile.notes && (
                <div>
                  <p className="text-text-muted">Notes</p>
                  <p className="font-medium whitespace-pre-wrap">{profile.notes}</p>
                </div>
              )}
              <div>
                <p className="text-text-muted">Customer since</p>
                <p className="font-medium">{formatDate(profile.createdAt)}</p>
              </div>
            </div>
          </div>

          {/* Vehicles */}
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-3">Vehicles</h2>
            {vehicles.length === 0 ? (
              <p className="text-sm text-text-muted">No vehicles on file.</p>
            ) : (
              <div className="space-y-4">
                {vehicles.map(v => (
                  <div key={v.id} className="border border-white/5 rounded-lg p-3 text-sm">
                    <p className="font-medium">
                      {v.year} {v.make} {v.model}
                      {v.isPrimary && (
                        <span className="ml-2 text-xs bg-blue-500/20 text-blue-400 px-1.5 py-0.5 rounded">Primary</span>
                      )}
                    </p>
                    <div className="mt-1 text-text-muted space-y-0.5">
                      {v.colour && <p>Colour: {v.colour}</p>}
                      {v.rego && <p>Rego: {v.rego}</p>}
                      {v.paintType && <p>Paint: {v.paintType}</p>}
                      {v.conditionNotes && <p>Condition: {v.conditionNotes}</p>}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Service History */}
        <div className="lg:col-span-2">
          <div className="glass-card p-5">
            <h2 className="text-lg font-semibold text-white mb-4">Service History</h2>
            {records.length === 0 ? (
              <p className="text-sm text-text-muted">No service records yet.</p>
            ) : (
              <div className="space-y-4">
                {records.map(record => {
                  const isExpanded = expandedRecord === record.id;
                  return (
                    <div key={record.id} className="border border-border-subtle rounded-lg">
                      <button
                        onClick={() => setExpandedRecord(isExpanded ? null : record.id)}
                        className="w-full flex items-center justify-between p-4 text-left hover:bg-white/5"
                      >
                        <div className="text-sm">
                          <p className="font-medium">
                            Service on {formatDate(record.completedDate)}
                          </p>
                          <p className="text-text-muted mt-0.5">
                            {(record.notes?.length ?? 0)} notes, {(record.products?.length ?? 0)} products
                          </p>
                        </div>
                        <span className="text-text-muted text-sm">{isExpanded ? 'Collapse' : 'Expand'}</span>
                      </button>

                      {isExpanded && (
                        <div className="border-t border-border-subtle p-4 space-y-6">
                          {/* Notes */}
                          <div>
                            <h3 className="text-sm font-semibold text-white mb-3">Notes</h3>
                            {record.notes && record.notes.length > 0 ? (
                              <div className="space-y-2 mb-4">
                                {record.notes.map(n => <NoteCard key={n.id} note={n} />)}
                              </div>
                            ) : (
                              <p className="text-sm text-text-muted mb-4">No notes yet.</p>
                            )}
                            <AddNoteForm
                              recordId={record.id}
                              token={token}
                              onAdded={(note) => handleNoteAdded(record.id, note)}
                            />
                          </div>

                          {/* Products */}
                          <div>
                            <h3 className="text-sm font-semibold text-white mb-3">Products Used</h3>
                            {record.products && record.products.length > 0 ? (
                              <div className="space-y-2 mb-4">
                                {record.products.map(p => (
                                  <div key={p.id} className="border border-white/5 rounded-lg p-3 text-sm">
                                    <p className="font-medium">{p.productName}</p>
                                    {p.notes && <p className="text-text-muted mt-0.5">{p.notes}</p>}
                                  </div>
                                ))}
                              </div>
                            ) : (
                              <p className="text-sm text-text-muted mb-4">No products recorded.</p>
                            )}
                            <ProductUsedForm
                              recordId={record.id}
                              token={token}
                              onAdded={(product) => handleProductAdded(record.id, product)}
                            />
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
