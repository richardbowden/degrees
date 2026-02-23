const statusStyles: Record<string, string> = {
  pending: 'bg-yellow-500/20 text-yellow-400',
  confirmed: 'bg-blue-500/20 text-blue-400',
  in_progress: 'bg-indigo-500/20 text-indigo-400',
  completed: 'bg-green-500/20 text-green-400',
  cancelled: 'bg-red-500/20 text-red-400',
  deposit_paid: 'bg-emerald-500/20 text-emerald-400',
  paid: 'bg-green-500/20 text-green-400',
  partial: 'bg-yellow-500/20 text-yellow-400',
  refunded: 'bg-gray-500/20 text-gray-400',
};

const statusLabels: Record<string, string> = {
  pending: 'Pending',
  confirmed: 'Confirmed',
  in_progress: 'In Progress',
  completed: 'Completed',
  cancelled: 'Cancelled',
  deposit_paid: 'Deposit Paid',
  paid: 'Paid',
  partial: 'Partial',
  refunded: 'Refunded',
};

export function StatusBadge({ status }: { status: string }) {
  const style = statusStyles[status] ?? 'bg-gray-500/20 text-gray-400';
  const label = statusLabels[status] ?? status;
  return (
    <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium ${style}`}>
      {label}
    </span>
  );
}
