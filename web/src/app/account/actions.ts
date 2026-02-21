'use server';

import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { api } from '@/lib/api';

export async function logout() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;

  if (token) {
    try {
      await api('/auth/logout', { method: 'POST', body: { sessionToken: token } });
    } catch {
      // Proceed with logout even if API call fails
    }
  }

  cookieStore.delete('session_token');
  redirect('/');
}
