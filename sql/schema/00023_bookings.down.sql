DROP INDEX IF EXISTS idx_booking_service_options_booking_service_id;
DROP TABLE IF EXISTS booking_service_options;

DROP INDEX IF EXISTS idx_booking_services_booking_id;
DROP TABLE IF EXISTS booking_services;

DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_scheduled_date;
DROP INDEX IF EXISTS idx_bookings_vehicle_id;
DROP INDEX IF EXISTS idx_bookings_customer_id;
DROP TABLE IF EXISTS bookings;

DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS booking_status;
