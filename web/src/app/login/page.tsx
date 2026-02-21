'use client';

import { useActionState, use } from 'react';
import Link from 'next/link';
import { loginAction } from './actions';

interface Props {
  searchParams: Promise<{ redirect?: string }>;
}

export default function LoginPage({ searchParams }: Props) {
  const params = use(searchParams);
  const redirectTo = params.redirect ?? '';
  const [state, formAction, pending] = useActionState(loginAction, null);

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Login</h1>
      <p className="text-gray-600 mb-8">Sign in to your 40 Degrees account.</p>

      <form action={formAction} className="space-y-4">
        {redirectTo && <input type="hidden" name="redirect" value={redirectTo} />}
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
            autoComplete="current-password"
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
          {pending ? 'Signing in...' : 'Sign In'}
        </button>
      </form>

      <div className="mt-6 text-sm text-gray-600 text-center space-y-2">
        <p>
          <Link href="/forgot-password" className="text-gray-900 font-medium hover:underline">
            Forgot your password?
          </Link>
        </p>
        <p>
          Don&apos;t have an account?{' '}
          <Link href="/register" className="text-gray-900 font-medium hover:underline">
            Create one
          </Link>
        </p>
      </div>
    </div>
  );
}
