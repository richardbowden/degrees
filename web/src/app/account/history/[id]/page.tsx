import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { ServiceRecord } from '@/lib/types';
import { formatDate } from '@/lib/format';
import { NoteCard } from '@/components/note-card';

export default async function HistoryDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const { record } = await api<{ record: ServiceRecord }>(`/me/history/${id}`, { token });

  const visibleNotes = (record.notes ?? []).filter(n => n.isVisibleToCustomer);

  return (
    <div>
      <Link
        href="/account/history"
        className="text-sm text-text-muted hover:text-text-secondary mb-4 inline-block"
      >
        &larr; Back to History
      </Link>

      <h1 className="text-2xl font-bold text-foreground mb-6">Service Record</h1>

      <div className="space-y-6">
        <div className="glass-card p-6">
          <h2 className="text-lg font-semibold text-foreground mb-2">Completed</h2>
          <p className="text-sm text-text-secondary">{formatDate(record.completedDate)}</p>
        </div>

        {visibleNotes.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-foreground mb-4">Notes</h2>
            <div className="space-y-3">
              {visibleNotes.map(note => (
                <NoteCard key={note.id} note={note} />
              ))}
            </div>
          </div>
        )}

        {record.products && record.products.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-foreground mb-4">Products Used</h2>
            <div className="border border-border-subtle rounded-lg overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-surface-input">
                  <tr>
                    <th className="text-left px-4 py-3 font-medium text-text-secondary">Product</th>
                    <th className="text-left px-4 py-3 font-medium text-text-secondary">Notes</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border-subtle">
                  {record.products.map(product => (
                    <tr key={product.id}>
                      <td className="px-4 py-3 text-foreground font-medium">{product.productName}</td>
                      <td className="px-4 py-3 text-text-secondary">{product.notes || '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
