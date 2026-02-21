import { cookies } from 'next/headers';
import { ScheduleClient } from './schedule-client';

export default async function AdminSchedulePage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Schedule</h1>
      <ScheduleClient token={token} />
    </div>
  );
}
