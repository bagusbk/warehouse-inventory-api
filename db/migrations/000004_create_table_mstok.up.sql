CREATE TABLE mstok (
    id SERIAL PRIMARY KEY,
    barang_id INTEGER REFERENCES master_barang(id),
    stok_akhir INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
