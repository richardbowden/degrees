import { cookies } from 'next/headers';
import { VehicleCategoriesClient } from './vehicle-categories-client';

export default async function AdminVehicleCategoriesPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value ?? '';

  return (
    <div>
      <h1 className="text-2xl font-bold text-white mb-6">Vehicle Categories</h1>
      <VehicleCategoriesClient token={token} />
    </div>
  );
}
