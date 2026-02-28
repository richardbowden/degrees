'use client';

import { useState } from 'react';
import { api } from '@/lib/api';
import type { ServiceNote } from '@/lib/types';

interface AddNoteFormProps {
  recordId: string;
  token: string;
  onAdded: (note: ServiceNote) => void;
}

export function AddNoteForm({ recordId, token, onAdded }: AddNoteFormProps) {
  const [content, setContent] = useState('');
  const [noteType, setNoteType] = useState('observation');
  const [isVisible, setIsVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!content.trim()) return;
    setLoading(true);
    setError('');
    try {
      const res = await api<{ note: ServiceNote }>(`/admin/records/${recordId}/notes`, {
        method: 'POST',
        body: { noteType, content: content.trim(), isVisibleToCustomer: isVisible },
        token,
      });
      onAdded(res.note);
      setContent('');
    } catch {
      setError('Failed to add note');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="flex gap-3">
        <select
          value={noteType}
          onChange={e => setNoteType(e.target.value)}
          className="bg-surface-input border border-border-subtle rounded-md px-3 py-1.5 text-sm text-foreground"
        >
          <option value="observation">Observation</option>
          <option value="recommendation">Recommendation</option>
          <option value="damage">Damage</option>
          <option value="follow_up">Follow Up</option>
        </select>
        <label className="flex items-center gap-1.5 text-sm text-text-secondary">
          <input
            type="checkbox"
            checked={isVisible}
            onChange={e => setIsVisible(e.target.checked)}
            className="rounded border-border-subtle"
          />
          Visible to customer
        </label>
      </div>
      <textarea
        value={content}
        onChange={e => setContent(e.target.value)}
        rows={3}
        placeholder="Add a note..."
        className="w-full bg-surface-input border border-border-subtle rounded-md px-3 py-2 text-sm text-foreground"
      />
      {error && <p className="text-sm text-red-400">{error}</p>}
      <button
        type="submit"
        disabled={loading || !content.trim()}
        className="btn-brand px-4 py-1.5 rounded-md text-sm font-medium disabled:opacity-50"
      >
        {loading ? 'Adding...' : 'Add Note'}
      </button>
    </form>
  );
}
