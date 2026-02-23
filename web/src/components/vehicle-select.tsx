'use client';

import { Vehicle } from '@/lib/types';

interface VehicleSelectProps {
  vehicles: Vehicle[];
  selectedId: string;
  onChange: (vehicleId: string) => void;
}

export function VehicleSelect({ vehicles, selectedId, onChange }: VehicleSelectProps) {
  if (vehicles.length === 0) {
    return <p className="text-sm text-text-muted">No vehicles found. Add a vehicle first.</p>;
  }

  return (
    <select
      value={selectedId}
      onChange={e => onChange(e.target.value)}
      className="w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500"
    >
      <option value="">Select a vehicle</option>
      {vehicles.map(v => (
        <option key={v.id} value={v.id}>
          {v.year} {v.make} {v.model} {v.colour ? `(${v.colour})` : ''} {v.rego ? `- ${v.rego}` : ''}
        </option>
      ))}
    </select>
  );
}
