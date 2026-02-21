import type { ServiceNote } from '@/lib/types';
import { formatDate } from '@/lib/format';

export function NoteCard({ note }: { note: ServiceNote }) {
  return (
    <div className="border border-gray-200 rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-medium text-gray-500 uppercase">{note.noteType}</span>
        <div className="flex items-center gap-2">
          {note.isVisibleToCustomer && (
            <span className="text-xs text-blue-600">Visible to customer</span>
          )}
          <span className="text-xs text-gray-400">{formatDate(note.createdAt)}</span>
        </div>
      </div>
      <p className="text-sm text-gray-700 whitespace-pre-wrap">{note.content}</p>
      {note.createdBy && (
        <p className="text-xs text-gray-400 mt-2">By: {note.createdBy}</p>
      )}
    </div>
  );
}
