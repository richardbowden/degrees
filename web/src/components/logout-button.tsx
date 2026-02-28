'use client';

import { useTransition } from 'react';
import { logout } from '@/app/account/actions';

interface Props {
  className?: string;
  children?: React.ReactNode;
}

export function LogoutButton({ className, children = 'Log Out' }: Props) {
  const [pending, startTransition] = useTransition();

  function handleLogout() {
    startTransition(async () => {
      await logout();
      window.location.href = '/';
    });
  }

  return (
    <button onClick={handleLogout} disabled={pending} className={className}>
      {pending ? 'Logging out…' : children}
    </button>
  );
}
