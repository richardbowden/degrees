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
        <p className="text-text-secondary mb-6">{state.success}</p>
        <Link href="/login" className="text-brand-400 font-medium hover:underline">
          Go to login
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto px-4 py-16">
      <h1 className="text-2xl font-bold mb-2">Create Account</h1>
      <p className="text-text-secondary mb-8">Join 40 Degrees Car Detailing.</p>

      <form action={formAction} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="first_name" className="block text-sm font-medium text-text-secondary mb-1">
              First Name
            </label>
            <input
              id="first_name"
              name="first_name"
              type="text"
              required
              autoComplete="given-name"
              className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </div>
          <div>
            <label htmlFor="surname" className="block text-sm font-medium text-text-secondary mb-1">
              Surname
            </label>
            <input
              id="surname"
              name="surname"
              type="text"
              required
              autoComplete="family-name"
              className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </div>
        </div>

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
            className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-text-secondary mb-1">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            required
            autoComplete="new-password"
            className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <label htmlFor="password_confirm" className="block text-sm font-medium text-text-secondary mb-1">
            Confirm Password
          </label>
          <input
            id="password_confirm"
            name="password_confirm"
            type="password"
            required
            autoComplete="new-password"
            className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
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
          {pending ? 'Creating account...' : 'Create Account'}
        </button>
      </form>

      <p className="mt-6 text-sm text-text-secondary text-center">
        Already have an account?{' '}
        <Link href="/login" className="text-brand-400 font-medium hover:underline">
          Sign in
        </Link>
      </p>
    </div>
  );
}
