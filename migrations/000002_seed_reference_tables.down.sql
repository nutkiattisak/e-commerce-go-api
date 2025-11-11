-- ===================================
-- Rollback: Clear Reference Tables
-- Version: 000002
-- ===================================

BEGIN;

DELETE FROM couriers;
DELETE FROM roles;
DELETE FROM shipment_status;
DELETE FROM refund_methods;
DELETE FROM refund_status;
DELETE FROM order_status;
DELETE FROM payment_methods;
DELETE FROM payment_status;

COMMIT;
