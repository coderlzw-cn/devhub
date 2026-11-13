package repository

import (
	"database/sql"
	"encoding/json"
	"testing"

	"gin-enterprise/internal/db/sqlc"
)

func TestSiteLogoURLJSON(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
		want  string
	}{
		{"null", nil, "null"},
		{"empty", "", `""`},
		{"value", "folder", `"folder"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var logoURL sql.NullString
			if err := logoURL.Scan(tc.value); err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(siteModel(sqlc.Bookmark{LogoUrl: logoURL}))
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if got := string(fields["logoUrl"]); got != tc.want {
				t.Fatalf("logoUrl = %s, want %s", got, tc.want)
			}
		})
	}
}
