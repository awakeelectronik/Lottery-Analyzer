package middleware

import (
	"net/http"
)

// CORS middleware que maneja Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		origin := r.Header.Get("Origin")

		// Permitir orígenes específicos o todos en desarrollo
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Métodos permitidos
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		// Headers permitidos
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		// Headers expuestos al cliente
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		// Permitir credenciales
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Tiempo de cache para preflight requests
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Manejar preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r)
	})
}

// CORSForProduction - versión más restrictiva para producción
func CORSForProduction(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Verificar si el origen está permitido
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Logging middleware para debug
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log básico de requests
		println("Request:", r.Method, r.URL.Path, "Origin:", r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
	})
}

// Recovery middleware para manejar panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				println("Panic recovered:", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
