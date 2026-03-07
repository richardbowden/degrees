import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { User, CustomerProfile } from '@/lib/types';
import { LogoutButton } from '@/components/logout-button';

const navItems = [
  { href: '/account', label: 'Dashboard' },
  { href: '/account/profile', label: 'Profile' },
  { href: '/account/vehicles', label: 'Vehicles' },
  { href: '/account/bookings', label: 'Bookings' },
  { href: '/account/history', label: 'History' },
];

export default async function AccountLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;

  if (!token) {
    redirect('/login');
  }

  let user: User | null = null;
  try {
    const { profile } = await api<{ profile: CustomerProfile }>('/me/profile', { token });
    const res = await api<{ user: User }>(`/user/${profile.userId}`, { token });
    user = res.user;
  } catch (err) {
    const apiErr = err as { status?: number };
    // Only redirect to login for auth failures — other errors (500, network) should
    // not bounce the user to login or they'll be stuck in an infinite loop.
    if (!apiErr.status || apiErr.status === 401 || apiErr.status === 403) {
      redirect('/login');
    }
    // For other errors, render the layout without user info rather than looping.
  }

  return (
    <div className="flex min-h-screen">
      <aside className="w-64 bg-surface-raised border-r border-border-subtle p-6">
        <div className="mb-8">
          <p className="text-sm text-text-muted">Signed in as</p>
          <p className="font-medium text-foreground truncate">{user ? `${user.firstName} ${user.surname}` : 'My Account'}</p>
          {user && <p className="text-sm text-text-muted truncate">{user.email}</p>}
        </div>
        <nav className="space-y-1">
          {navItems.map(item => (
            <Link
              key={item.href}
              href={item.href}
              className="block px-3 py-2 rounded-md text-sm font-medium text-text-secondary hover:bg-surface-hover hover:text-foreground"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="mt-8 pt-8 border-t border-border-subtle">
          <LogoutButton className="w-full text-left px-3 py-2 rounded-md text-sm font-medium text-red-400 hover:bg-surface-hover" />
        </div>
      </aside>
      <div className="flex-1 p-8">
        {children}
      </div>
    </div>
  );
}
