import { cookies } from 'next/headers';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Vehicle } from '@/lib/types';

export default async function VehiclesPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value!;

  const { vehicles } = await api<{ vehicles: Vehicle[] }>('/me/vehicles', { token });

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-white">Vehicles</h1>
        <Link
          href="/account/vehicles/add"
          className="btn-brand px-4 py-2 text-sm font-medium rounded-md"
        >
          Add Vehicle
        </Link>
      </div>

      {vehicles.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-text-muted text-lg mb-4">No vehicles added yet.</p>
          <p className="text-text-muted text-sm mb-6">
            Add your vehicle to get started with booking a detail.
          </p>
          <Link
            href="/account/vehicles/add"
            className="btn-brand px-6 py-3 text-sm font-medium rounded-md"
          >
            Add Your First Vehicle
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {vehicles.map(vehicle => (
            <Link
              key={vehicle.id}
              href={`/account/vehicles/${vehicle.id}`}
              className="block border border-border-subtle rounded-lg p-6 hover:bg-white/5 transition-colors"
            >
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-white">
                    {vehicle.year} {vehicle.make} {vehicle.model}
                  </h2>
                  <p className="text-sm text-text-muted mt-1">
                    {vehicle.colour}
                    {vehicle.rego && ` · ${vehicle.rego}`}
                    {vehicle.paintType && ` · ${vehicle.paintType}`}
                  </p>
                </div>
                {vehicle.isPrimary && (
                  <span className="text-xs bg-white/10 text-text-secondary px-2 py-0.5 rounded-full shrink-0">
                    Primary
                  </span>
                )}
              </div>
              <span className="inline-block mt-4 text-sm text-brand-400">
                View details &rarr;
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
