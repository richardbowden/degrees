import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { User, CustomerProfile } from '@/lib/types';
import { logout } from './actions';

const navItems = [
  { href: '/account', label: 'Dashboard' },
  { href: '/account/profile', label: 'Profile' },
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
      <aside className="w-64 bg-gray-50 border-r border-gray-200 p-6">
        <div className="mb-8">
          <p className="text-sm text-gray-500">Signed in as</p>
          <p className="font-medium text-gray-900 truncate">{user.firstName} {user.surname}</p>
          <p className="text-sm text-gray-500 truncate">{user.email}</p>
        </div>
        <nav className="space-y-1">
          {navItems.map(item => (
            <Link
              key={item.href}
              href={item.href}
              className="block px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-200 hover:text-gray-900"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="mt-8 pt-8 border-t border-gray-200">
          <form action={logout}>
            <button
              type="submit"
              className="w-full text-left px-3 py-2 rounded-md text-sm font-medium text-red-600 hover:bg-red-50"
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
