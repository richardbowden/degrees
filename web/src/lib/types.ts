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
  priceTiers: ServicePriceTier[];
}

export interface VehicleCategory {
  id: string;
  name: string;
  slug: string;
  description: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface ServicePriceTier {
  serviceId: string;
  vehicleCategoryId: string;
  categoryName: string;
  categorySlug: string;
  price: number;
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

export interface BookingServiceOptionItem {
  id: string;
  serviceOptionId: string;
  optionName: string;
  priceAtBooking: number;
}

export interface BookingServiceItem {
  id: string;
  serviceId: string;
  serviceName: string;
  serviceSlug: string;
  priceAtBooking: number;
  options: BookingServiceOptionItem[];
}

export interface BookingCustomerInfo {
  userId: string;
  phone: string;
  name: string;
}

export interface BookingVehicleInfo {
  make: string;
  model: string;
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
  vehicleCategoryId: string;
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

export interface ServicePhoto {
  id: string;
  serviceRecordId: string;
  photoType: string;
  url: string;
  caption: string;
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
  photos: ServicePhoto[];
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
