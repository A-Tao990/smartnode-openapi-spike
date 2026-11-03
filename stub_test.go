package spike

// STUB: re-implements only the observable HTTP behaviour of snroute (Read route +
// requireToken + response.WriteResponse envelope) with hard-coded fake data.
// It is NOT the real snroute and has no chain access.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/A-Tao990/smartnode-openapi-spike/api"
)

type stub struct{ fail bool }

func (s stub) GetNetworkNodeFee(w http.ResponseWriter, r *http.Request) {
	resp := api.NodeFeeResponse{Status: "success", NodeFee: 0.14, MinNodeFee: 0.05, TargetNodeFee: 0.1, MaxNodeFee: 0.2}
	code := http.StatusOK
	if s.fail {
		resp = api.NodeFeeResponse{Status: "error", Error: "boom"}
		code = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}

// requireToken mirrors snroute: header "Authorization: Bearer <tok>",
// 401 + "WWW-Authenticate: Bearer" + APIResponse-shaped body.
func requireToken(tok string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") || strings.TrimSpace(h[len("Bearer "):]) != tok {
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"status":"error","error":"unauthorized"}` + "\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearer(tok string) api.RequestEditorFn {
	return func(_ context.Context, r *http.Request) error {
		r.Header.Set("Authorization", "Bearer "+tok)
		return nil
	}
}

func newClient(t *testing.T, url string) *api.ClientWithResponses {
	t.Helper()
	c, err := api.NewClientWithResponses(url)
	if err != nil {
		t.Fatalf("NewClientWithResponses: %v", err)
	}
	return c
}

func TestContract(t *testing.T) {
	const token = "t0k"
	srv := httptest.NewServer(requireToken(token, api.Handler(stub{})))
	defer srv.Close()
	c := newClient(t, srv.URL)

	t.Run("200 typed response with valid token", func(t *testing.T) {
		got, err := c.GetNetworkNodeFeeWithResponse(context.Background(), bearer(token))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if got.JSON200 == nil {
			t.Fatalf("JSON200 is nil, status %d, body %s", got.StatusCode(), got.Body)
		}
		if got.JSON200.Status != "success" || got.JSON200.MaxNodeFee != 0.2 {
			t.Fatalf("unexpected body: %+v", got.JSON200)
		}
	})

	t.Run("401 with wrong token", func(t *testing.T) {
		got, err := c.GetNetworkNodeFeeWithResponse(context.Background(), bearer("nope"))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if got.StatusCode() != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", got.StatusCode())
		}
		if v := got.HTTPResponse.Header.Get("WWW-Authenticate"); v != "Bearer" {
			t.Fatalf("WWW-Authenticate = %q, want Bearer", v)
		}
	})

	t.Run("401 without Authorization header", func(t *testing.T) {
		got, err := c.GetNetworkNodeFeeWithResponse(context.Background())
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if got.StatusCode() != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", got.StatusCode())
		}
	})

	// The generated std mux is method-bound (GET only). The real snroute mux is not
	// (see README, Limitations); this test documents the generated behaviour only.
	t.Run("non-GET is not served by the handler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/network/node-fee", nil)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("POST status = %d, want 405", resp.StatusCode)
		}
	})
}

func TestErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(api.Handler(stub{fail: true}))
	defer srv.Close()
	c := newClient(t, srv.URL)

	got, err := c.GetNetworkNodeFeeWithResponse(context.Background())
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if got.JSON500 == nil {
		t.Fatalf("JSON500 is nil, status %d, body %s", got.StatusCode(), got.Body)
	}
	if got.JSON500.Status != "error" || got.JSON500.Error != "boom" {
		t.Fatalf("unexpected envelope: %+v", got.JSON500)
	}
}
