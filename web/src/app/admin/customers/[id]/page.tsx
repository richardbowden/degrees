import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import type { CustomerProfile, Vehicle, ServiceRecord } from '@/lib/types';
import { CustomerDetailClient } from './customer-detail-client';

export default async function AdminCustomerDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  let profile: CustomerProfile | null = null;
  let vehicles: Vehicle[] = [];
  let records: ServiceRecord[] = [];
  let error = '';

  try {
    const [customerRes, historyRes] = await Promise.all([
      api<{ profile: CustomerProfile; vehicles: Vehicle[] }>(`/admin/customers/${id}`, { token }),
      api<{ records: ServiceRecord[] }>(`/admin/customers/${id}/history`, { token }),
    ]);
    profile = customerRes.profile;
    vehicles = customerRes.vehicles ?? [];
    records = historyRes.records ?? [];
  } catch {
    error = 'Failed to load customer';
  }

  if (error || !profile) {
    return (
      <div>
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Customer Detail</h1>
        <p className="text-red-600 text-sm">{error || 'Customer not found'}</p>
      </div>
    );
  }

  return (
    <CustomerDetailClient
      profile={profile}
      vehicles={vehicles}
      initialRecords={records}
      token={token}
    />
  );
}
