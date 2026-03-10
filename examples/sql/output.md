# Code Summary: examples/sql

*brf.it dev*

---

## Files

### examples/sql/schema.sql

```sql
CREATE TABLE products (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
)
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    total DECIMAL(10, 2),
    status VARCHAR(50) DEFAULT 'pending',
    ordered_at TIMESTAMP DEFAULT NOW()
)
CREATE FUNCTION revenue_by_category(cat VARCHAR)
RETURNS DECIMAL LANGUAGE sql
CREATE VIEW top_products
CREATE INDEX idx_products_category ON products (category)
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered')
```

---

