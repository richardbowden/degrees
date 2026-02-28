'use client';

import { useActionState } from 'react';
import Link from 'next/link';
import { forgotPasswordAction } from './actions';

export default function ForgotPasswordPage() {
  const [state, formAction, pending] = useActionState(forgotPasswordAction, null);

  if (state?.success) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">Check Your Email</h1>
        <p className="text-text-secondary mb-6">{state.success}</p>
        <Link href="/login" className="text-brand-400 font-medium hover:underline">
          Back to login
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Reset Password</h1>
      <p className="text-text-secondary mb-8">
        Enter your email and we&apos;ll send you a link to reset your password.
      </p>

      <form action={formAction} className="space-y-4">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-text-secondary mb-1">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            required
            autoComplete="email"
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
          {pending ? 'Sending...' : 'Send Reset Link'}
        </button>
      </form>

      <p className="mt-6 text-sm text-text-secondary text-center">
        <Link href="/login" className="text-brand-400 font-medium hover:underline">
          Back to login
        </Link>
      </p>
    </div>
  );
}
