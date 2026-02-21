'use server';

import { api } from '@/lib/api';

export async function registerAction(
  prevState: { error?: string; success?: string } | null,
  formData: FormData
) {
  const email = formData.get('email') as string;
  const firstName = formData.get('first_name') as string;
  const surname = formData.get('surname') as string;
  const password = formData.get('password') as string;
  const passwordConfirm = formData.get('password_confirm') as string;

  if (password !== passwordConfirm) {
    return { error: 'Passwords do not match' };
  }

  try {
    const data = await api<{ message: string; signupId: string }>('/auth/register', {
      method: 'POST',
      body: { email, firstName, surname, password, passwordConfirm },
    });

    return { success: data.message || 'Check your email to verify your account.' };
  } catch (err: unknown) {
    const error = err as { detail?: string };
    return { error: error.detail || 'Registration failed' };
  }
}
