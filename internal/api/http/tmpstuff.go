package apihttp

// import (
// 	"encoding/json"
// 	"net/http"

// 	"git.moonshot.com/richard/opcounter/internal/errs"
// 	"git.moonshot.com/richard/opcounter/internal/users"
// 	"github.com/alexedwards/scs/v2"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/httplog"
// 	"github.com/google/uuid"
// 	"github.com/markbates/goth/gothic"
// )

// type endpoints struct {
// 	service        users.UserService
// 	sessionManager *scs.SessionManager
// }

// func NewEndpoints(svc users.UserService, sessionManager *scs.SessionManager) *endpoints {
// 	return &endpoints{
// 		service:        svc,
// 		sessionManager: sessionManager,
// 	}
// }

// func (e *endpoints) Endpoints() http.Handler {
// 	r := chi.NewRouter()
// 	// r.Get("/register", e.RegisterUserWithPassword)

// 	r.Group(func(r chi.Router) {
// 		r.Use(LogInMiddleware(e.sessionManager))
// 		r.Get("/login", e.Login)

// 	})

// 	r.Group(func(r chi.Router) {
// 		r.Use(e.sessionManager.LoadAndSave)
// 		r.Get("/logout", e.Logout)
// 	})

// 	r.Get("/oauth/callback", e.OAuthCallback)

// 	r.Get("/checkusername", e.CheckUsernameAvailability)

// 	return r
// }

// func (e *endpoints) CheckUsernameAvailability(w http.ResponseWriter, r *http.Request) {
// 	type request struct {
// 		Username string `json:"username"`
// 	}
// 	var req request
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "Bad request reading body", http.StatusBadRequest)
// 		return
// 	}

// 	response := map[string]bool{"available": true}

// 	json.NewEncoder(w).Encode(response)
// }

// func (e *endpoints) RegisterUserWithPassword(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("hello"))
// }

// func (e *endpoints) Login(w http.ResponseWriter, r *http.Request) {
// 	// try to get the user without re-authenticating
// 	// if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
// 	// 	fmt.Println(gothUser.Email)
// 	// } else {
// 	gothic.BeginAuthHandler(w, r)
// 	// }
// }

// func (e *endpoints) Logout(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("logged out"))
// 	w.WriteHeader(http.StatusOK)
// 	gothic.Logout(w, r)
// }

// func (e *endpoints) OAuthCallback(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	log := httplog.LogEntry(ctx)

// 	user, err := gothic.CompleteUserAuth(w, r)

// 	if err != nil {
// 		e := errs.E(errs.Internal, "failed to complete user auth", err)
// 		errs.HTTPErrorResponse(w, log, e)
// 		return
// 	}

// 	var loggingInUser *users.User
// 	loggingInUser, exists, err := e.service.GetUser(ctx, user.Email)

// 	if err != nil {
// 		e := errs.E(errs.Database, "failed to get user in OAuthCallback", err)
// 		errs.HTTPErrorResponse(w, log, e)
// 		return
// 	}

// 	if !exists {
// 		n := users.NewSignInParams{
// 			FirstName: user.FirstName,
// 			Surname:   user.LastName,
// 			Email:     user.Email,
// 		}

// 		loggingInUser, err = e.service.NewSignIn(ctx, n)

// 		if err != nil {
// 			e := errs.E(errs.Database, "issue trying to new sign in", err)
// 			errs.HTTPErrorResponse(w, log, e)
// 			return
// 		}
// 	}

// 	log.Info().Str("acc_num", loggingInUser.AccNum).Msg("")

// 	ctx, err = e.sessionManager.Load(ctx, "opcounter")

// 	if err != nil {
// 		log.Error().Err(err).Msg("failed to get new session")
// 		return
// 	}

// 	uid := uuid.UUID{}
// 	uid = loggingInUser.ID.Bytes

// 	e.sessionManager.Put(ctx, "user_id", uid.String())
// 	ss, tt, err := e.sessionManager.Commit(ctx)
// 	if err != nil {
// 		log.Error().Err(err).Msg("")
// 		return
// 	}

// 	e.sessionManager.WriteSessionCookie(ctx, w, ss, tt)

// 	dd, ok := e.sessionManager.Get(ctx, "user_id").(string)

// 	if !ok {
// 		log.Error().Msg("failed to get user_id")
// 		return
// 	}

// 	_ = dd

// 	// log.Debug()
// 	_ = ss
// 	_ = tt

// 	// if uuu.Enabled {
// 	// 	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
// 	// 	return
// 	// }

// 	// if err != nil {
// 	// 	log.Error().Err(err).Msg("")
// 	// }

// 	// _ = u
// 	// session, err := e.sessionManager.Load(ctx, "sid")
// 	// _ = session
// 	// if err != nil {
// 	// 	e := errs.E(errs.Internal, "failed to load a new session", err)
// 	// 	errs.HTTPErrorResponse(w, log, e)
// 	// 	return
// 	// }

// 	// b := make([]byte, 33)
// 	// _, err = rand.Read(b)

// 	// token := base64.URLEncoding.EncodeToString(b)

// 	// e.sessionManager.Put(ctx, "sid", token)

// 	// e.sessionManager.Load()

// 	// s, err := sessions.Start(w, r, true)

// 	// if err != nil {
// 	// 	e := errs.E(errs.Internal, "failed to make session", err)
// 	// 	errs.HTTPErrorResponse(w, log, e)
// 	// 	return
// 	// }

// 	// err = s.LogIn(u, true, w)

// 	// if err != nil {
// 	// 	e := errs.E(errs.Internal, "failed to log session in", err)
// 	// 	errs.HTTPErrorResponse(w, log, e)
// 	// 	return
// 	// }

// 	http.Redirect(w, r.WithContext(ctx), "/", http.StatusTemporaryRedirect)
// }
