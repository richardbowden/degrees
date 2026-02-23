import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import type { CustomerProfile, Vehicle } from '@/lib/types';
import { ProfileClient } from './profile-client';

export default async function ProfilePage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const [profileRes, vehiclesRes] = await Promise.all([
    api<{ profile: CustomerProfile }>('/me/profile', { token }),
    api<{ vehicles: Vehicle[] }>('/me/vehicles', { token }),
  ]);

  return (
    <div>
      <h1 className="text-2xl font-bold text-white mb-6">Profile</h1>
      <ProfileClient
        initialProfile={profileRes.profile}
        initialVehicles={vehiclesRes.vehicles}
      />
    </div>
  );
}
