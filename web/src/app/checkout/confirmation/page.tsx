import Link from 'next/link';

interface Props {
  searchParams: Promise<{ booking_id?: string }>;
}

export default async function ConfirmationPage({ searchParams }: Props) {
  const params = await searchParams;
  const bookingId = params.booking_id ?? '';

  return (
    <div className="max-w-md mx-auto px-4 py-16 text-center">
      <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
        <span className="text-green-600 text-2xl">&#10003;</span>
      </div>
      <h1 className="text-2xl font-bold mb-2">Booking Confirmed</h1>
      <p className="text-gray-600 mb-2">
        Your booking has been confirmed and your deposit received.
      </p>
      {bookingId && (
        <p className="text-sm text-gray-400 mb-8">Booking ID: {bookingId}</p>
      )}
      <div className="flex flex-col gap-3">
        <Link
          href="/account/bookings"
          className="block bg-gray-900 text-white py-2.5 rounded font-semibold hover:bg-gray-800"
        >
          View My Bookings
        </Link>
        <Link
          href="/services"
          className="text-gray-600 hover:underline text-sm"
        >
          Browse more services
        </Link>
      </div>
    </div>
  );
}
