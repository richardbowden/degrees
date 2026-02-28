import Link from 'next/link';
import { cookies } from 'next/headers';
import { CartIcon } from '@/components/cart-icon';
import { ThemeToggle } from '@/components/theme-toggle';
import { logout } from '@/app/account/actions';

export async function Nav() {
  const cookieStore = await cookies();
  const isLoggedIn = !!cookieStore.get('session_token')?.value;

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
          <CartIcon />
          {isLoggedIn ? (
            <>
              <Link href="/account" className="text-text-secondary hover:text-foreground transition-colors">
                Account
              </Link>
              <form action={logout}>
                <button
                  type="submit"
                  className="text-text-secondary hover:text-foreground transition-colors"
                >
                  Log Out
                </button>
              </form>
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
