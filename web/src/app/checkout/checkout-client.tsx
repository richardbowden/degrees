'use client';

import { useState, useCallback, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api, ApiError } from '@/lib/api';
import { VehicleSelect } from '@/components/vehicle-select';
import { BookingCalendar } from '@/components/booking-calendar';
import { TimeSlotGrid } from '@/components/time-slot-grid';
import { VehicleForm } from '@/components/vehicle-form';
import { formatPrice, formatDate, formatTime } from '@/lib/format';
import type { Vehicle, AvailableSlot, Cart, VehicleFormData } from '@/lib/types';
import { addVehicleAction, submitBookingAction } from './checkout-actions';

type Step = 'vehicle' | 'date' | 'time' | 'confirm';

interface Props {
  initialVehicles: Vehicle[];
  initialCart: Cart | null;
  serviceDurations: Record<string, number>;
}

export function CheckoutClient({ initialVehicles, initialCart, serviceDurations }: Props) {
  const router = useRouter();
  const [step, setStep] = useState<Step>('vehicle');
  const [vehicles, setVehicles] = useState<Vehicle[]>(initialVehicles);
  const [cart, setCart] = useState<Cart | null>(initialCart);
  // If server-side returned an empty cart (can happen when the guest cart wasn't merged),
  // do a client-side fetch which correctly reads the cart session from document.cookie.
  const [cartLoading, setCartLoading] = useState(!initialCart?.items?.length);

  useEffect(() => {
    if (!initialCart?.items?.length) {
      api<{ cart: Cart }>('/cart')
        .then(data => { if (data.cart) setCart(data.cart); })
        .catch(() => {})
        .finally(() => setCartLoading(false));
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  const [selectedVehicle, setSelectedVehicle] = useState('');
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [notes, setNotes] = useState('');
  const [slots, setSlots] = useState<AvailableSlot[]>([]);
  const [slotsLoading, setSlotsLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [showAddVehicle, setShowAddVehicle] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const items = cart?.items ?? [];
  const subtotal = Number(cart?.subtotal ?? 0);
  const depositAmount = Math.floor(subtotal * 30 / 100);

  const totalDuration = items.reduce((sum, item) => {
    const duration = serviceDurations[item.serviceId] || 60;
    return sum + duration * Number(item.quantity);
  }, 0) || 60;

  const selectedVehicleObj = vehicles.find(v => v.id === selectedVehicle) ?? null;

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

  async function handleAddVehicle(data: VehicleFormData) {
    setError(null);
    const result = await addVehicleAction(data);
    if (result.error) {
      setError(result.error);
      return;
    }
    if (result.vehicle) {
      setVehicles(prev => [...prev, result.vehicle!]);
      setShowAddVehicle(false);
      setSelectedVehicle(result.vehicle.id);
      // Don't auto-advance — let user confirm their vehicle choice
    }
  }

  function handleVehicleContinue() {
    if (!selectedVehicle) {
      setError('Please select a vehicle to continue.');
      return;
    }
    setError(null);
    setStep('date');
  }

  function handleDateSelect(date: string) {
    setSelectedDate(date);
    setStep('time');
    fetchSlots(date);
  }

  function handleTimeSelect(time: string) {
    setSelectedTime(time);
    setStep('confirm');
  }

  async function handleSubmit() {
    if (!selectedVehicle || !selectedDate || !selectedTime) return;
    setSubmitting(true);
    setError(null);
    const result = await submitBookingAction({
      vehicleId: selectedVehicle,
      scheduledDate: selectedDate,
      scheduledTime: selectedTime,
      notes,
      cartSession: cart?.sessionToken ?? '',
    });
    if (result.error) {
      setError(result.error);
      setSubmitting(false);
      return;
    }
    if (result.bookingId) {
      router.push(`/checkout/deposit?booking_id=${result.bookingId}`);
    }
  }

  const steps: { key: Step; label: string; number: number }[] = [
    { key: 'vehicle', label: 'Vehicle', number: 1 },
    { key: 'date', label: 'Date', number: 2 },
    { key: 'time', label: 'Time', number: 3 },
    { key: 'confirm', label: 'Confirm', number: 4 },
  ];

  const stepIndex = steps.findIndex(s => s.key === step);

  if (cartLoading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-text-muted">Loading checkout...</p>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-4">Checkout</h1>
        <p className="text-text-muted mb-6">Your cart is empty.</p>
        <a href="/services" className="text-brand-400 font-medium hover:underline">Browse services</a>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">Checkout</h1>

      {/* Step indicator */}
      <div className="flex items-center gap-2 mb-8">
        {steps.map((s, i) => (
          <div key={s.key} className="flex items-center">
            <button
              type="button"
              onClick={() => i < stepIndex && setStep(s.key)}
              disabled={i >= stepIndex}
              className={`flex items-center gap-2 text-sm font-medium ${
                i < stepIndex ? 'text-foreground cursor-pointer' : i === stepIndex ? 'text-foreground' : 'text-text-muted'
              }`}
            >
              <span
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs ${
                  i < stepIndex
                    ? 'bg-green-500 text-white'
                    : i === stepIndex
                      ? 'bg-brand-500 text-white'
                      : 'bg-surface-raised text-text-muted'
                }`}
              >
                {i < stepIndex ? '\u2713' : s.number}
              </span>
              <span className="hidden sm:inline">{s.label}</span>
            </button>
            {i < steps.length - 1 && (
              <div className="w-8 h-px bg-border-subtle mx-2" />
            )}
          </div>
        ))}
      </div>

      {error && (
        <p className="text-red-400 text-sm mb-4">{error}</p>
      )}

      {/* Step 1: Vehicle */}
      {step === 'vehicle' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Select Your Vehicle</h2>
          {showAddVehicle ? (
            <VehicleForm
              onSubmit={handleAddVehicle}
              onCancel={() => setShowAddVehicle(false)}
            />
          ) : (
            <>
              <VehicleSelect
                vehicles={vehicles}
                selectedId={selectedVehicle}
                onChange={id => { setSelectedVehicle(id); setError(null); }}
              />
              <button
                type="button"
                onClick={() => setShowAddVehicle(true)}
                className="mt-3 text-sm font-medium text-brand-400 underline hover:text-brand-500"
              >
                + Add a new vehicle
              </button>
              <div className="mt-6">
                <button
                  type="button"
                  onClick={handleVehicleContinue}
                  disabled={!selectedVehicle}
                  className="w-full btn-brand py-3 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Continue
                </button>
              </div>
            </>
          )}
        </div>
      )}

      {/* Step 2: Date */}
      {step === 'date' && (
        <div>
          <h2 className="text-lg font-semibold mb-1">Pick a Date</h2>
          <p className="text-sm text-text-muted mb-4">Bookings require at least 48 hours notice.</p>
          <BookingCalendar
            selectedDate={selectedDate}
            onSelect={handleDateSelect}
          />
        </div>
      )}

      {/* Step 3: Time */}
      {step === 'time' && (
        <div>
          <h2 className="text-lg font-semibold mb-1">Pick a Time</h2>
          {selectedDate && (
            <p className="text-sm text-text-muted mb-4">{formatDate(selectedDate)}</p>
          )}
          {slotsLoading ? (
            <p className="text-text-muted text-sm">Loading available slots...</p>
          ) : (
            <TimeSlotGrid
              slots={slots}
              selectedTime={selectedTime}
              onSelect={handleTimeSelect}
            />
          )}
        </div>
      )}

      {/* Step 4: Confirm */}
      {step === 'confirm' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Confirm Your Booking</h2>

          {/* Booking summary */}
          <div className="border border-border-subtle rounded-lg divide-y divide-border-subtle mb-6">
            {/* Services */}
            <div className="p-4">
              <p className="text-xs text-text-muted uppercase tracking-wide mb-2">Services</p>
              <div className="space-y-1">
                {items.map(item => (
                  <div key={item.id} className="flex justify-between text-sm">
                    <span className="text-foreground">
                      {item.serviceName}
                      {Number(item.quantity) > 1 && <span className="text-text-muted"> ×{item.quantity}</span>}
                    </span>
                    <span className="text-text-secondary">{formatPrice(Number(item.servicePrice) * Number(item.quantity))}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Vehicle */}
            {selectedVehicleObj && (
              <div className="p-4">
                <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Vehicle</p>
                <p className="text-sm text-foreground">
                  {selectedVehicleObj.year} {selectedVehicleObj.make} {selectedVehicleObj.model}
                  {selectedVehicleObj.colour && ` (${selectedVehicleObj.colour})`}
                  {selectedVehicleObj.rego && ` — ${selectedVehicleObj.rego}`}
                </p>
              </div>
            )}

            {/* Date & Time */}
            <div className="p-4 grid grid-cols-2 gap-4">
              <div>
                <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Date</p>
                <p className="text-sm text-foreground">{selectedDate ? formatDate(selectedDate) : '—'}</p>
              </div>
              <div>
                <p className="text-xs text-text-muted uppercase tracking-wide mb-1">Time</p>
                <p className="text-sm text-foreground">{selectedTime ? formatTime(selectedTime) : '—'}</p>
              </div>
            </div>

            {/* Totals */}
            <div className="p-4 space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-text-secondary">Subtotal</span>
                <span>{formatPrice(subtotal)}</span>
              </div>
              <div className="flex justify-between text-sm font-medium text-brand-400">
                <span>Deposit (30%, due on arrival)</span>
                <span>{formatPrice(depositAmount)}</span>
              </div>
              <div className="flex justify-between text-sm font-semibold pt-1 border-t border-border-subtle mt-1">
                <span>Total</span>
                <span>{formatPrice(subtotal)}</span>
              </div>
            </div>
          </div>

          {/* Optional notes */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-text-secondary mb-1">
              Special requests <span className="text-text-muted font-normal">(optional)</span>
            </label>
            <textarea
              value={notes}
              onChange={e => setNotes(e.target.value)}
              rows={3}
              placeholder="Parking instructions, specific concerns, etc."
              className="w-full bg-surface-input border border-border-subtle rounded px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </div>

          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="w-full btn-brand py-3"
          >
            {submitting ? 'Creating booking...' : 'Confirm Booking'}
          </button>
        </div>
      )}
    </div>
  );
}
