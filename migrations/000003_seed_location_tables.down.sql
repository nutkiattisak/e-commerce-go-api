-- ===================================
-- Migration Rollback: Seed Location Tables
-- Version: 000003
-- Description: Remove location data
-- ===================================

BEGIN;

DELETE FROM sub_districts;
DELETE FROM districts;
DELETE FROM provinces;

COMMIT;
