'use server';

import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import type { User } from '@/lib/types';

export async function loginAction(
  prevState: { error?: string; redirectTo?: string } | null,
  formData: FormData
) {
  const email = formData.get('email') as string;
  const password = formData.get('password') as string;
  const redirectTo = formData.get('redirect') as string;

  try {
    const data = await api<{ sessionToken: string; user: User }>('/auth/login', {
      method: 'POST',
      body: { email, password },
    });

    const cookieStore = await cookies();
    cookieStore.set('session_token', data.sessionToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 60 * 60 * 24 * 30,
    });

    // Return the redirect URL — the client does a hard navigation so the
    // root layout re-renders and the Nav reflects the new auth state.
    return { redirectTo: redirectTo && redirectTo.startsWith('/') ? redirectTo : '/account' };
  } catch (err: unknown) {
    const error = err as { detail?: string; message?: string };
    return { error: error.detail || error.message || 'Login failed' };
  }
}
