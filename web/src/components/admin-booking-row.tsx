import Link from 'next/link';
import type { Booking } from '@/lib/types';
import { formatPrice, formatTime } from '@/lib/format';
import { StatusBadge } from '@/components/status-badge';

export function AdminBookingRow({ booking }: { booking: Booking }) {
  return (
    <tr className="border-b border-white/5 hover:bg-white/5">
      <td className="py-3 px-4 text-sm">{formatTime(booking.scheduledTime)}</td>
      <td className="py-3 px-4 text-sm font-medium">
        <Link href={`/admin/bookings/${booking.id}`} className="text-brand-400 hover:text-brand-500">
          {booking.customer?.name ?? 'Unknown'}
        </Link>
      </td>
      <td className="py-3 px-4 text-sm text-text-secondary">
        {booking.vehicle ? `${booking.vehicle.year} ${booking.vehicle.make} ${booking.vehicle.model}` : 'N/A'}
      </td>
      <td className="py-3 px-4 text-sm text-text-secondary">
        {booking.services?.map(s => s.serviceName).join(', ') ?? 'N/A'}
      </td>
      <td className="py-3 px-4 text-sm">{formatPrice(booking.totalAmount)}</td>
      <td className="py-3 px-4"><StatusBadge status={booking.status} /></td>
      <td className="py-3 px-4"><StatusBadge status={booking.paymentStatus} /></td>
      <td className="py-3 px-4 text-sm">
        <Link href={`/admin/bookings/${booking.id}`} className="text-brand-400 hover:underline">
          View
        </Link>
      </td>
    </tr>
  );
}
