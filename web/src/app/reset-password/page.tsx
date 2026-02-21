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
        <p className="text-gray-600 mb-6">No reset token provided.</p>
        <Link href="/forgot-password" className="text-gray-900 font-medium hover:underline">
          Request a new reset link
        </Link>
      </div>
    );
  }

  if (state?.success) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <span className="text-green-600 text-2xl">&#10003;</span>
        </div>
        <h1 className="text-2xl font-bold mb-4">Password Reset</h1>
        <p className="text-gray-600 mb-6">{state.success}</p>
        <Link
          href="/login"
          className="inline-block bg-gray-900 text-white px-6 py-2.5 rounded font-semibold hover:bg-gray-800"
        >
          Sign In
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Set New Password</h1>
      <p className="text-gray-600 mb-8">Enter your new password below.</p>

      <form action={formAction} className="space-y-4">
        <input type="hidden" name="token" value={token} />

        <div>
          <label htmlFor="new_password" className="block text-sm font-medium text-gray-700 mb-1">
            New Password
          </label>
          <input
            id="new_password"
            name="new_password"
            type="password"
            required
            autoComplete="new-password"
            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
          />
        </div>

        <div>
          <label htmlFor="new_password_confirm" className="block text-sm font-medium text-gray-700 mb-1">
            Confirm New Password
          </label>
          <input
            id="new_password_confirm"
            name="new_password_confirm"
            type="password"
            required
            autoComplete="new-password"
            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
          />
        </div>

        {state?.error && (
          <p className="text-red-600 text-sm">{state.error}</p>
        )}

        <button
          type="submit"
          disabled={pending}
          className="w-full bg-gray-900 text-white py-2.5 rounded font-semibold hover:bg-gray-800 disabled:opacity-50"
        >
          {pending ? 'Resetting...' : 'Reset Password'}
        </button>
      </form>
    </div>
  );
}
