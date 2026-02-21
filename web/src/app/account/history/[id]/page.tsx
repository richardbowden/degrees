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
        className="text-sm text-gray-500 hover:text-gray-700 mb-4 inline-block"
      >
        &larr; Back to History
      </Link>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Service Record</h1>

      <div className="space-y-6">
        <div className="border border-gray-200 rounded-lg p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-2">Completed</h2>
          <p className="text-sm text-gray-700">{formatDate(record.completedDate)}</p>
        </div>

        {visibleNotes.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Notes</h2>
            <div className="space-y-3">
              {visibleNotes.map(note => (
                <NoteCard key={note.id} note={note} />
              ))}
            </div>
          </div>
        )}

        {record.products && record.products.length > 0 && (
          <div>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Products Used</h2>
            <div className="border border-gray-200 rounded-lg overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="text-left px-4 py-3 font-medium text-gray-700">Product</th>
                    <th className="text-left px-4 py-3 font-medium text-gray-700">Notes</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {record.products.map(product => (
                    <tr key={product.id}>
                      <td className="px-4 py-3 text-gray-900 font-medium">{product.productName}</td>
                      <td className="px-4 py-3 text-gray-600">{product.notes || '-'}</td>
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
