-- Initialize databases for all microservices

-- Auth Service Database
CREATE DATABASE IF NOT EXISTS auth_db;

-- Shop Service Database
CREATE DATABASE IF NOT EXISTS shop_db;

-- Customer Service Database
CREATE DATABASE IF NOT EXISTS customer_db;

-- Inventory Service Database
CREATE DATABASE IF NOT EXISTS inventory_db;

-- Order Service Database
CREATE DATABASE IF NOT EXISTS order_db;

-- Payment Service Database
CREATE DATABASE IF NOT EXISTS payment_db;

-- Shipping Service Database
CREATE DATABASE IF NOT EXISTS shipping_db;

-- Notification Service Database
CREATE DATABASE IF NOT EXISTS notification_db;

-- Review Service Database
CREATE DATABASE IF NOT EXISTS review_db;

-- Chat Service Database
CREATE DATABASE IF NOT EXISTS chat_db;

-- Search Service Database
CREATE DATABASE IF NOT EXISTS search_db;

-- Admin Service Database
CREATE DATABASE IF NOT EXISTS admin_db;

\echo 'All databases created successfully'
