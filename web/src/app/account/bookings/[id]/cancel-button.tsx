'use client';

import { useState } from 'react';
import { api, ApiError } from '@/lib/api';
import { useRouter } from 'next/navigation';

export function CancelButton({ bookingId }: { bookingId: string }) {
  const [cancelling, setCancelling] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  async function handleCancel() {
    if (!confirm('Are you sure you want to cancel this booking? Cancellations within 24 hours may forfeit the deposit.')) {
      return;
    }

    setCancelling(true);
    setError('');
    try {
      await api(`/me/bookings/${bookingId}/cancel`, { method: 'POST' });
      router.refresh();
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to cancel booking.');
    } finally {
      setCancelling(false);
    }
  }

  return (
    <div>
      {error && <p className="text-sm text-red-400 mb-2">{error}</p>}
      <button
        onClick={handleCancel}
        disabled={cancelling}
        className="px-4 py-2 border border-red-400/30 text-red-400 text-sm font-medium rounded-md hover:bg-red-400/10 disabled:opacity-50"
      >
        {cancelling ? 'Cancelling...' : 'Cancel Booking'}
      </button>
    </div>
  );
}
