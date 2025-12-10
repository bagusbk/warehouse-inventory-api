CREATE TABLE jual_detail (
    id SERIAL PRIMARY KEY,
    jual_header_id INTEGER REFERENCES jual_header(id),
    barang_id INTEGER REFERENCES master_barang(id),
    qty INTEGER NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL
);
