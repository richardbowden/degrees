'use client';

import { useState } from 'react';
import { Vehicle } from '@/lib/types';

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

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    onSubmit({ make, model, year, colour, rego, paintType, conditionNotes, isPrimary });
  }

  const inputClass = 'w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1';

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
      <label className="flex items-center gap-2 text-sm text-gray-700">
        <input
          type="checkbox"
          checked={isPrimary}
          onChange={e => setIsPrimary(e.target.checked)}
          className="rounded border-gray-300"
        />
        Primary vehicle
      </label>
      <div className="flex gap-3">
        <button
          type="submit"
          className="bg-gray-900 text-white px-4 py-2 rounded text-sm font-semibold hover:bg-gray-800"
        >
          {vehicle ? 'Update Vehicle' : 'Add Vehicle'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="border border-gray-300 px-4 py-2 rounded text-sm text-gray-700 hover:bg-gray-50"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
