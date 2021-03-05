package traefik_jwt_plugin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	traefik_jwt_plugin "github.com/team-carepay/traefik-jwt-plugin"
)

func TestServeHTTPOK(t *testing.T) {
	cfg := traefik_jwt_plugin.CreateConfig()
	cfg.PayloadFields = []string { "exp" }
	cfg.Keys = []string { "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv\nvkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc\naT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy\ntvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0\ne+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb\nV6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9\nMwIDAQAB\n-----END PUBLIC KEY-----" }
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })

	jwt, err := traefik_jwt_plugin.New(ctx, next, cfg, "test-traefik-jwt-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header["Authorization"] = []string{"Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.JlX3gXGyClTBFciHhknWrjo7SKqyJ5iBO0n-3S2_I7cIgfaZAeRDJ3SQEbaPxVC7X8aqGCOM-pQOjZPKUJN8DMFrlHTOdqMs0TwQ2PRBmVAxXTSOZOoEhD4ZNCHohYoyfoDhJDP4Qye_FCqu6POJzg0Jcun4d3KW04QTiGxv2PkYqmB7nHxYuJdnqE3704hIS56pc_8q6AW0WIT0W-nIvwzaSbtBU9RgaC7ZpBD2LiNE265UBIFraMDF8IAFw9itZSUCTKg1Q-q27NwwBZNGYStMdIBDor2Bsq5ge51EkWajzZ7ALisVp-bskzUsqUf77ejqX_CBAqkNdH1Zebn93A"}

	jwt.ServeHTTP(recorder, req)

	if nextCalled == false {
		t.Fatal("next.ServeHTTP was not called")
	}
}

func TestServeHTTPInvalidSignature(t *testing.T) {
	cfg := traefik_jwt_plugin.CreateConfig()
	cfg.PayloadFields = []string { "exp" }
	cfg.Keys = []string { "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv\nvkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc\naT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy\ntvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0\ne+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb\nV6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9\nMwIDAQAB\n-----END PUBLIC KEY-----" }
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })

	jwt, err := traefik_jwt_plugin.New(ctx, next, cfg, "test-traefik-jwt-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header["Authorization"] = []string{"Bearer A.B.C"}

	jwt.ServeHTTP(recorder, req)

	if nextCalled == true {
		t.Fatal("next.ServeHTTP was called")
	}
}

func TestServeHTTPMissingExp(t *testing.T) {
	cfg := traefik_jwt_plugin.CreateConfig()
	cfg.PayloadFields = []string { "exp" }
	cfg.Required = true
	cfg.Keys = []string { "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv\nvkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc\naT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy\ntvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0\ne+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb\nV6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9\nMwIDAQAB\n-----END PUBLIC KEY-----" }
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })

	jwt, err := traefik_jwt_plugin.New(ctx, next, cfg, "test-traefik-jwt-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header["Authorization"] = []string{"Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.JlX3gXGyClTBFciHhknWrjo7SKqyJ5iBO0n-3S2_I7cIgfaZAeRDJ3SQEbaPxVC7X8aqGCOM-pQOjZPKUJN8DMFrlHTOdqMs0TwQ2PRBmVAxXTSOZOoEhD4ZNCHohYoyfoDhJDP4Qye_FCqu6POJzg0Jcun4d3KW04QTiGxv2PkYqmB7nHxYuJdnqE3704hIS56pc_8q6AW0WIT0W-nIvwzaSbtBU9RgaC7ZpBD2LiNE265UBIFraMDF8IAFw9itZSUCTKg1Q-q27NwwBZNGYStMdIBDor2Bsq5ge51EkWajzZ7ALisVp-bskzUsqUf77ejqX_CBAqkNdH1Zebn93A"}

	jwt.ServeHTTP(recorder, req)

	if nextCalled == true {
		t.Fatal("next.ServeHTTP was called")
	}
}

func TestServeHTTPAllowed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/data/testok" {
			t.Fatal(fmt.Sprintf("Path incorrect: %s", r.URL.Path))
		}
		param1 := r.URL.Query()["Param1"]
		if len(param1) != 2 || param1[0] != "foo" || param1[1] != "bar" {
			t.Fatal(fmt.Sprintf("Parameters incorrect, expected foo,bar but got %s", strings.Join(param1, ",")))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "{ \"result\": { \"allow\": true } }")
	}))
	defer ts.Close()
	cfg := traefik_jwt_plugin.CreateConfig()
	cfg.OpaUrl = fmt.Sprintf("%s/v1/data/testok?Param1=foo&Param1=bar", ts.URL)
	cfg.OpaAllowField = "allow"
	ctx := context.Background()
	nextCalled := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { nextCalled = true })

	opa, err := traefik_jwt_plugin.New(ctx, next, cfg, "test-traefik-jwt-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	opa.ServeHTTP(recorder, req)

	if recorder.Code == http.StatusForbidden {
		t.Fatal("Exptected OK")
	}
	if nextCalled == false {
		t.Fatal("next.ServeHTTP was not called")
	}
}

func TestServeHTTPForbidden(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "{ \"result\": { \"allow\": false } }")
	}))
	defer ts.Close()
	cfg := traefik_jwt_plugin.CreateConfig()
	cfg.OpaUrl = ts.URL
	cfg.OpaAllowField = "allow"
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) { t.Fatal("Should not chain HTTP call") })

	opa, err := traefik_jwt_plugin.New(ctx, next, cfg, "test-traefik-jwt-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	opa.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatal("Exptected Forbidden")
	}
}