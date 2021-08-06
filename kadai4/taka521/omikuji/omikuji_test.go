package omikuji

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestHandler(t *testing.T) {
	original := Execute
	defer func() {
		Execute = original
	}()

	type want struct {
		res  *response
		code int
	}
	tests := []struct {
		name  string
		want  want
		setup func()
	}{
		{
			name: "200 OK",
			want: want{
				res:  &response{Result: "大吉"},
				code: http.StatusOK,
			},
			setup: func() {
				Execute = func(t time.Time) *response {
					return &response{Result: "大吉"}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			w := httptest.NewRecorder()
			Handler(w, nil)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			if tt.want.code != res.StatusCode {
				t.Errorf("mismatch StatusCode: want = %d, got = %d", tt.want.code, res.StatusCode)
			}

			got := &response{}
			_ = json.Unmarshal(body, got)
			if diff := cmp.Diff(tt.want.res, got); diff != "" {
				t.Errorf("mismatch response body (-want +got):\n%s", diff)
			}
		})
	}
}
