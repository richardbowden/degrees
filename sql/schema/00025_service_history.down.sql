DROP INDEX IF EXISTS idx_service_photos_service_record_id;
DROP TABLE IF EXISTS service_photos;

DROP INDEX IF EXISTS idx_service_products_used_service_record_id;
DROP TABLE IF EXISTS service_products_used;

DROP INDEX IF EXISTS idx_service_notes_service_record_id;
DROP TABLE IF EXISTS service_notes;

DROP INDEX IF EXISTS idx_service_records_vehicle_id;
DROP INDEX IF EXISTS idx_service_records_customer_id;
DROP INDEX IF EXISTS idx_service_records_booking_id;
DROP TABLE IF EXISTS service_records;
