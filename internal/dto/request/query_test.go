package request_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"gin-enterprise/internal/dto/request"
)

func TestSiteQueryParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := humagin.New(engine, huma.DefaultConfig("test", "1"))
	var got request.ListSitesInput
	huma.Register(api, huma.Operation{Method: "GET", Path: "/sites"}, func(_ context.Context, input *request.ListSitesInput) (*struct{ Body string }, error) {
		got = *input
		return &struct{ Body string }{Body: "ok"}, nil
	})
	for _, tc := range []struct {
		query  string
		status int
		check  func(*testing.T)
	}{
		{"", 200, func(t *testing.T) {
			if got.IsPinned.Pointer() != nil || got.Skip() != 0 || got.Limit() != 10 {
				t.Fatalf("unexpected defaults: %+v", got)
			}
		}},
		{"?title=Go&is_pinned=false&page=3&page_size=20&offset=0", 200, func(t *testing.T) {
			if got.Title != "Go" || got.IsPinned.Pointer() == nil || *got.IsPinned.Pointer() {
				t.Fatalf("lost query values: %+v", got)
			}
			if got.Skip() != 0 || got.PageNumber() != 1 || got.Limit() != 20 {
				t.Fatalf("explicit offset=0 ignored: %+v", got)
			}
		}},
		{"?page=3&page_size=20", 200, func(t *testing.T) {
			if got.Skip() != 40 || got.PageNumber() != 3 {
				t.Fatalf("incorrect page: %+v", got)
			}
		}},
		{"?category_id=bad", 422, nil},
		{"?user_id=bad", 422, nil},
		{"?is_pinned=bad", 422, nil},
		{"?page_size=101", 422, nil},
		{"?offset=-1", 422, nil},
	} {
		t.Run(tc.query, func(t *testing.T) {
			res := httptest.NewRecorder()
			engine.ServeHTTP(res, httptest.NewRequest("GET", "/sites"+tc.query, nil))
			if res.Code != tc.status {
				t.Fatalf("status=%d want=%d body=%s", res.Code, tc.status, res.Body.String())
			}
			if tc.check != nil {
				tc.check(t)
			}
		})
	}
	spec := api.OpenAPI().Paths["/sites"].Get
	for _, p := range spec.Parameters {
		if p.Name == "is_pinned" && p.Schema.Type != "boolean" {
			t.Fatalf("is_pinned schema type = %s", p.Schema.Type)
		}
	}
}
