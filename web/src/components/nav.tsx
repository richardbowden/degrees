import Link from 'next/link';
import { cookies } from 'next/headers';
import { CartIcon } from '@/components/cart-icon';
import { ThemeToggle } from '@/components/theme-toggle';
import { LogoutButton } from '@/components/logout-button';
import { api } from '@/lib/api';
import type { User, CustomerProfile } from '@/lib/types';

async function getSysopStatus(token: string): Promise<boolean> {
  try {
    const { profile } = await api<{ profile: CustomerProfile }>('/me/profile', { token });
    const { user } = await api<{ user: User }>(`/user/${profile.userId}`, { token });
    return user.sysop;
  } catch {
    return false;
  }
}

export async function Nav() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value;
  const isLoggedIn = !!token;
  const isSysop = token ? await getSysopStatus(token) : false;

  return (
    <nav className="sticky top-0 z-50 bg-surface/95 backdrop-blur-lg border-b border-border-subtle">
      <div className="max-w-6xl mx-auto px-4 flex items-center justify-between h-16">
        <Link href="/" className="text-xl font-bold text-brand-gradient">
          40 Degrees
        </Link>
        <div className="flex items-center gap-6">
          <Link href="/services" className="text-text-secondary hover:text-foreground transition-colors">
            Services
          </Link>
          {isSysop && (
            <Link
              href="/admin"
              className="flex items-center gap-1.5 text-brand-400 hover:text-brand-500 transition-colors font-medium"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
                />
              </svg>
              Admin
            </Link>
          )}
          <CartIcon />
          {isLoggedIn ? (
            <>
              <Link href="/account" className="text-text-secondary hover:text-foreground transition-colors">
                Account
              </Link>
              <LogoutButton className="text-text-secondary hover:text-foreground transition-colors" />
            </>
          ) : (
            <Link href="/login" className="text-text-secondary hover:text-foreground transition-colors">
              Login
            </Link>
          )}
          <ThemeToggle />
        </div>
      </div>
    </nav>
  );
}
