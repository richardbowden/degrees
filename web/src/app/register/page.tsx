'use client';

import { useActionState } from 'react';
import Link from 'next/link';
import { registerAction } from './actions';

export default function RegisterPage() {
  const [state, formAction, pending] = useActionState(registerAction, null);

  if (state?.success) {
    return (
      <div className="max-w-md mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">Account Created</h1>
        <p className="text-gray-600 mb-6">{state.success}</p>
        <Link href="/login" className="text-gray-900 font-medium hover:underline">
          Go to login
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Create Account</h1>
      <p className="text-gray-600 mb-8">Join 40 Degrees Car Detailing.</p>

      <form action={formAction} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="first_name" className="block text-sm font-medium text-gray-700 mb-1">
              First Name
            </label>
            <input
              id="first_name"
              name="first_name"
              type="text"
              required
              autoComplete="given-name"
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
            />
          </div>
          <div>
            <label htmlFor="surname" className="block text-sm font-medium text-gray-700 mb-1">
              Surname
            </label>
            <input
              id="surname"
              name="surname"
              type="text"
              required
              autoComplete="family-name"
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
            />
          </div>
        </div>

        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            required
            autoComplete="email"
            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            autoComplete="new-password"
            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900"
          />
        </div>

        <div>
          <label htmlFor="password_confirm" className="block text-sm font-medium text-gray-700 mb-1">
            Confirm Password
          </label>
          <input
            id="password_confirm"
            name="password_confirm"
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
          {pending ? 'Creating account...' : 'Create Account'}
        </button>
      </form>

      <p className="mt-6 text-sm text-gray-600 text-center">
        Already have an account?{' '}
        <Link href="/login" className="text-gray-900 font-medium hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
