'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { api, ApiError } from '@/lib/api';
import { VehicleSelect } from '@/components/vehicle-select';
import { BookingCalendar } from '@/components/booking-calendar';
import { TimeSlotGrid } from '@/components/time-slot-grid';
import type { Vehicle, AvailableSlot, Booking, Cart } from '@/lib/types';

type Step = 'vehicle' | 'date' | 'time' | 'notes';

export default function CheckoutPage() {
  const router = useRouter();
  const [step, setStep] = useState<Step>('vehicle');
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [cart, setCart] = useState<Cart | null>(null);
  const [selectedVehicle, setSelectedVehicle] = useState('');
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [notes, setNotes] = useState('');
  const [slots, setSlots] = useState<AvailableSlot[]>([]);
  const [loading, setLoading] = useState(true);
  const [slotsLoading, setSlotsLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const [vehiclesData, cartData] = await Promise.all([
          api<{ vehicles: Vehicle[] }>('/me/vehicles'),
          api<{ cart: Cart }>('/cart'),
        ]);
        setVehicles(vehiclesData.vehicles ?? []);
        setCart(cartData.cart);
      } catch (err: unknown) {
        const apiErr = err as ApiError;
        setError(apiErr.detail || 'Failed to load checkout data');
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const totalDuration = cart?.items?.reduce((sum, item) => sum + 60 * Number(item.quantity), 0) ?? 60;

  const fetchSlots = useCallback(async (date: string) => {
    setSlotsLoading(true);
    setSlots([]);
    setSelectedTime(null);
    try {
      const data = await api<{ slots: AvailableSlot[] }>(
        `/checkout/available-slots?date=${date}&durationMinutes=${totalDuration}`
      );
      setSlots(data.slots ?? []);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to load time slots');
    } finally {
      setSlotsLoading(false);
    }
  }, [totalDuration]);

  function handleVehicleSelect(vehicleId: string) {
    setSelectedVehicle(vehicleId);
    setStep('date');
  }

  function handleDateSelect(date: string) {
    setSelectedDate(date);
    setStep('time');
    fetchSlots(date);
  }

  function handleTimeSelect(time: string) {
    setSelectedTime(time);
    setStep('notes');
  }

  async function handleSubmit() {
    if (!selectedVehicle || !selectedDate || !selectedTime) return;
    setSubmitting(true);
    setError(null);
    try {
      const data = await api<{ booking: Booking }>('/checkout', {
        method: 'POST',
        body: {
          vehicleId: selectedVehicle,
          scheduledDate: selectedDate,
          scheduledTime: selectedTime,
          notes,
        },
      });
      router.push(`/checkout/deposit?booking_id=${data.booking.id}`);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to create booking');
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-gray-500">Loading checkout...</p>
      </div>
    );
  }

  const steps: { key: Step; label: string; number: number }[] = [
    { key: 'vehicle', label: 'Vehicle', number: 1 },
    { key: 'date', label: 'Date', number: 2 },
    { key: 'time', label: 'Time', number: 3 },
    { key: 'notes', label: 'Confirm', number: 4 },
  ];

  const stepIndex = steps.findIndex(s => s.key === step);

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">Checkout</h1>

      <div className="flex items-center gap-2 mb-8">
        {steps.map((s, i) => (
          <div key={s.key} className="flex items-center">
            <button
              type="button"
              onClick={() => i <= stepIndex && setStep(s.key)}
              disabled={i > stepIndex}
              className={`flex items-center gap-2 text-sm font-medium ${
                i <= stepIndex ? 'text-gray-900' : 'text-gray-400'
              }`}
            >
              <span
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs ${
                  i < stepIndex
                    ? 'bg-green-600 text-white'
                    : i === stepIndex
                      ? 'bg-gray-900 text-white'
                      : 'bg-gray-200 text-gray-500'
                }`}
              >
                {i < stepIndex ? '\u2713' : s.number}
              </span>
              <span className="hidden sm:inline">{s.label}</span>
            </button>
            {i < steps.length - 1 && (
              <div className="w-8 h-px bg-gray-300 mx-2" />
            )}
          </div>
        ))}
      </div>

      {error && (
        <p className="text-red-600 text-sm mb-4">{error}</p>
      )}

      {step === 'vehicle' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Select Vehicle</h2>
          <VehicleSelect
            vehicles={vehicles}
            selectedId={selectedVehicle}
            onChange={handleVehicleSelect}
          />
        </div>
      )}

      {step === 'date' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Pick a Date</h2>
          <BookingCalendar
            selectedDate={selectedDate}
            onSelect={handleDateSelect}
          />
        </div>
      )}

      {step === 'time' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Pick a Time</h2>
          {slotsLoading ? (
            <p className="text-gray-500 text-sm">Loading available slots...</p>
          ) : (
            <TimeSlotGrid
              slots={slots}
              selectedTime={selectedTime}
              onSelect={handleTimeSelect}
            />
          )}
        </div>
      )}

      {step === 'notes' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Any Special Requests?</h2>
          <textarea
            value={notes}
            onChange={e => setNotes(e.target.value)}
            rows={4}
            placeholder="Parking instructions, specific concerns, etc."
            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900 mb-6"
          />
          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="w-full bg-gray-900 text-white py-3 rounded font-semibold hover:bg-gray-800 disabled:opacity-50"
          >
            {submitting ? 'Creating booking...' : 'Create Booking'}
          </button>
        </div>
      )}
    </div>
  );
}
