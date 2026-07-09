package router

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/middleware"
	"server/internal/proxy"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	mux *chi.Mux
	cfg *config.Config
}

func New(cfg *config.Config, rdb *redis.Client) *Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestId) // Fixed name
	r.Use(middleware.Tracing)
	r.Use(middleware.PrometheusMetrics) // Should be early
	r.Use(middleware.Logging)
	r.Use(middleware.CORS())
	r.Use(middleware.Compress())
	r.Use(middleware.RateLimitRedis(rdb, 50))

	return &Router{
		mux: r,
		cfg: cfg,
	}

}

func (r *Router) Setup() {
	r.mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("OK"))
	})

	r.mux.Handle("/metrics", promhttp.Handler())

	protected := chi.NewRouter()
	// protected.Use(middleware.JWTAuth(r.cfg.JWTSECRET))
	for _, route := range r.cfg.Routes {

		service, exists := r.cfg.Services[route.Service]
		if !exists {
			log.Printf("Warning: Service '%s' not found for route %s", route.Service, route.Path)
			continue
		}

		if len(service.Instances) == 0 {
			log.Printf("Warning: No instances found for service '%s'", route.Service)
			continue
		}

		serviceConfig := proxy.Service{
			ServiceName: route.Service,
			Instances:   service.Instances,
			RetryCount:  3,
			Timeout:     10 * time.Second,
		}
		target := route.Service

		// create the proxy handler
		proxyHandler, err := proxy.NewHandler(serviceConfig)
		if err != nil {
			log.Printf("Failed to create proxy for %s -> %s: %v", route.Path, target, err)
			continue
		}

		// register the routes
		protected.Handle(route.Path, proxyHandler)
		protected.Handle(route.Path+"/*", proxyHandler)

		log.Printf(" Route registered: %s → %s (%s)", route.Path, target, route.Service)

	}

	r.mux.Mount("/", protected)
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
