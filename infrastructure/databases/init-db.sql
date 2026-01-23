-- ========================================
-- Initialize Databases for All Microservices
-- ========================================
-- This script creates separate databases for each microservice
-- following the database-per-service pattern

\echo 'Creating databases for microservices...'

-- ========================================
-- Create Databases
-- ========================================

-- Auth Service Database
SELECT 'CREATE DATABASE auth_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth_db')\gexec

-- Shop Service Database
SELECT 'CREATE DATABASE shop_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shop_db')\gexec

-- Customer Service Database
SELECT 'CREATE DATABASE customer_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'customer_db')\gexec

-- Inventory Service Database
SELECT 'CREATE DATABASE inventory_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'inventory_db')\gexec

-- Order Service Database
SELECT 'CREATE DATABASE order_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'order_db')\gexec

-- Payment Service Database
SELECT 'CREATE DATABASE payment_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'payment_db')\gexec

-- Shipping Service Database
SELECT 'CREATE DATABASE shipping_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shipping_db')\gexec

-- Notification Service Database
SELECT 'CREATE DATABASE notification_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'notification_db')\gexec

-- Review Service Database
SELECT 'CREATE DATABASE review_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'review_db')\gexec

-- Chat Service Database
SELECT 'CREATE DATABASE chat_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'chat_db')\gexec

-- Search Service Database
SELECT 'CREATE DATABASE search_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'search_db')\gexec

-- Admin Service Database
SELECT 'CREATE DATABASE admin_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'admin_db')\gexec

\echo 'All databases created successfully!'

-- ========================================
-- Grant Privileges
-- ========================================

\c auth_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c shop_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c customer_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c inventory_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c order_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c payment_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c shipping_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c notification_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c review_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c chat_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c search_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
\c admin_db
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;

\c postgres
\echo 'All privileges granted!'

-- ========================================
-- Database Summary
-- ========================================

\echo ''
\echo '========================================='
\echo 'Database Initialization Complete'
\echo '========================================='
\echo ''
\echo 'Created databases:'
\echo '  1. auth_db         - Authentication & Authorization'
\echo '  2. shop_db         - Shop Management'
\echo '  3. customer_db     - Customer Management'
\echo '  4. inventory_db    - Inventory Management'
\echo '  5. order_db        - Order Management'
\echo '  6. payment_db      - Payment Processing'
\echo '  7. shipping_db     - Shipping Management'
\echo '  8. notification_db - Notification Service'
\echo '  9. review_db       - Review Management'
\echo ' 10. chat_db         - Chat Service'
\echo ' 11. search_db       - Search Service'
\echo ' 12. admin_db        - Admin Service'
\echo ''
\echo 'All databases are ready for use!'
\echo '========================================='
