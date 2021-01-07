package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

//type key int
//
//var (
//	UserIDKey key
////	UserIDKey
//)

//// ContextKey represents a context key
type ContextKey string

const (
	// UserIDKey is the key for the user id of the authenticated user
	UserIDKey ContextKey = "userID"
	ReqIDKey  ContextKey = "reqID"
)

func GetTokenFromRequest(r *http.Request, ctx context.Context) (*jwt.Token, error) {
	h := r.Header.Get("Authorization")

	if h == "" {
		return nil, fmt.Errorf("Auth header empty")
	}

	parts := strings.SplitN(h, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, fmt.Errorf("Invalid auth header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		//id, err := GetHMACKey()
		//if err != nil {
		//	return nil, err
		//}
		return []byte("hoge"), nil
	})
}

func GetTokenFromCookie(r *http.Request, ctx context.Context) (*jwt.Token, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		//id, err := GetHMACKey()
		//if err != nil {
		//	return nil, err
		//}
		return []byte("hoge"), nil
	})
}

type User struct {
	ID     string
	TeamID string
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := GetTokenFromCookie(r, ctx)

		if err != nil {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		log.Printf("token %v", token)

		if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			log.Printf("Get Claims: %v", err)
			http.Error(w, http.StatusText(403), 403)
			return
		} else if id, ok := claims["custom.id"].(string); !ok {
			log.Printf("Get Sub: %v", err)
			http.Error(w, http.StatusText(403), 403)
			return
		} else if team, ok := claims["custom.namespace"].(string); !ok {
			log.Printf("Get Sub: %v", err)
			http.Error(w, http.StatusText(403), 403)
			return
		} else {
			log.Printf("claim id is %v", id)
			log.Printf("claim namespace is %v", team)
			//key := datastore.NewKey(ctx, "User", sub, 0, nil)
			//u := new(User)
			//if err := datastore.Get(ctx, key, u); err != nil {
			//	return err
			//}
			//e.Set("User", u)
			u := &User{
				ID:     id,
				TeamID: team,
			}
			ctx = context.WithValue(r.Context(), UserIDKey, u)
			log.Printf("ctx %v", ctx)

		}
		log.Printf("ctx %v", ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
