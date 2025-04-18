CREATE TABLE IF NOT EXISTS payments
(
    id           varchar(36) primary key,
    order_id     varchar(36)    not null,
    customer_id  varchar(36)    not null,
    amount       decimal(10, 2) not null,
    currency     varchar(3)     not null,
    status       varchar(20)    not null,
    processed_at timestamp      not null,
    error_msg    text,
    created_at   timestamp      not null,
    updated_at   timestamp      not null
);

create index if not exists idx_payments_order_id on payments (order_id);
create index if not exists idx_payments_customer_id on payments (customer_id);
create index if not exists idx_payments_status on payments (status);
