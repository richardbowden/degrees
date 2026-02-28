import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { User, CustomerProfile } from '@/lib/types';
import { logout } from './actions';

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

  let user: User;
  try {
    const { profile } = await api<{ profile: CustomerProfile }>('/me/profile', { token });
    const res = await api<{ user: User }>(`/user/${profile.userId}`, { token });
    user = res.user;
  } catch {
    redirect('/login');
  }

  return (
    <div className="flex min-h-screen">
      <aside className="w-64 bg-surface-raised border-r border-border-subtle p-6">
        <div className="mb-8">
          <p className="text-sm text-text-muted">Signed in as</p>
          <p className="font-medium text-foreground truncate">{user.firstName} {user.surname}</p>
          <p className="text-sm text-text-muted truncate">{user.email}</p>
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
          <form action={logout}>
            <button
              type="submit"
              className="w-full text-left px-3 py-2 rounded-md text-sm font-medium text-red-400 hover:bg-surface-hover"
            >
              Log Out
            </button>
          </form>
        </div>
      </aside>
      <div className="flex-1 p-8">
        {children}
      </div>
    </div>
  );
}
