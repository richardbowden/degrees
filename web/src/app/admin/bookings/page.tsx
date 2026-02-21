import { cookies } from 'next/headers';
import { BookingsClient } from './bookings-client';

export default async function AdminBookingsPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Bookings</h1>
      <BookingsClient token={token} />
    </div>
  );
}
