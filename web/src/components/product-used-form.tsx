'use client';

import { useState } from 'react';
import { api } from '@/lib/api';
import type { ProductUsed } from '@/lib/types';

interface ProductUsedFormProps {
  recordId: string;
  token: string;
  onAdded: (product: ProductUsed) => void;
}

export function ProductUsedForm({ recordId, token, onAdded }: ProductUsedFormProps) {
  const [productName, setProductName] = useState('');
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!productName.trim()) return;
    setLoading(true);
    setError('');
    try {
      const res = await api<{ product: ProductUsed }>(`/admin/records/${recordId}/products`, {
        method: 'POST',
        body: { productName: productName.trim(), notes: notes.trim() },
        token,
      });
      onAdded(res.product);
      setProductName('');
      setNotes('');
    } catch {
      setError('Failed to add product');
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="flex gap-3">
        <input
          type="text"
          value={productName}
          onChange={e => setProductName(e.target.value)}
          placeholder="Product name (e.g. Bowden's Auto Body Gel)"
          className="flex-1 border border-gray-300 rounded-md px-3 py-1.5 text-sm"
        />
      </div>
      <input
        type="text"
        value={notes}
        onChange={e => setNotes(e.target.value)}
        placeholder="Usage notes (optional)"
        className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
      />
      {error && <p className="text-sm text-red-600">{error}</p>}
      <button
        type="submit"
        disabled={loading || !productName.trim()}
        className="bg-gray-900 text-white px-4 py-1.5 rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-50"
      >
        {loading ? 'Adding...' : 'Add Product'}
      </button>
    </form>
  );
}
