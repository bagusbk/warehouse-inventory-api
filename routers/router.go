package routers

import (
	"warehouse/handlers"
	"warehouse/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Warehouse API is running",
		})
	})

	userHandler := handlers.NewUserHandler()
	barangHandler := handlers.NewBarangHandler()
	stokHandler := handlers.NewStokHandler()
	pembelianHandler := handlers.NewPembelianHandler()
	penjualanHandler := handlers.NewPenjualanHandler()
	laporanHandler := handlers.NewLaporanHandler()

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", userHandler.GetProfile)

			protected.GET("/barang", barangHandler.GetAll)
			protected.GET("/barang/stok", barangHandler.GetAllWithStok)
			protected.GET("/barang/:id", barangHandler.GetByID)

			protected.GET("/stok", stokHandler.GetAll)
			protected.GET("/stok/:barang_id", stokHandler.GetByBarang)
			protected.GET("/history-stok", stokHandler.GetHistory)
			protected.GET("/history-stok/:barang_id", stokHandler.GetHistoryByBarang)

			protected.POST("/pembelian", pembelianHandler.Create)
			protected.GET("/pembelian", pembelianHandler.GetAll)
			protected.GET("/pembelian/:id", pembelianHandler.GetByID)

			protected.POST("/penjualan", penjualanHandler.Create)
			protected.GET("/penjualan", penjualanHandler.GetAll)
			protected.GET("/penjualan/:id", penjualanHandler.GetByID)

			protected.GET("/laporan/stok", laporanHandler.LaporanStok)
			protected.GET("/laporan/pembelian", laporanHandler.LaporanPembelian)
			protected.GET("/laporan/penjualan", laporanHandler.LaporanPenjualan)

			admin := protected.Group("")
			admin.Use(middleware.AdminOnly())
			{
				admin.POST("/barang", barangHandler.Create)
				admin.PUT("/barang/:id", barangHandler.Update)
				admin.DELETE("/barang/:id", barangHandler.Delete)
			}
		}
	}

	return r
}
