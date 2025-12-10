INSERT INTO history_stok (barang_id, user_id, jenis_transaksi, jumlah, stok_sebelum, stok_sesudah, keterangan) VALUES
(1, 2, 'masuk', 2, 0, 2, 'Pembelian BLI001'),
(2, 2, 'masuk', 10, 0, 10, 'Pembelian BLI001'),
(3, 3, 'masuk', 5, 0, 5, 'Pembelian BLI002'),
(4, 3, 'masuk', 3, 0, 3, 'Pembelian BLI002'),
(5, 3, 'masuk', 4, 0, 4, 'Pembelian BLI002'),
(1, 2, 'keluar', 1, 2, 1, 'Penjualan JUAL001'),
(2, 2, 'keluar', 2, 10, 8, 'Penjualan JUAL001'),
(3, 2, 'keluar', 1, 5, 4, 'Penjualan JUAL001'),
(2, 3, 'keluar', 5, 8, 3, 'Penjualan JUAL002'),
(4, 3, 'keluar', 1, 3, 2, 'Penjualan JUAL002');