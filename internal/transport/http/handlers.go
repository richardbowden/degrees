package thttp

type Handlers struct {
	Users *UserHandler
	Auth  *AuthN
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
