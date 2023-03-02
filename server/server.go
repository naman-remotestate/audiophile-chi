package server

import (
	"audiophile/handlers"
	"audiophile/middlewares"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	ReadTimeOut       = 5 * time.Minute
	ReadHeaderTimeOut = 30 * time.Second
	WriteTimeout      = 5 * time.Minute
)

func SetUpRoutes() *Server {
	r := chi.NewRouter()
	r.Route("/public", func(public chi.Router) {
		public.Post("/register", handlers.RegisterUser)
		public.Post("/login", handlers.Login)
		public.Put("/logout", handlers.Logout)
	})
	r.Route("/admin", func(admin chi.Router) {
		admin.Use(middlewares.Auth)
		admin.Use(middlewares.IsAdmin)
		admin.Get("/users", handlers.GetAllUsers)
		admin.Get("/users/{userId}/all-address", handlers.GetAllAddressOfUser)
		admin.Post("/inventory", handlers.AddToInventory)
		admin.Put("/inventory/{productId}", handlers.UpdateProductInfo)
		admin.Delete("/inventory/{productId}", handlers.DeleteProductFromInventory)
		admin.Post("/product/image/{productId}", handlers.UploadProductImage)
		admin.Get("/product/image/{productId}", handlers.GetProductImages)
		admin.Put("/product/image/{productImage}", handlers.UpdateImage)
		admin.Get("/products", handlers.GetAllProducts)
	})
	r.Route("/user", func(user chi.Router) {
		user.Use(middlewares.Auth)
		user.Get("/cart/all", handlers.GetCartDetails)
		user.Post("/cart/{productId}", handlers.AddProductToCart)
		user.Delete("/cart/{productId}", handlers.RemoveProductFromCart)
		user.Get("/address", handlers.GetAllAddress)
		user.Post("/address", handlers.AddAddress)
		user.Put("/address/{addressId}", handlers.UpdateAddress)
		user.Delete("/address/{addressId}", handlers.RemoveAddress)
		user.Get("/orders", handlers.AllOrders)
		user.Post("/order", handlers.PlaceOrder)
		user.Put("/order/{orderId}", handlers.CancelOrder)
		user.Get("/", handlers.GetAllProducts)
	})
	r.Route("/products", func(common chi.Router) {
		common.Use(middlewares.Auth)
		common.Get("/", handlers.GetAllProducts)
		common.Get("/{productId}", handlers.GetProduct)
	})

	return &Server{
		Router: r,
	}
}

func (svc *Server) Start(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       ReadTimeOut,
		ReadHeaderTimeout: ReadHeaderTimeOut,
		WriteTimeout:      WriteTimeout,
	}
	return svc.server.ListenAndServe()
}
