package repository

import (
	"database/sql"
	"encoding/json"
	"testing"

	"gin-enterprise/internal/db/sqlc"
)

func TestCategoryIconJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		value any
		want string
	}{
		{"null", nil, "null"},
		{"empty", "", `""`},
		{"value", "folder", `"folder"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var icon sql.NullString
			if err := icon.Scan(tc.value); err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(categoryModel(sqlc.Category{Icon: icon}))
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(data, &fields); err != nil {
				t.Fatal(err)
			}
			if got := string(fields["icon"]); got != tc.want {
				t.Fatalf("icon = %s, want %s", got, tc.want)
			}
		})
	}
}
