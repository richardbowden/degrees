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
      <h4 className="font-semibold text-sm text-text-secondary">Add-ons</h4>
      {options.filter(o => o.isActive).map(option => (
        <label
          key={option.id}
          className={`flex items-start gap-3 p-3 border rounded-lg cursor-pointer transition-colors ${
            selectedIds.includes(option.id)
              ? 'border-brand-500 bg-brand-500/10'
              : 'border-border-subtle hover:bg-surface-hover'
          }`}
        >
          <input
            type="checkbox"
            checked={selectedIds.includes(option.id)}
            onChange={() => toggle(option.id)}
            className="mt-1 accent-brand-500"
          />
          <div className="flex-1">
            <div className="flex items-center justify-between">
              <span className="font-medium text-sm">{option.name}</span>
              <span className="text-sm text-brand-400">+{formatPrice(option.price)}</span>
            </div>
            {option.description && (
              <p className="text-xs text-text-muted mt-1">{option.description}</p>
            )}
          </div>
        </label>
      ))}
    </div>
  );
}
