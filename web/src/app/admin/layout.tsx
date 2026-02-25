import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { User } from '@/lib/types';

interface ProfileResponse {
  profile: { userId: string };
}

const NAV_ITEMS = [
  { href: '/admin', label: 'Dashboard' },
  { href: '/admin/bookings', label: 'Bookings' },
  { href: '/admin/customers', label: 'Customers' },
  { href: '/admin/services', label: 'Services' },
  { href: '/admin/vehicle-categories', label: 'Vehicle Categories' },
  { href: '/admin/schedule', label: 'Schedule' },
  { href: '/admin/users', label: 'Users' },
];

export default async function AdminLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;

  if (!token) redirect('/login?redirect=/admin');

  try {
    const { profile } = await api<ProfileResponse>('/me/profile', { token });
    const { user } = await api<{ user: User }>(`/user/${profile.userId}`, { token });
    if (!user.sysop) redirect('/account');
  } catch {
    redirect('/account');
  }

  return (
    <div className="flex min-h-screen">
      <aside className="w-64 bg-surface-raised text-white p-6 flex-shrink-0">
        <div className="mb-8">
          <h1 className="text-lg font-bold">40 Degrees</h1>
          <p className="text-xs text-text-muted mt-1">Admin Panel</p>
        </div>
        <nav className="space-y-1">
          {NAV_ITEMS.map(item => (
            <Link
              key={item.href}
              href={item.href}
              className="block px-3 py-2 rounded-md text-sm font-medium text-text-secondary hover:text-white hover:bg-white/5 transition-colors"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="mt-auto pt-8">
          <Link href="/" className="text-xs text-text-muted hover:text-white">
            Back to site
          </Link>
        </div>
      </aside>
      <div className="flex-1 p-8 min-h-screen">{children}</div>
    </div>
  );
}
