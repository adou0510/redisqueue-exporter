package main

import (
	"reflect"
	"testing"
)

func TestCheckLen(t *testing.T) {
	type args struct {
		dsn RedisDSN
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test testqueue",
			args: args{
				dsn: RedisDSN{
					Addr:     "127.0.0.1:6379",
					Password: "",
					DB:       0,
					Key:      "testqueue",
				},
			},
			want: map[string]int64{
				"testqueue": 1,
			},
			wantErr: false,
		},
		{
			name: "test test*",
			args: args{
				dsn: RedisDSN{
					Addr:     "127.0.0.1:6379",
					Password: "",
					DB:       0,
					Key:      "test*",
				},
			},
			want: map[string]int64{
				"test":        1,
				"testqueue":   1,
				"test:queue1": 1,
				"test:queue2": 1,
				"test:queue3": 1,
				"test:queue4": 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckLen(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
