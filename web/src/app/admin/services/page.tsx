import { cookies } from 'next/headers';
import { ServicesClient } from './services-client';

export default async function AdminServicesPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-white mb-6">Services</h1>
      <ServicesClient token={token} />
    </div>
  );
}
