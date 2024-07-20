package auth_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meian/atgo/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	existDir   string
	noExistDir string
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	existDir = dir
	defer func() {
		os.RemoveAll(existDir)
	}()
	dir, err = os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	noExistDir = dir
	if err := os.RemoveAll(noExistDir); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestRead(t *testing.T) {
	type args struct {
		file string
	}
	type want struct {
		err      bool
		user     string
		password string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid1",
			args: args{
				file: filepath.Join(testDataDir(), "valid1"),
			},
			want: want{
				user:     "sample-user",
				password: "sample-password",
			},
		},
		{
			name: "valid2",
			args: args{
				file: filepath.Join(testDataDir(), "valid2"),
			},
			want: want{
				user:     "user2",
				password: "abcdefghijklmn",
			},
		},
		{
			name: "file not found",
			args: args{
				file: filepath.Join(testDataDir(), "not found"),
			},
			want: want{
				err: true,
			},
		},
		{
			name: "illegal base64",
			args: args{
				file: filepath.Join(testDataDir(), "illegal-base64"),
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			user, password, err := auth.Read(context.Background(), tt.args.file)
			t.Logf("user: %s", user)
			t.Logf("password: %s", password)
			if tt.want.err {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(tt.want.user, user)
			assert.Equal(tt.want.password, password)
		})
	}
}

func TestWrite(t *testing.T) {
	type args struct {
		file     string
		user     string
		password string
	}
	type want struct {
		err      bool
		testdata string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid1",
			args: args{
				file:     filepath.Join(existDir, "valid1"),
				user:     "sample-user",
				password: "sample-password",
			},
			want: want{
				testdata: "valid1",
			},
		},
		{
			name: "valid2",
			args: args{
				file:     filepath.Join(existDir, "valid2"),
				user:     "user2",
				password: "abcdefghijklmn",
			},
			want: want{
				testdata: "valid2",
			},
		},
		{
			name: "no parent directory",
			args: args{
				file:     filepath.Join(noExistDir, "credential"),
				user:     "user",
				password: "password",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			err := auth.Write(context.Background(), tt.args.file, tt.args.user, tt.args.password)
			if tt.want.err {
				assert.Error(err)
				return
			}

			if !assert.NoError(err) {
				return
			}
			if !assert.FileExists(tt.args.file) {
				return
			}
			cf, err := os.ReadFile(tt.args.file)
			require.NoError(err)

			wantfile := filepath.Join(testDataDir(), tt.want.testdata)
			require.FileExists(wantfile)
			testdata, err := os.ReadFile(wantfile)
			require.NoError(err)

			assert.Equal(string(testdata), string(cf))
		})
	}
}

func testDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	bp := filepath.Dir(b)
	return filepath.Join(bp, "testdata", "credentials")
}
