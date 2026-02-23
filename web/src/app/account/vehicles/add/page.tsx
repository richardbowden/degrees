'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import { VehicleForm, VehicleFormData } from '@/components/vehicle-form';
import type { Vehicle } from '@/lib/types';

export default function AddVehiclePage() {
  const router = useRouter();
  const [error, setError] = useState('');

  async function handleSubmit(data: VehicleFormData) {
    setError('');
    try {
      await api<{ vehicle: Vehicle }>('/me/vehicles', {
        method: 'POST',
        body: data,
      });
      router.push('/account/vehicles');
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to add vehicle.');
    }
  }

  return (
    <div>
      <Link
        href="/account/vehicles"
        className="text-sm text-text-muted hover:text-white mb-4 inline-block"
      >
        &larr; Back to Vehicles
      </Link>
      <h1 className="text-2xl font-bold text-white mb-6">Add Vehicle</h1>

      {error && <p className="text-sm text-red-400 mb-4">{error}</p>}

      <div className="max-w-2xl">
        <VehicleForm
          onSubmit={handleSubmit}
          onCancel={() => router.push('/account/vehicles')}
        />
      </div>
    </div>
  );
}
