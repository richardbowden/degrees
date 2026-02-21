'use client';

import { useState, use } from 'react';
import { useRouter } from 'next/navigation';
import { api, ApiError } from '@/lib/api';
import { formatPrice } from '@/lib/format';

interface Props {
  searchParams: Promise<{ booking_id?: string }>;
}

export default function DepositPage({ searchParams }: Props) {
  const params = use(searchParams);
  const router = useRouter();
  const bookingId = params.booking_id ?? '';
  const [loading, setLoading] = useState(false);
  const [depositAmount, setDepositAmount] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handlePayDeposit() {
    if (!bookingId) return;
    setLoading(true);
    setError(null);
    try {
      const data = await api<{ clientSecret: string; depositAmount: number }>('/checkout/deposit', {
        method: 'POST',
        body: { bookingId: bookingId },
      });
      setDepositAmount(data.depositAmount);
      // Stripe integration TBD - for now redirect to confirmation
      router.push(`/checkout/confirmation?booking_id=${bookingId}`);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to process deposit');
      setLoading(false);
    }
  }

  if (!bookingId) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <p className="text-gray-500">No booking ID provided.</p>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-4">Pay Deposit</h1>
      <p className="text-gray-600 mb-6">
        A 30% deposit is required to confirm your booking.
      </p>

      <div className="border border-gray-200 rounded-lg p-6 mb-6">
        <p className="text-sm text-gray-500">Booking</p>
        <p className="font-mono text-sm mb-4">{bookingId}</p>
        {depositAmount !== null && (
          <div>
            <p className="text-sm text-gray-500">Deposit Amount</p>
            <p className="text-xl font-bold">{formatPrice(depositAmount)}</p>
          </div>
        )}
      </div>

      {error && (
        <p className="text-red-600 text-sm mb-4">{error}</p>
      )}

      <button
        type="button"
        onClick={handlePayDeposit}
        disabled={loading}
        className="w-full bg-gray-900 text-white py-3 rounded font-semibold hover:bg-gray-800 disabled:opacity-50"
      >
        {loading ? 'Processing...' : 'Pay Deposit'}
      </button>
    </div>
  );
}
