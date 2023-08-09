CREATE SCHEMA processout;

CREATE TABLE processout.merchant (
    id SERIAL primary key,
    name text,
);

CREATE TABLE processout.payments (
    id int primary key,
    card_number text,
    card_type text,
    card_name text,
    cvv text, 
    card_expiry_month int,
    card_expirt_year int,
    amount int, 
    details text,
);
