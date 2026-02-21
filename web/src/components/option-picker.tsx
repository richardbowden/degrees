'use client';

import { DetailingServiceOption } from '@/lib/types';
import { formatPrice } from '@/lib/format';

interface OptionPickerProps {
  options: DetailingServiceOption[];
  selectedIds: string[];
  onChange: (selectedIds: string[]) => void;
}

export function OptionPicker({ options, selectedIds, onChange }: OptionPickerProps) {
  function toggle(id: string) {
    if (selectedIds.includes(id)) {
      onChange(selectedIds.filter(s => s !== id));
    } else {
      onChange([...selectedIds, id]);
    }
  }

  if (options.length === 0) return null;

  return (
    <div className="space-y-3">
      <h4 className="font-semibold text-sm text-gray-700">Add-ons</h4>
      {options.filter(o => o.isActive).map(option => (
        <label
          key={option.id}
          className="flex items-start gap-3 p-3 border border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50"
        >
          <input
            type="checkbox"
            checked={selectedIds.includes(option.id)}
            onChange={() => toggle(option.id)}
            className="mt-1"
          />
          <div className="flex-1">
            <div className="flex items-center justify-between">
              <span className="font-medium text-sm">{option.name}</span>
              <span className="text-sm text-gray-600">+{formatPrice(option.price)}</span>
            </div>
            {option.description && (
              <p className="text-xs text-gray-500 mt-1">{option.description}</p>
            )}
          </div>
        </label>
      ))}
    </div>
  );
}
