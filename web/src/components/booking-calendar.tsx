'use client';

import { useState } from 'react';

interface BookingCalendarProps {
  selectedDate: string | null;
  onSelect: (date: string) => void;
}

function getDaysInMonth(year: number, month: number): number {
  return new Date(year, month + 1, 0).getDate();
}

function getFirstDayOfWeek(year: number, month: number): number {
  return new Date(year, month, 1).getDay();
}

function formatDateStr(year: number, month: number, day: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
}

const MONTH_NAMES = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
];

const DAY_LABELS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

export function BookingCalendar({ selectedDate, onSelect }: BookingCalendarProps) {
  const today = new Date();
  const [viewYear, setViewYear] = useState(today.getFullYear());
  const [viewMonth, setViewMonth] = useState(today.getMonth());

  const daysInMonth = getDaysInMonth(viewYear, viewMonth);
  const firstDay = getFirstDayOfWeek(viewYear, viewMonth);
  const todayStr = formatDateStr(today.getFullYear(), today.getMonth(), today.getDate());

  function prevMonth() {
    if (viewMonth === 0) {
      setViewMonth(11);
      setViewYear(viewYear - 1);
    } else {
      setViewMonth(viewMonth - 1);
    }
  }

  function nextMonth() {
    if (viewMonth === 11) {
      setViewMonth(0);
      setViewYear(viewYear + 1);
    } else {
      setViewMonth(viewMonth + 1);
    }
  }

  const cells: (number | null)[] = [];
  for (let i = 0; i < firstDay; i++) cells.push(null);
  for (let d = 1; d <= daysInMonth; d++) cells.push(d);

  return (
    <div className="w-full max-w-sm">
      <div className="flex items-center justify-between mb-4">
        <button
          type="button"
          onClick={prevMonth}
          className="px-2 py-1 text-text-secondary hover:text-foreground"
          aria-label="Previous month"
        >
          &larr;
        </button>
        <span className="font-semibold">
          {MONTH_NAMES[viewMonth]} {viewYear}
        </span>
        <button
          type="button"
          onClick={nextMonth}
          className="px-2 py-1 text-text-secondary hover:text-foreground"
          aria-label="Next month"
        >
          &rarr;
        </button>
      </div>
      <div className="grid grid-cols-7 gap-1 text-center text-xs text-text-muted mb-1">
        {DAY_LABELS.map(d => (
          <div key={d} className="py-1">{d}</div>
        ))}
      </div>
      <div className="grid grid-cols-7 gap-1">
        {cells.map((day, i) => {
          if (day === null) {
            return <div key={`empty-${i}`} />;
          }
          const dateStr = formatDateStr(viewYear, viewMonth, day);
          const isPast = dateStr <= todayStr;
          const isSelected = dateStr === selectedDate;

          return (
            <button
              key={dateStr}
              type="button"
              disabled={isPast}
              onClick={() => onSelect(dateStr)}
              className={`py-2 rounded text-sm ${
                isSelected
                  ? 'bg-brand-500 text-white font-semibold'
                  : isPast
                    ? 'text-foreground/20 cursor-not-allowed'
                    : 'hover:bg-surface-hover text-text-secondary'
              }`}
            >
              {day}
            </button>
          );
        })}
      </div>
    </div>
  );
}
