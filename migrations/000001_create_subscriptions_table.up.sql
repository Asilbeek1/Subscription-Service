CREATE TABLE subscriptions(
    id SERIAL PRIMARY KEY
    service_name TEXT NOT NULL,
    pirce INTEGER NOT NULL CHECK(price > 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);