package main

import (
	"encoding/json"
	"errors"
	"net/http"

	employeside "github.com/Employes-Side/employee-side"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{

	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

		aud := "https://employee-side"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

		if !checkAud {
			return token, errors.New("invalid audience")
		}

		iss := "https://dev-37ekikkl0pu50lfc.us.auth0.com/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("invalid issuer")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

func jwtMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := jwtMiddleware.CheckJWT(w, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getPemCert(token *jwt.Token) (string, error) {

	cert := ""
	resp, err := http.Get("https://dev-37ekikkl0pu50lfc.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}

	defer resp.Body.Close()

	var jwks = employeside.Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"

		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
