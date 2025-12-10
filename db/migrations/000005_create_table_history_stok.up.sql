CREATE TABLE history_stok (
    id SERIAL PRIMARY KEY,
    barang_id INTEGER REFERENCES master_barang(id),
    user_id INTEGER REFERENCES users(id),
    jenis_transaksi VARCHAR(50) NOT NULL, -- 'masuk', 'keluar', 'adjustment'
    jumlah INTEGER NOT NULL,
    stok_sebelum INTEGER NOT NULL,
    stok_sesudah INTEGER NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
