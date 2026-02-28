'use client';

import { useActionState, use } from 'react';
import Link from 'next/link';
import { resetPasswordAction } from './actions';

interface Props {
  searchParams: Promise<{ token?: string }>;
}

export default function ResetPasswordPage({ searchParams }: Props) {
  const params = use(searchParams);
  const token = params.token ?? '';
  const [state, formAction, pending] = useActionState(resetPasswordAction, null);

  if (!token) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">Invalid Link</h1>
        <p className="text-text-secondary mb-6">No reset token provided.</p>
        <Link href="/forgot-password" className="text-brand-400 font-medium hover:underline">
          Request a new reset link
        </Link>
      </div>
    );
  }

  if (state?.success) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <div className="w-16 h-16 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-6">
          <span className="text-green-400 text-2xl">&#10003;</span>
        </div>
        <h1 className="text-2xl font-bold mb-4">Password Reset</h1>
        <p className="text-text-secondary mb-6">{state.success}</p>
        <Link
          href="/login"
          className="inline-block btn-brand px-6 py-2.5"
        >
          Sign In
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Set New Password</h1>
      <p className="text-text-secondary mb-8">Enter your new password below.</p>

      <form action={formAction} className="space-y-4">
        <input type="hidden" name="token" value={token} />

        <div>
          <label htmlFor="new_password" className="block text-sm font-medium text-text-secondary mb-1">
            New Password
          </label>
          <input
            id="new_password"
            name="new_password"
            type="password"
            required
            autoComplete="new-password"
            className="w-full bg-surface-input border border-border-subtle rounded px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <label htmlFor="new_password_confirm" className="block text-sm font-medium text-text-secondary mb-1">
            Confirm New Password
          </label>
          <input
            id="new_password_confirm"
            name="new_password_confirm"
            type="password"
            required
            autoComplete="new-password"
            className="w-full bg-surface-input border border-border-subtle rounded px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        {state?.error && (
          <p className="text-red-400 text-sm">{state.error}</p>
        )}

        <button
          type="submit"
          disabled={pending}
          className="w-full btn-brand py-2.5"
        >
          {pending ? 'Resetting...' : 'Reset Password'}
        </button>
      </form>
    </div>
  );
}
