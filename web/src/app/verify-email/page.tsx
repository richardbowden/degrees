'use client';

import { useState, useEffect, use } from 'react';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';

interface Props {
  searchParams: Promise<{ token?: string }>;
}

export default function VerifyEmailPage({ searchParams }: Props) {
  const params = use(searchParams);
  const token = params.token ?? '';
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (!token) {
      setStatus('error');
      setMessage('No verification token provided.');
      return;
    }

    api<{ message: string }>('/auth/verify-email', {
      method: 'POST',
      body: { token },
    })
      .then(data => {
        setStatus('success');
        setMessage(data.message || 'Your email has been verified.');
      })
      .catch((err: ApiError) => {
        setStatus('error');
        setMessage(err.detail || 'Verification failed. The link may have expired.');
      });
  }, [token]);

  return (
    <div className="max-w-md mx-auto px-4 py-16 text-center">
      {status === 'loading' && (
        <>
          <h1 className="text-2xl font-bold mb-4">Verifying Email</h1>
          <p className="text-gray-500">Please wait...</p>
        </>
      )}

      {status === 'success' && (
        <>
          <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <span className="text-green-600 text-2xl">&#10003;</span>
          </div>
          <h1 className="text-2xl font-bold mb-4">Email Verified</h1>
          <p className="text-gray-600 mb-6">{message}</p>
          <Link
            href="/login"
            className="inline-block bg-gray-900 text-white px-6 py-2.5 rounded font-semibold hover:bg-gray-800"
          >
            Sign In
          </Link>
        </>
      )}

      {status === 'error' && (
        <>
          <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <span className="text-red-600 text-2xl">&#10007;</span>
          </div>
          <h1 className="text-2xl font-bold mb-4">Verification Failed</h1>
          <p className="text-gray-600 mb-6">{message}</p>
          <Link
            href="/register"
            className="text-gray-900 font-medium hover:underline"
          >
            Try registering again
          </Link>
        </>
      )}
    </div>
  );
}
