CREATE TABLE beli_header (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(100) UNIQUE NOT NULL,
    supplier VARCHAR(200) NOT NULL,
    total DECIMAL(15,2) DEFAULT 0,
    user_id INTEGER REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'selesai',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
