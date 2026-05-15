-- 02.meals.sql
-- Tables for meals and meal pictures

CREATE TABLE meals (
    id BIGINT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL,  -- in kobo
    enabled BOOL DEFAULT FALSE,
    category TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    score DOUBLE PRECISION NOT NULL DEFAULT 50,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_meals_vendor_id ON meals(vendor_id);
CREATE INDEX idx_meals_category ON meals(category);
CREATE INDEX idx_meals_status ON meals(status);

CREATE TABLE meal_pictures (
    id BIGINT PRIMARY KEY,
    meal_id BIGINT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    public_id TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_meal_pictures_meal_id ON meal_pictures(meal_id);

CREATE TABLE user_purchased_meals (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, meal_id)
);

CREATE INDEX idx_user_purchased_meals_user_id ON user_purchased_meals(user_id);
CREATE INDEX idx_user_purchased_meals_meal_id ON user_purchased_meals(meal_id);


CREATE TABLE cart_items (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_id BIGINT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    quantity SMALLINT NOT NULL
);

CREATE INDEX idx_cart_item_customer_id ON cart_items(customer_id);
CREATE INDEX idx_cart_item_customer_id_meal_id ON cart_items(customer_id, meal_id);