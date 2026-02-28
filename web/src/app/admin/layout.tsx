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
      <aside className="w-64 bg-surface-raised text-foreground p-6 flex-shrink-0 border-r border-border-subtle">
        <div className="mb-8">
          <Link href="/" className="text-lg font-bold text-brand-gradient block">
            40 Degrees
          </Link>
          <div className="flex items-center gap-1.5 mt-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-3.5 h-3.5 text-brand-400"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
              />
            </svg>
            <span className="text-xs font-semibold text-brand-400 uppercase tracking-wider">Admin</span>
          </div>
        </div>
        <nav className="space-y-1">
          {NAV_ITEMS.map(item => (
            <Link
              key={item.href}
              href={item.href}
              className="block px-3 py-2 rounded-md text-sm font-medium text-text-secondary hover:text-foreground hover:bg-surface-hover transition-colors"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="mt-auto pt-8">
          <Link href="/" className="text-xs text-text-muted hover:text-text-secondary transition-colors">
            ← Back to site
          </Link>
        </div>
      </aside>
      <div className="flex-1 p-8 min-h-screen">{children}</div>
    </div>
  );
}
