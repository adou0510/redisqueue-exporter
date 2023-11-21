package main

import (
	"reflect"
	"testing"
)

func TestParseRedisDSN(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisDSN
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test adou:123456@127.0.0.1:6379/0/testqueue",
			args: args{
				dsn: "adou:123456@127.0.0.1:6379/0/testqueue",
			},
			want: &RedisDSN{
				Addr:     "127.0.0.1:6379",
				Password: "adou:123456",
				DB:       0,
				Key:      "testqueue",
			},
			wantErr: false,
		},
		{
			name: "test 127.0.0.1:6379/0/testqueue",
			args: args{
				dsn: "127.0.0.1:6379/0/testqueue",
			},
			want: &RedisDSN{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       0,
				Key:      "testqueue",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRedisDSN(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRedisDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRedisDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}
