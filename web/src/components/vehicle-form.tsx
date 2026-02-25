'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Vehicle, VehicleCategory } from '@/lib/types';

interface VehicleFormProps {
  vehicle?: Vehicle;
  onSubmit: (data: VehicleFormData) => void;
  onCancel: () => void;
}

export interface VehicleFormData {
  make: string;
  model: string;
  year: number;
  colour: string;
  rego: string;
  paintType: string;
  conditionNotes: string;
  isPrimary: boolean;
  vehicleCategoryId: string;
}

export function VehicleForm({ vehicle, onSubmit, onCancel }: VehicleFormProps) {
  const [make, setMake] = useState(vehicle?.make ?? '');
  const [model, setModel] = useState(vehicle?.model ?? '');
  const [year, setYear] = useState(vehicle?.year ?? new Date().getFullYear());
  const [colour, setColour] = useState(vehicle?.colour ?? '');
  const [rego, setRego] = useState(vehicle?.rego ?? '');
  const [paintType, setPaintType] = useState(vehicle?.paintType ?? '');
  const [conditionNotes, setConditionNotes] = useState(vehicle?.conditionNotes ?? '');
  const [isPrimary, setIsPrimary] = useState(vehicle?.isPrimary ?? false);
  const [vehicleCategoryId, setVehicleCategoryId] = useState(vehicle?.vehicleCategoryId ?? '');
  const [vehicleCategories, setVehicleCategories] = useState<VehicleCategory[]>([]);

  useEffect(() => {
    api<{ vehicleCategories: VehicleCategory[] }>('/catalogue/vehicle-categories')
      .then(res => setVehicleCategories(res.vehicleCategories ?? []))
      .catch(() => {});
  }, []);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    onSubmit({ make, model, year, colour, rego, paintType, conditionNotes, isPrimary, vehicleCategoryId });
  }

  const inputClass = 'w-full bg-white/5 border border-border-subtle rounded px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500';
  const labelClass = 'block text-sm font-medium text-text-secondary mb-1';

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label className={labelClass}>Make</label>
          <input type="text" required value={make} onChange={e => setMake(e.target.value)} className={inputClass} placeholder="Toyota" />
        </div>
        <div>
          <label className={labelClass}>Model</label>
          <input type="text" required value={model} onChange={e => setModel(e.target.value)} className={inputClass} placeholder="Camry" />
        </div>
        <div>
          <label className={labelClass}>Year</label>
          <input type="number" required min={1900} max={2100} value={year} onChange={e => setYear(Number(e.target.value))} className={inputClass} />
        </div>
        <div>
          <label className={labelClass}>Colour</label>
          <input type="text" required value={colour} onChange={e => setColour(e.target.value)} className={inputClass} placeholder="White" />
        </div>
        <div>
          <label className={labelClass}>Registration</label>
          <input type="text" value={rego} onChange={e => setRego(e.target.value)} className={inputClass} placeholder="1ABC234" />
        </div>
        <div>
          <label className={labelClass}>Paint Type</label>
          <input type="text" value={paintType} onChange={e => setPaintType(e.target.value)} className={inputClass} placeholder="Metallic" />
        </div>
        {vehicleCategories.length > 0 && (
          <div>
            <label className={labelClass}>Vehicle Size</label>
            <select
              required
              value={vehicleCategoryId}
              onChange={e => setVehicleCategoryId(e.target.value)}
              className={inputClass}
            >
              <option value="">Select size...</option>
              {vehicleCategories.map(vc => (
                <option key={vc.id} value={vc.id}>{vc.name}</option>
              ))}
            </select>
          </div>
        )}
      </div>
      <div>
        <label className={labelClass}>Condition Notes</label>
        <textarea
          value={conditionNotes}
          onChange={e => setConditionNotes(e.target.value)}
          className={inputClass}
          rows={3}
          placeholder="Any scratches, dents, or special considerations..."
        />
      </div>
      <label className="flex items-center gap-2 text-sm text-text-secondary">
        <input
          type="checkbox"
          checked={isPrimary}
          onChange={e => setIsPrimary(e.target.checked)}
          className="rounded border-border-subtle accent-brand-500"
        />
        Primary vehicle
      </label>
      <div className="flex gap-3">
        <button
          type="submit"
          className="btn-brand px-4 py-2 text-sm"
        >
          {vehicle ? 'Update Vehicle' : 'Add Vehicle'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="border border-border-subtle px-4 py-2 rounded text-sm text-text-secondary hover:bg-white/5"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
