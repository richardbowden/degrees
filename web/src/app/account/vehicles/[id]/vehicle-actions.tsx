'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api, ApiError } from '@/lib/api';
import { VehicleForm, VehicleFormData } from '@/components/vehicle-form';
import type { Vehicle } from '@/lib/types';

interface VehicleActionsProps {
  vehicle: Vehicle;
}

export function VehicleActions({ vehicle }: VehicleActionsProps) {
  const router = useRouter();
  const [editing, setEditing] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState('');

  async function handleEdit(data: VehicleFormData) {
    setError('');
    try {
      await api<{ vehicle: Vehicle }>(`/me/vehicles/${vehicle.id}`, {
        method: 'PUT',
        body: data,
      });
      router.refresh();
      setEditing(false);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to update vehicle.');
    }
  }

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete this vehicle?')) return;
    setDeleting(true);
    setError('');
    try {
      await api(`/me/vehicles/${vehicle.id}`, { method: 'DELETE' });
      router.push('/account/vehicles');
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to delete vehicle.');
      setDeleting(false);
    }
  }

  if (editing) {
    return (
      <div>
        {error && <p className="text-sm text-red-400 mb-4">{error}</p>}
        <VehicleForm
          vehicle={vehicle}
          onSubmit={handleEdit}
          onCancel={() => setEditing(false)}
        />
      </div>
    );
  }

  return (
    <div>
      {error && <p className="text-sm text-red-400 mb-4">{error}</p>}
      <div className="flex gap-3">
        <button
          onClick={() => setEditing(true)}
          className="btn-brand px-4 py-2 text-sm font-medium rounded-md"
        >
          Edit Vehicle
        </button>
        <button
          onClick={handleDelete}
          disabled={deleting}
          className="border border-red-400/30 text-red-400 px-4 py-2 rounded-md text-sm font-medium hover:bg-red-400/10 disabled:opacity-50"
        >
          {deleting ? 'Deleting...' : 'Delete Vehicle'}
        </button>
      </div>
    </div>
  );
}
