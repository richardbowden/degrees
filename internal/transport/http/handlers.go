package thttp

import (
	fastmail "github.com/typewriterco/p402/internal/email/genericsmtp"
)

// Handlers contains HTTP handlers for functionality not provided by gRPC-Gateway
// Most API endpoints are handled by gRPC-Gateway, these are for special cases
type Handlers struct {
	SMTP *fastmail.Client // SMTP admin - not in proto, requires direct HTTP access
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// func (h *Handlers) RegisterPublicRoutes(r chi.Router) {
// 	if h.Users != nil {
// 		r.Route("/users", func(r chi.Router) {
// 			h.Users.RegisterPublicRoutes(r)
// 		})
// 	}
// }
//
// func (h *Handlers) RegisterProtectedRoutes(r chi.Router) {
// 	if h.Users != nil {
// 		r.Route("/users", func(r chi.Router) {
// 			h.Users.RegisterProtectedRoutes(r)
// 		})
// 	}
// }
