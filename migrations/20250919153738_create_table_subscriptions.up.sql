CREATE TABLE IF NOT EXISTS effective_mobile.subscriptions (
    id UUID PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
