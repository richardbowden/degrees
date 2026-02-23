'use client';

import { useState, useEffect } from 'react';
import { api, ApiError } from '@/lib/api';
import type { User } from '@/lib/types';
import { formatDate } from '@/lib/format';

export function UsersClient({ token }: { token: string }) {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  useEffect(() => {
    api<{ users: User[] }>('/admin/users', { token })
      .then(data => setUsers(data.users ?? []))
      .catch((err: ApiError) => setError(err.detail || 'Failed to load users'))
      .finally(() => setLoading(false));
  }, [token]);

  async function toggleEnabled(user: User) {
    setActionLoading(user.id);
    setError('');
    const endpoint = user.enabled
      ? `/user/${user.id}/disable`
      : `/user/${user.id}/enable`;
    try {
      const res = await api<{ user: User }>(endpoint, {
        method: 'POST',
        body: {},
        token,
      });
      setUsers(prev => prev.map(u => (u.id === user.id ? res.user : u)));
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to update user');
    } finally {
      setActionLoading(null);
    }
  }

  if (loading) {
    return <p className="text-text-muted text-sm">Loading users...</p>;
  }

  if (error && users.length === 0) {
    return <p className="text-red-400 text-sm">{error}</p>;
  }

  return (
    <>
      {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
      <div className="glass-card overflow-hidden">
        <table className="w-full text-left">
          <thead className="bg-white/5 border-b border-border-subtle">
            <tr>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Name</th>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Email</th>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Joined</th>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Status</th>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase">Role</th>
              <th className="py-3 px-4 text-xs font-medium text-text-muted uppercase"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border-subtle">
            {users.map(user => (
              <tr key={user.id}>
                <td className="py-3 px-4 text-sm font-medium text-white">
                  {user.firstName} {user.surname}
                </td>
                <td className="py-3 px-4 text-sm text-text-secondary">{user.email}</td>
                <td className="py-3 px-4 text-sm text-text-secondary">{formatDate(user.createdOn)}</td>
                <td className="py-3 px-4">
                  <span
                    className={`inline-block px-2 py-0.5 text-xs font-medium rounded-full ${
                      user.enabled
                        ? 'bg-green-500/20 text-green-400'
                        : 'bg-red-500/20 text-red-400'
                    }`}
                  >
                    {user.enabled ? 'Active' : 'Disabled'}
                  </span>
                </td>
                <td className="py-3 px-4">
                  {user.sysop && (
                    <span className="inline-block px-2 py-0.5 text-xs font-medium rounded-full bg-purple-500/20 text-purple-400">
                      Admin
                    </span>
                  )}
                </td>
                <td className="py-3 px-4 text-right">
                  <button
                    onClick={() => toggleEnabled(user)}
                    disabled={actionLoading === user.id}
                    className={`text-sm font-medium disabled:opacity-50 ${
                      user.enabled
                        ? 'text-red-400 hover:text-red-500'
                        : 'text-green-400 hover:text-green-500'
                    }`}
                  >
                    {actionLoading === user.id
                      ? 'Updating...'
                      : user.enabled
                        ? 'Disable'
                        : 'Enable'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
