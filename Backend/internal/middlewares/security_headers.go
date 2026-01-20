package middlewares

import "net/http"

// SecurityHeaders adds essential HTTP security headers.
// This middleware is SAFE for production and VAPT compliant.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// XSS protection (legacy but required for audits)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer privacy
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Disable powerful browser features
		w.Header().Set(
			"Permissions-Policy",
			"geolocation=(), microphone=(), camera=(), payment=()",
		)

		// Content Security Policy (safe default)
		// Adjust only if you load external assets
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'",
		)

		// Enforce HTTPS (ONLY enable when HTTPS is active)
		if r.TLS != nil {
			w.Header().Set(
				"Strict-Transport-Security",
				"max-age=63072000; includeSubDomains; preload",
			)
		}

		next.ServeHTTP(w, r)
	})
}
