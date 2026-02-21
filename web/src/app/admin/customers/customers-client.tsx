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
      } catch {
        setError('Failed to load customers');
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
    } catch {
      setError('Failed to load more customers');
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

  if (loading) return <p className="text-sm text-gray-500">Loading customers...</p>;
  if (error && customers.length === 0) return <p className="text-red-600 text-sm">{error}</p>;

  return (
    <div>
      <div className="mb-4">
        <input
          type="text"
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Search by phone, suburb, address..."
          className="w-full max-w-md border border-gray-300 rounded-md px-3 py-2 text-sm"
        />
      </div>

      {filtered.length === 0 ? (
        <p className="text-sm text-gray-500">No customers found.</p>
      ) : (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Phone</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Address</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Suburb</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase">Postcode</th>
                <th className="py-3 px-4 text-xs font-medium text-gray-500 uppercase"></th>
              </tr>
            </thead>
            <tbody>
              {filtered.map(c => (
                <tr key={c.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-3 px-4 text-sm">{c.phone || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-gray-600">{c.address || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-gray-600">{c.suburb || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm text-gray-600">{c.postcode || 'N/A'}</td>
                  <td className="py-3 px-4 text-sm">
                    <Link href={`/admin/customers/${c.id}`} className="text-blue-600 hover:underline">
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
            className="bg-white border border-gray-300 text-gray-700 px-4 py-2 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50"
          >
            {loadingMore ? 'Loading...' : 'Load More'}
          </button>
        </div>
      )}
    </div>
  );
}
