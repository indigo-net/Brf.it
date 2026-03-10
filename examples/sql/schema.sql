-- E-commerce database schema

CREATE TABLE products (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    total DECIMAL(10, 2),
    status VARCHAR(50) DEFAULT 'pending',
    ordered_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id BIGINT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1
);

-- Get total revenue by category
CREATE FUNCTION revenue_by_category(cat VARCHAR)
RETURNS DECIMAL AS $$
SELECT COALESCE(SUM(p.price * oi.quantity), 0)
FROM products p
JOIN order_items oi ON oi.product_id = p.id
WHERE p.category = cat;
$$ LANGUAGE sql;

CREATE VIEW top_products AS
SELECT p.name, COUNT(*) as order_count
FROM products p
JOIN order_items oi ON oi.product_id = p.id
GROUP BY p.name
ORDER BY order_count DESC;

CREATE INDEX idx_products_category ON products (category);

CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered');
