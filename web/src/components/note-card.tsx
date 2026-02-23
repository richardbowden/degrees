import type { ServiceNote } from '@/lib/types';
import { formatDate } from '@/lib/format';

export function NoteCard({ note }: { note: ServiceNote }) {
  return (
    <div className="border border-border-subtle rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-medium text-text-muted uppercase">{note.noteType}</span>
        <div className="flex items-center gap-2">
          {note.isVisibleToCustomer && (
            <span className="text-xs text-brand-400">Visible to customer</span>
          )}
          <span className="text-xs text-text-muted">{formatDate(note.createdAt)}</span>
        </div>
      </div>
      <p className="text-sm text-text-secondary whitespace-pre-wrap">{note.content}</p>
      {note.createdBy && (
        <p className="text-xs text-text-muted mt-2">By: {note.createdBy}</p>
      )}
    </div>
  );
}
