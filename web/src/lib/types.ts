export interface User {
  id: string;
  firstName: string;
  middleName: string;
  surname: string;
  email: string;
  enabled: boolean;
  sysop: boolean;
  signUpStage: number;
  createdOn: string;
  updatedAt: string;
}

export interface ServiceCategory {
  id: string;
  name: string;
  slug: string;
  description: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface DetailingServiceOption {
  id: string;
  serviceId: string;
  name: string;
  description: string;
  price: number;
  isActive: boolean;
  sortOrder: number;
  createdAt: string;
}

export interface DetailingService {
  id: string;
  categoryId: string;
  name: string;
  slug: string;
  description: string;
  shortDesc: string;
  basePrice: number;
  durationMinutes: number;
  isActive: boolean;
  sortOrder: number;
  categoryName: string;
  options: DetailingServiceOption[];
  createdAt: string;
  updatedAt: string;
}

export interface CartItem {
  id: string;
  serviceId: string;
  vehicleId: string;
  quantity: number;
  serviceName: string;
  servicePrice: number;
  optionIds: string[];
  createdAt: string;
}

export interface Cart {
  id: string;
  sessionToken: string;
  items: CartItem[];
  subtotal: number;
  expiresAt: string;
}

export interface AvailableSlot {
  date: string;
  time: string;
  availableDurationMins: number;
}

export interface BookingServiceItem {
  serviceId: string;
  serviceName: string;
  price: number;
  options: string[];
}

export interface BookingCustomerInfo {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface BookingVehicleInfo {
  id: string;
  make: string;
  model: string;
  year: number;
  colour: string;
  rego: string;
}

export interface Booking {
  id: string;
  customerId: string;
  vehicleId: string;
  scheduledDate: string;
  scheduledTime: string;
  estimatedDurationMins: number;
  status: string;
  paymentStatus: string;
  subtotal: number;
  depositAmount: number;
  totalAmount: number;
  notes: string;
  services: BookingServiceItem[];
  customer: BookingCustomerInfo;
  vehicle: BookingVehicleInfo;
  createdAt: string;
  updatedAt: string;
}

export interface CustomerProfile {
  id: string;
  userId: string;
  phone: string;
  address: string;
  suburb: string;
  postcode: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

export interface Vehicle {
  id: string;
  customerId: string;
  make: string;
  model: string;
  year: number;
  colour: string;
  rego: string;
  paintType: string;
  conditionNotes: string;
  isPrimary: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ServiceNote {
  id: string;
  serviceRecordId: string;
  noteType: string;
  content: string;
  isVisibleToCustomer: boolean;
  createdBy: string;
  createdAt: string;
}

export interface ProductUsed {
  id: string;
  serviceRecordId: string;
  productName: string;
  notes: string;
  createdAt: string;
}

export interface ServiceRecord {
  id: string;
  bookingId: string;
  customerId: string;
  vehicleId: string;
  completedDate: string;
  notes: ServiceNote[];
  products: ProductUsed[];
  createdAt: string;
  updatedAt: string;
}

export interface ScheduleDay {
  id: string;
  dayOfWeek: number;
  openTime: string;
  closeTime: string;
  isOpen: boolean;
  bufferMinutes: number;
}

export interface Blackout {
  id: string;
  date: string;
  reason: string;
}
