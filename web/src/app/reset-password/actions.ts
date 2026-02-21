'use server';

import { api } from '@/lib/api';

export async function resetPasswordAction(
  prevState: { error?: string; success?: string } | null,
  formData: FormData,
) {
  const token = formData.get('token') as string;
  const newPassword = formData.get('new_password') as string;
  const newPasswordConfirm = formData.get('new_password_confirm') as string;

  if (newPassword !== newPasswordConfirm) {
    return { error: 'Passwords do not match.' };
  }

  try {
    const data = await api<{ message: string }>('/auth/complete-password-reset', {
      method: 'POST',
      body: { token, newPassword, newPasswordConfirm },
    });
    return { success: data.message || 'Your password has been reset.' };
  } catch (err: unknown) {
    const error = err as { detail?: string };
    return { error: error.detail || 'Password reset failed. The link may have expired.' };
  }
}
