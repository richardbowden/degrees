import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { api } from '@/lib/api';
import type { CustomerProfile } from '@/lib/types';
import { ProfileClient } from './profile-client';

export default async function ProfilePage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  if (!token) redirect('/login');

  try {
    const { profile } = await api<{ profile: CustomerProfile }>('/me/profile', { token });
    return (
      <div>
        <h1 className="text-2xl font-bold mb-6">Profile</h1>
        <ProfileClient initialProfile={profile} />
      </div>
    );
  } catch (err) {
    const apiErr = err as { status?: number };
    if (!apiErr.status || apiErr.status === 401 || apiErr.status === 403) {
      redirect('/login');
    }
    return (
      <div>
        <h1 className="text-2xl font-bold mb-6">Profile</h1>
        <p className="text-red-400 text-sm">Failed to load profile. Please try again.</p>
      </div>
    );
  }
}
