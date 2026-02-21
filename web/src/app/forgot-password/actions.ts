'use server';

import { api } from '@/lib/api';

export async function forgotPasswordAction(
  prevState: { error?: string; success?: string } | null,
  formData: FormData,
) {
  const email = formData.get('email') as string;

  try {
    const data = await api<{ message: string }>('/auth/reset-password', {
      method: 'POST',
      body: { email },
    });
    return { success: data.message || 'If an account with that email exists, a reset link has been sent.' };
  } catch (err: unknown) {
    const error = err as { detail?: string };
    return { error: error.detail || 'Failed to send reset email.' };
  }
}
