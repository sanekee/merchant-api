package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	type args struct {
		fn  http.HandlerFunc
		jwt string
	}
	okFunc := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "auth success",
			args: args{
				fn:  okFunc,
				jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMSIsIm5hbWUiOiJVc2VyIEEiLCJpYXQiOjE2NDA0MTY3MTN9.SbdB7XjwUDk2iNKegVPG7OEvodf5btXP1UjVCXXWHo0",
			},
			want: http.StatusOK,
		},
		{
			name: "invalid token",
			args: args{
				fn:  nil,
				jwt: "12345",
			},
			want: http.StatusForbidden,
		},
		{
			name: "invalid signature",
			args: args{
				fn:  nil,
				jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMSIsIm5hbWUiOiJVc2VyIEEiLCJpYXQiOjE2NDA0MTY3MTN9.douUzUc4Rbb7kAPgRPgZN9AUJTPZCCAI9vqJSdXsjrA",
			},
			want: http.StatusForbidden,
		},
		{
			name: "user id not found",
			args: args{
				fn:  nil,
				jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiVXNlciBBIiwiaWF0IjoxNjQwNDE2NzEzfQ.PPCvySNOlSDmjg7GWJk-89F6zhHZhebqBrd7Qcc9UqI",
			},
			want: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := JWTAuth(tt.args.fn)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tt.args.jwt))
			rec := httptest.NewRecorder()
			fn.ServeHTTP(rec, req)

			assert.Equal(t, tt.want, rec.Code)
		})
	}
}
