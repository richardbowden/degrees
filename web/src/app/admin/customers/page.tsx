import { cookies } from 'next/headers';
import { CustomersClient } from './customers-client';

export default async function AdminCustomersPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Customers</h1>
      <CustomersClient token={token} />
    </div>
  );
}
