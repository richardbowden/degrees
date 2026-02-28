'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { CustomerProfile } from '@/lib/types';

const PAGE_SIZE = 50;

export function CustomersClient({ token }: { token: string }) {
  const [customers, setCustomers] = useState<CustomerProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState('');
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [search, setSearch] = useState('');

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError('');
      try {
        const res = await api<{ customers: CustomerProfile[] }>(
          `/admin/customers?limit=${PAGE_SIZE}&offset=0`,
          { token },
        );
        const list = res.customers ?? [];
        setCustomers(list);
        setOffset(PAGE_SIZE);
        setHasMore(list.length === PAGE_SIZE);
      } catch (err: unknown) {
        const apiErr = err as { status?: number; detail?: string; message?: string };
        setError(apiErr?.detail || apiErr?.message || 'Failed to load customers');
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [token]);

  async function loadMore() {
    setLoadingMore(true);
    try {
      const res = await api<{ customers: CustomerProfile[] }>(
        `/admin/customers?limit=${PAGE_SIZE}&offset=${offset}`,
        { token },
      );
      const list = res.customers ?? [];
      setCustomers(prev => [...prev, ...list]);
      setOffset(prev => prev + PAGE_SIZE);
      setHasMore(list.length === PAGE_SIZE);
    } catch (err: unknown) {
      const apiErr = err as { status?: number; detail?: string; message?: string };
      setError(apiErr?.detail || apiErr?.message || 'Failed to load more customers');
    } finally {
      setLoadingMore(false);
    }
  }

  const filtered = search.trim()
    ? customers.filter(c => {
        const q = search.toLowerCase();
        return (
          c.phone?.toLowerCase().includes(q) ||
          c.suburb?.toLowerCase().includes(q) ||
          c.address?.toLowerCase().includes(q) ||
          c.userId?.toLowerCase().includes(q)
        );
      })
    : customers;

  if (loading) return <p className="text-sm text-text-muted">Loading customers...</p>;
  if (error && customers.length === 0) return <p className="text-red-400 text-sm">{error}</p>;

  return (
    <div>
      <div className="mb-4">
        <input
          type="text"
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Search by phone, suburb, address..."
          className="w-full max-w-md bg-surface-input border border-border-subtle rounded-md px-3 py-2 text-sm text-foreground"
        />
      </div>

      {filtered.length === 0 ? (
        <p className="text-sm text-text-muted">No customers found.</p>
      ) : (
        <div className="glass-card overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-surface-input border-b border-border-subtle">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Phone</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Address</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Suburb</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Postcode</th>
                <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase"></th>
              </tr>
            </thead>
            <tbody>
              {filtered.map(c => (
                <tr key={c.id} className="border-b border-border-subtle hover:bg-surface-hover">
                  <td className="py-3 px-4 text-sm">{c.phone || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-text-secondary">{c.address || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-text-secondary">{c.suburb || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-text-secondary">{c.postcode || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm">
                    <Link href={`/admin/customers/${c.id}`} className="text-brand-400 hover:underline">
                      View
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {hasMore && !search.trim() && (
        <div className="mt-4">
          <button
            onClick={loadMore}
            disabled={loadingMore}
            className="bg-surface-input border border-border-subtle text-text-secondary px-4 py-2 rounded-md text-sm hover:bg-surface-hover disabled:opacity-50"
          >
            {loadingMore ? 'Loading...' : 'Load More'}
          </button>
        </div>
      )}
    </div>
  );
}
