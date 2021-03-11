package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/google/go-github/github"
	"github.com/hayashiki/go-pkg/slack/auth"
	"github.com/hayashiki/mentions/pkg/config"
	ghSvc "github.com/hayashiki/mentions/pkg/github"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/repository"
	"github.com/hayashiki/mentions/pkg/usecase"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	UserScopes = []string{"identity.basic"}
	//"incoming-webhook"
	TeamScopes = []string{"chat:write", "users.profile:read", "users:read"}
)

type adminHandler func(w http.ResponseWriter, r *http.Request) error

type App struct {
	isDev            bool
	config           config.Config
	userRepo         repository.UserRepository
	teamRepo         repository.TeamRepository
	repoRepo         repository.RepoRepository
	installationRepo repository.InstallationRepository
	taskRepo         repository.TaskRepository
	token            auth.Token
	ghSvc            ghSvc.Github
	ghAppSvc         ghSvc.Github
	appsTransport    *ghinstallation.AppsTransport
}

func NewApp(
	config config.Config,
	ghSvc ghSvc.Github,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	repoRepo repository.RepoRepository,
	taskRepo repository.TaskRepository,
	installationRepo repository.InstallationRepository,
	ghAppSvc ghSvc.Github,
) Server {
	return &App{
		isDev:            true,
		config:           config,
		ghSvc:            ghSvc,
		userRepo:         userRepo,
		teamRepo:         teamRepo,
		taskRepo:         taskRepo,
		repoRepo:         repoRepo,
		installationRepo: installationRepo,
		ghAppSvc:         ghAppSvc,
	}
}

type httpError struct {
	err     error
	message string
	code    int
}

func newHttpError(err error, message string, code int) (herr *httpError) {
	herr = new(httpError)
	herr.err = err
	herr.message = message
	herr.code = code

	return herr
}

func (e *httpError) Error() string {
	return e.err.Error()
}

func (h adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		log.WithError(err).Error("debug err")
		if err, ok := err.(*httpError); ok {
			if len(err.message) > 0 {
				log.WithError(err).Errorf("%s: %s", err.message, err.Error())
			} else {
				log.WithError(err)
			}
			http.Error(w, err.Error(), err.code)
		} else {
			log.WithError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type (
	// Server interface
	Server interface {
		Handler() chi.Router
	}
)

func (a *App) Handler() chi.Router {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2001-07-06 12:03:15"
	formatter.FullTimestamp = true
	log.SetLevel(log.DebugLevel)

	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
	)
	r.Method(http.MethodPost, "/webhook/github", adminHandler(a.postWebhook))
	r.Method(http.MethodPost, "/webhook/github/apps", adminHandler(a.postAppsWebhook))
	//r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
	//	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	//})

	r.Method(http.MethodGet, "/slack/team/auth", adminHandler(a.slackTeamAuthorize))
	r.Method(http.MethodGet, "/slack/user/auth", adminHandler(a.slackUserAuthorize))
	r.Method(http.MethodGet, "/slack/team/callback", adminHandler(a.slackTeamCallback))
	r.Method(http.MethodGet, "/slack/user/callback", adminHandler(a.slackUserCallback))

	r.Route("/api", func(r chi.Router) {
		r.Use(AuthorizationMiddleware)
		r.Method(http.MethodPost, "/users", adminHandler(a.CreateUser))
		r.Method(http.MethodGet, "/users", adminHandler(a.ListUser))
		r.Method(http.MethodPatch, "/users/{id}", adminHandler(a.UpdateUser))
		r.Method(http.MethodGet, "/slack/users", adminHandler(a.ListSlackAPIUser))
	})
	//r.Method(http.MethodGet, "/signup", adminHandler(a.loginHandler))
	//r.Method(http.MethodGet, "/", adminHandler(a.appHandler))

	r.Get("/github/callback", GithubCallback)
	r.Method(http.MethodGet, "/github/installation/callback", adminHandler(a.InstallationCallback))

	r.NotFound(notFoundHandler)

	return r
}

type InstallationResponse struct {
	installationID string
	setupAction    string
}

func (a *App) InstallationCallback(w http.ResponseWriter, r *http.Request) error {
	log.Debug("Installation Callback")
	resp := &InstallationResponse{
		installationID: r.FormValue("installation_id"),
		setupAction:    r.FormValue("setup_action"),
	}

	int, err := strconv.Atoi(resp.installationID)

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, a.config.GithubAppID, int64(int), a.config.GithubAppPrivateKeyFileName)
	if err != nil {
		log.WithError(err).Error("failed to read keyFile")
		return err
	}
	client := github.NewClient(&http.Client{Transport: itr})
	inst := &model.Installation{
		ID: int64(int),
		//https://api.github.com/app/installations/xxxx: 401 A JSON web token could not be decoded []
		// なぜか取得できないのでName取得は諦めた
		//Name:      ins.Account.GetName(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = a.installationRepo.Put(r.Context(), inst)
	if err != nil {
		log.WithError(err).Errorf("failed to put installation %v", inst)
		return err
	}

	// ページング処理はしていない、1000件も取得していれば十分という判断
	repos, _, err := client.Apps.ListRepos(context.Background(), &github.ListOptions{
		PerPage: 1000,
	})
	for _, repo := range repos {
		parts := strings.Split(repo.GetFullName(), "/")
		dsRepo := &model.Repo{
			ID:       repo.GetID(),
			Owner:    parts[0],
			Name:     parts[1],
			FullName: repo.GetFullName(),
		}
		err := a.repoRepo.Put(r.Context(), dsRepo)
		if err != nil {
			log.WithError(err).Errorf("failed to put repo %v", repo)
			return err
		}
	}
	// TODO: 正常登録されました的なHTMLを返す
	w.Write([]byte("OK"))

	return nil
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	bufBody := new(bytes.Buffer)
	bufBody.ReadFrom(r.Body)
	body := bufBody.String()
	log.Debugf("Called Github App Callback %v", body)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) error {
	// TODO: validate
	req := struct {
		SlackID  string `json:"slackId"`
		GithubID string `json:"githubId"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	u := a.currentUser(r.Context())

	log.Printf("req.SlackID %v", req.SlackID)

	userCreate := usecase.NewUserCreate(a.teamRepo, a.userRepo)
	input := usecase.UserCreateInput{
		SlackTeamID: u.TeamID,
		UserID:      req.SlackID,
		GithubID:    req.GithubID,
		SlackID:     req.SlackID,
	}

	if user, err := userCreate.Do(r.Context(), input); err != nil {
		return err
	} else {
		jsonResponse(w, http.StatusCreated, user)
		//w.Write([]byte("ok"))
		return nil
	}
}

// jsonResponse(w, http.StatusOK, items)的な
func jsonResponse(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": http.StatusText(http.StatusNotFound)})
}

func (a *App) postWebhook(w http.ResponseWriter, r *http.Request) error {

	webhookProcess := usecase.NewWebhookProcess(a.config, a.ghSvc, a.ghAppSvc, a.userRepo, a.taskRepo, a.repoRepo)
	if err := webhookProcess.Do(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (a *App) postAppsWebhook(w http.ResponseWriter, r *http.Request) error {

	webhookProcess := usecase.NewWebhookProcess(a.config, a.ghSvc, a.ghAppSvc, a.userRepo, a.taskRepo, a.repoRepo)
	if err := webhookProcess.Do(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

//
//type user struct {
//	ID   int64
//}
//
//// TODO errをかえり値に
//func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
//	session, err := app.sessions.New(r, defaultSessionID)
//	if err != nil {
//		return
//	}
//	session.Options.MaxAge = -1
//	if err := session.Save(r, w); err != nil {
//		return
//	}
//	return
//}
//
//func (app *App) authMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if err := func() error {
//			session, err := app.sessions.Get(r, defaultSessionID)
//			if err != nil {
//				return err
//			}
//			if value, exists := session.Values[userSessionKey]; exists {
//				if u, ok := value.(*user); ok {
//					r = r.WithContext(context.WithValue(r.Context(), contextKeyUser, u))
//					return nil
//				}
//			}
//			authHeader := r.Header.Get("Authorization")
//			bearerPrefix := "Bearer "
//			if strings.HasPrefix(authHeader, bearerPrefix) {
//				token := strings.TrimPrefix(authHeader, bearerPrefix)
//				log.Printf("token is %s", token)
//			}
//			return nil
//		}(); err != nil {
//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			return
//		}
//	})
//}
//
//var cookieNameToken = "TOKEN"
//
//func (app *App) authCallback(w http.ResponseWriter, r *http.Request) {
//	body := struct {
//		Code string
//		State string
//	}{}
//	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
//		return
//		//return err
//	}
//	defer r.Body.Close()
//
//	// saveUser
//	//u, err := app.entity.SaveUser(r.Context(), userInfo.ID, userInfo.Login)
//	defaultSession, err := app.sessions.New(r, defaultSessionID )
//	if err != nil {
//		return
//		//return err
//	}
//	defaultSession.Values[userSessionKey] = &user{
//		ID:   userInfo.ID,
//		Role: u.Role,
//	}
//
//	// ユーザがいなければ作成する
//	user :=  model.User{ID: resp.ID, TeamID: resp.Team}
//	user.saveしたい
//
//	//key := datastore.NewKey(ctx, "User", userinfo.Sub, 0, nil)
//	//u := &User{ID: userinfo.Sub}
//	//err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
//	//	err := datastore.Get(ctx, key, u)
//	//	if err != nil && err != datastore.ErrNoSuchEntity {
//	//		log.Debugf(ctx, "user exists: %v", u)
//	//		return err
//	//	}
//	//	_, err = datastore.Put(ctx, key, u)
//	//	return err
//	//}, nil)
//	//if err != nil {
//	//	log.Errorf(ctx, "Transaction failed: %v", err)
//	//	return err
//	//}
//	//log.Debugf(ctx, "user created: %v", u)
//
//	//jwt生成して、セッションにいれたい
//
//	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512", jwt.MapClaims{
//		"sub": resp.ID,
//		"exp": time.Now().Add(time.Hour * 24).Unix(),
//	}))
//
//	hmackey, err := GetHMACKey()
//	if err != nil {
//		return err
//	}
//
//	signedToken, err := token.SignedString(hmackey.Bytes())
//	if err != nil {
//		return err
//	}
//	http.SetCookie(w, &http.Cookie{
//		Name:  cookieNameToken,
//		Value: signedToken,
//		Path:  "/",
//	})
//
//	if err := defaultSession.Save(r, w); err != nil {
//		return &appError{err, "failed to save session"}
//	}
//
//	return nil
//}
//
//func GetHMACKey() (uuid.UUID, error) {
//	key := os.Getenv("HMAC_KEY")
//	if key == "" {
//		return uuid.Must(uuid.NewV4()), nil
//	}
//	return uuid.FromString(key)
//}
//
//
