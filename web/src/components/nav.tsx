import Link from 'next/link';
import { cookies } from 'next/headers';
import { CartIcon } from '@/components/cart-icon';
import { logout } from '@/app/account/actions';

export async function Nav() {
  const cookieStore = await cookies();
  const isLoggedIn = !!cookieStore.get('session_token')?.value;

  return (
    <nav className="border-b border-gray-200 bg-white">
      <div className="max-w-6xl mx-auto px-4 flex items-center justify-between h-16">
        <Link href="/" className="text-xl font-bold text-gray-900">
          40 Degrees
        </Link>
        <div className="flex items-center gap-6">
          <Link href="/services" className="text-gray-700 hover:text-gray-900">
            Services
          </Link>
          <CartIcon />
          {isLoggedIn ? (
            <>
              <Link href="/account" className="text-gray-700 hover:text-gray-900">
                Account
              </Link>
              <form action={logout}>
                <button
                  type="submit"
                  className="text-gray-700 hover:text-gray-900"
                >
                  Log Out
                </button>
              </form>
            </>
          ) : (
            <Link href="/login" className="text-gray-700 hover:text-gray-900">
              Login
            </Link>
          )}
        </div>
      </div>
    </nav>
  );
}
