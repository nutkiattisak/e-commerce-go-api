-- ===================================
-- Migration: Create Initial Tables
-- Version: 000001
-- Description: Create all database tables for E-commerce API
-- ===================================

BEGIN;

-- ============================================
-- 1. Reference Tables (Master Data)
-- ============================================

-- Payment Status
CREATE TABLE IF NOT EXISTS payment_status (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Payment Methods
CREATE TABLE IF NOT EXISTS payment_methods (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Order Status
CREATE TABLE IF NOT EXISTS order_status (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Refund Status
CREATE TABLE IF NOT EXISTS refund_status (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Refund Methods
CREATE TABLE IF NOT EXISTS refund_methods (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Shipment Status
CREATE TABLE IF NOT EXISTS shipment_status (
    id INTEGER NOT NULL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

-- Roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- ============================================
-- 2. Location Tables
-- ============================================

-- Provinces
CREATE TABLE IF NOT EXISTS provinces (
    id INTEGER NOT NULL PRIMARY KEY,
    name_th VARCHAR(150) NOT NULL,
    name_en VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6)
);

-- Districts
CREATE TABLE IF NOT EXISTS districts (
    id SERIAL PRIMARY KEY,
    province_id INTEGER NOT NULL,
    name_th VARCHAR(150) NOT NULL,
    name_en VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_districts_province_id ON districts(province_id);

-- Sub Districts
CREATE TABLE IF NOT EXISTS sub_districts (
    id SERIAL PRIMARY KEY,
    zipcode INTEGER NOT NULL,
    name_th VARCHAR(150) NOT NULL,
    name_en VARCHAR(150) NOT NULL,
    district_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sub_districts_district_id ON sub_districts(district_id);
CREATE INDEX IF NOT EXISTS idx_sub_districts_zipcode ON sub_districts(zipcode);

-- ============================================
-- 3. User Management Tables
-- ============================================

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    image_url TEXT,
    created_at TIMESTAMPTZ(6) DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;

-- User Roles (Many-to-Many)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ(6) DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Addresses
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    name TEXT,
    line1 TEXT,
    line2 TEXT,
    sub_district_id INTEGER NOT NULL,
    district_id INTEGER NOT NULL,
    province_id INTEGER NOT NULL,
    zipcode INTEGER,
    phone_number VARCHAR(15),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sub_district_id) REFERENCES sub_districts(id),
    FOREIGN KEY (district_id) REFERENCES districts(id),
    FOREIGN KEY (province_id) REFERENCES provinces(id)
);

CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id) WHERE deleted_at IS NULL;

-- ============================================
-- 4. Shop Management Tables
-- ============================================

-- Shops
CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    address TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_shops_user_id ON shops(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_shops_is_active ON shops(is_active) WHERE deleted_at IS NULL;

-- ============================================
-- 5. Product Management Tables
-- ============================================

-- Products
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    price NUMERIC(10, 2) NOT NULL,
    stock_qty INTEGER NOT NULL DEFAULT 0,
    shop_id UUID NOT NULL,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ(6) DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    CONSTRAINT products_price_positive CHECK (price >= 0),
    CONSTRAINT products_stock_positive CHECK (stock_qty >= 0)
);

CREATE INDEX IF NOT EXISTS idx_products_shop_id ON products(shop_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active) WHERE deleted_at IS NULL;

-- ============================================
-- 6. Cart Tables
-- ============================================

-- Carts
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Cart Items
CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT cart_items_quantity_positive CHECK (qty > 0),
    UNIQUE(cart_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);

-- ============================================
-- 7. Order Tables
-- ============================================

-- Orders (Main Order - ครอบคลุมหลายร้าน)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    address_id INTEGER,
    grand_total NUMERIC(10, 2) NOT NULL,
    shipping_name VARCHAR(255) NOT NULL,
    shipping_phone VARCHAR(15) NOT NULL,
    shipping_line1 TEXT NOT NULL,
    shipping_line2 TEXT,
    shipping_sub_district VARCHAR(100) NOT NULL,
    shipping_district VARCHAR(100) NOT NULL,
    shipping_province VARCHAR(100) NOT NULL,
    shipping_zipcode VARCHAR(5) NOT NULL,
    payment_method_id INTEGER,
    created_at TIMESTAMPTZ(6),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT orders_total_price_positive CHECK (grand_total >= 0)
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

-- Shop Orders (Sub Order - แยกตามร้าน)
CREATE TABLE IF NOT EXISTS shop_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    subtotal NUMERIC(10, 2) NOT NULL,
    shipping NUMERIC(10, 2) NOT NULL,
    grand_total NUMERIC(10, 2) NOT NULL,
    order_number VARCHAR(20) NOT NULL,
    order_status_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    FOREIGN KEY (order_status_id) REFERENCES order_status(id),
    CONSTRAINT shop_orders_total_price_positive CHECK (grand_total >= 0)
);

CREATE INDEX IF NOT EXISTS idx_shop_orders_order_id ON shop_orders(order_id);
CREATE INDEX IF NOT EXISTS idx_shop_orders_shop_id ON shop_orders(shop_id);
CREATE INDEX IF NOT EXISTS idx_shop_orders_status ON shop_orders(order_status_id);

-- Order Items
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    shop_order_id UUID NOT NULL,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    unit_price NUMERIC(10, 2) NOT NULL,
    subtotal NUMERIC(10, 2) NOT NULL,
    FOREIGN KEY (shop_order_id) REFERENCES shop_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT order_items_quantity_positive CHECK (qty > 0),
    CONSTRAINT order_items_price_positive CHECK (unit_price >= 0)
);

CREATE INDEX IF NOT EXISTS idx_order_items_shop_order_id ON order_items(shop_order_id);

-- Order Logs (Timeline)
CREATE TABLE IF NOT EXISTS order_logs (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    shop_order_id UUID,
    order_status_id INTEGER,
    note TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ(6) DEFAULT NOW(),
    FOREIGN KEY (shop_order_id) REFERENCES shop_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (order_status_id) REFERENCES order_status(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_order_logs_shop_order_id ON order_logs(shop_order_id);

-- ============================================
-- 8. Payment Tables
-- ============================================

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    payment_method_id INTEGER,
    payment_status_id INTEGER NOT NULL DEFAULT 1,
    amount NUMERIC(10, 2) NOT NULL,
    transaction_id TEXT,
    paid_at TIMESTAMPTZ(6),
    expires_at TIMESTAMPTZ(6),
    created_at TIMESTAMPTZ(6) DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id),
    FOREIGN KEY (payment_status_id) REFERENCES payment_status(id),
    CONSTRAINT payments_amount_positive CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(payment_status_id);

-- ============================================
-- 9. Courier & Shipment Tables
-- ============================================

-- Couriers
CREATE TABLE IF NOT EXISTS couriers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_url TEXT,
    rate NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6)
);

-- Shop Couriers (Many-to-Many)
CREATE TABLE IF NOT EXISTS shop_couriers (
    id SERIAL PRIMARY KEY,
    shop_id UUID NOT NULL,
    courier_id INTEGER NOT NULL,
    rate NUMERIC(10, 2),
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ(6),
    FOREIGN KEY (shop_id) REFERENCES shops(id) ON DELETE CASCADE,
    FOREIGN KEY (courier_id) REFERENCES couriers(id) ON DELETE CASCADE,
    CONSTRAINT shop_couriers_fee_positive CHECK (rate >= 0)
);

-- Shipments
CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_order_id UUID NOT NULL,
    courier_id INTEGER NOT NULL,
    shipment_status_id INTEGER NOT NULL DEFAULT 1,
    tracking_no VARCHAR(100) NOT NULL,
    shipped_at TIMESTAMPTZ(6),
    delivered_at TIMESTAMPTZ(6),
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL,
    FOREIGN KEY (shop_order_id) REFERENCES shop_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (courier_id) REFERENCES couriers(id),
    FOREIGN KEY (shipment_status_id) REFERENCES shipment_status(id)
);

CREATE INDEX IF NOT EXISTS idx_shipments_shop_order_id ON shipments(shop_order_id);
CREATE INDEX IF NOT EXISTS idx_shipments_tracking_number ON shipments(tracking_no);

-- ============================================
-- 10. Refund Tables
-- ============================================

-- Refunds
CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_order_id UUID NOT NULL,
    payment_id UUID,
    amount NUMERIC(10, 2) NOT NULL,
    refund_method_id INTEGER,
    refund_status_id INTEGER NOT NULL DEFAULT 1,
    reason TEXT,
    bank_account VARCHAR(100),
    bank_name VARCHAR(100),
    transaction_id TEXT,
    refunded_at TIMESTAMPTZ(6),
    created_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    FOREIGN KEY (shop_order_id) REFERENCES shop_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (refund_status_id) REFERENCES refund_status(id),
    FOREIGN KEY (refund_method_id) REFERENCES refund_methods(id),
    CONSTRAINT refunds_amount_positive CHECK (amount >= 0)
);

CREATE INDEX IF NOT EXISTS idx_refunds_shop_order_id ON refunds(shop_order_id);
CREATE INDEX IF NOT EXISTS idx_refunds_status ON refunds(refund_status_id);

COMMIT;
