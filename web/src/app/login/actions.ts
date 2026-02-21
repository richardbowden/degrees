'use server';

import { cookies } from 'next/headers';
import { api } from '@/lib/api';
import { redirect } from 'next/navigation';
import type { User } from '@/lib/types';

export async function loginAction(prevState: { error?: string } | null, formData: FormData) {
  const email = formData.get('email') as string;
  const password = formData.get('password') as string;

  try {
    // Backend returns camelCase (protobuf JSON): sessionToken
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
  } catch (err: unknown) {
    const error = err as { detail?: string; message?: string };
    return { error: error.detail || error.message || 'Login failed' };
  }

  redirect('/account');
}
