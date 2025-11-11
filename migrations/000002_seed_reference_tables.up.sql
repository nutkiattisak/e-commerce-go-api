-- ===================================
-- Migration: Seed Reference Tables
-- Version: 000002
-- Description: Insert master data for all reference tables
-- ===================================

BEGIN;

INSERT INTO payment_status (id, code, name) VALUES
  (1, 'PENDING_PAYMENT', 'รอชำระเงิน'),
  (2, 'PROCESSING', 'กำลังดำเนินการ'),
  (3, 'PAID', 'ชำระเงินแล้ว'),
  (4, 'FAILED', 'ชำระเงินไม่สำเร็จ'),
  (5, 'CANCELLED', 'ยกเลิก'),
  (6, 'REFUNDED', 'คืนเงินแล้ว'),
  (7, 'EXPIRED', 'หมดอายุ')
ON CONFLICT (id) DO NOTHING;

INSERT INTO payment_methods (id, code, name) VALUES
  (1, 'CREDIT_CARD', 'บัตรเครดิต'),
  (2, 'COD', 'เก็บเงินปลายทาง'),
  (3, 'BANK_TRANSFER', 'โอนเงินผ่านธนาคาร'),
  (4, 'PROMPTPAY', 'พร้อมเพย์')
ON CONFLICT (id) DO NOTHING;

INSERT INTO order_status (id, code, name) VALUES
  (1, 'PENDING', 'รอดำเนินการ'),
  (2, 'PROCESSING', 'กำลังเตรียมสินค้า'),
  (3, 'SHIPPED', 'จัดส่งแล้ว'),
  (4, 'DELIVERED', 'ส่งสำเร็จ'),
  (5, 'COMPLETED', 'เสร็จสิ้น'),
  (6, 'CANCELLED', 'ยกเลิก')
ON CONFLICT (id) DO NOTHING;

INSERT INTO refund_status (id, code, name) VALUES
  (1, 'PENDING', 'รอดำเนินการ'),
  (2, 'APPROVED', 'อนุมัติแล้ว'),
  (3, 'COMPLETED', 'คืนเงินแล้ว'),
  (4, 'REJECTED', 'ปฏิเสธ')
ON CONFLICT (id) DO NOTHING;

-- Refund Methods
INSERT INTO refund_methods (id, code, name) VALUES
  (1, 'BANK_TRANSFER', 'โอนเงินเข้าบัญชี'),
  (2, 'CREDIT_CARD', 'คืนเข้าบัตรเครดิต')
ON CONFLICT (id) DO NOTHING;

-- Shipment Status
INSERT INTO shipment_status (id, code, name) VALUES
  (1, 'IN_TRANSIT', 'กำลังจัดส่ง'),
  (2, 'DELIVERED', 'จัดส่งสำเร็จ'),
  (3, 'FAILED_DELIVERY', 'จัดส่งไม่สำเร็จ')
ON CONFLICT (id) DO NOTHING;

-- Roles
INSERT INTO roles (name) VALUES
  ('ADMIN'),
  ('USER'),
  ('SHOP')
ON CONFLICT (name) DO NOTHING;

-- Couriers
INSERT INTO couriers (name, image_url, rate) VALUES
  ('Kerry Express', NULL, 50.00),
  ('Flash Express', NULL, 45.00),
  ('Thailand Post', NULL, 40.00),
  ('J&T Express', NULL, 48.00),
  ('SCG Express', NULL, 52.00)
ON CONFLICT DO NOTHING;

DO $$
DECLARE
  payment_status_count INT;
  payment_method_count INT;
  order_status_count INT;
  refund_status_count INT;
  refund_method_count INT;
  shipment_status_count INT;
  role_count INT;
  courier_count INT;
BEGIN
  SELECT COUNT(*) INTO payment_status_count FROM payment_status;
  SELECT COUNT(*) INTO payment_method_count FROM payment_methods;
  SELECT COUNT(*) INTO order_status_count FROM order_status;
  SELECT COUNT(*) INTO refund_status_count FROM refund_status;
  SELECT COUNT(*) INTO refund_method_count FROM refund_methods;
  SELECT COUNT(*) INTO shipment_status_count FROM shipment_status;
  SELECT COUNT(*) INTO role_count FROM roles;
  SELECT COUNT(*) INTO courier_count FROM couriers;

  RAISE NOTICE '==========================================';
  RAISE NOTICE 'Seed Data Summary:';
  RAISE NOTICE '==========================================';
  RAISE NOTICE 'Payment Status: % records', payment_status_count;
  RAISE NOTICE 'Payment Methods: % records', payment_method_count;
  RAISE NOTICE 'Order Status: % records', order_status_count;
  RAISE NOTICE 'Refund Status: % records', refund_status_count;
  RAISE NOTICE 'Refund Methods: % records', refund_method_count;
  RAISE NOTICE 'Shipment Status: % records', shipment_status_count;
  RAISE NOTICE 'Roles: % records', role_count;
  RAISE NOTICE 'Couriers: % records', courier_count;
  RAISE NOTICE '==========================================';
  RAISE NOTICE 'Seed data inserted successfully!';
END $$;

COMMIT;
