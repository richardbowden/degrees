import { cookies } from 'next/headers';
import { UsersClient } from './users-client';

export default async function AdminUsersPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Users</h1>
      <UsersClient token={token} />
    </div>
  );
}
