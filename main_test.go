package main

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_valid(t *testing.T) {
	type args struct {
		gocap string
		pkg   string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]error
		wantErr bool
	}{
		{
			name: "CapabilityNotProvided",
			args: args{
				gocap: `github.com/cugu/gocap/test/pkg ()

github.com/cugu/gocap/test/pkg/aws (network)
github.com/cugu/gocap/test/pkg/logging ()
github.com/cugu/gocap/test/pkg/sqlite (file)`,
				pkg: "github.com/cugu/gocap/test/pkg",
			},
			want: map[string][]error{
				"github.com/cugu/gocap/test/pkg/logging": {&CapabilityNotProvided{Capability: "execute"}},
			},
			wantErr: false,
		},
		{
			name: "UnnecessaryCapability",
			args: args{
				gocap: `github.com/cugu/gocap/test/pkg/aws (network, syscall)`,
				pkg:   "github.com/cugu/gocap/test/pkg/aws",
			},
			want: map[string][]error{
				"github.com/cugu/gocap/test/pkg/aws": {&UnnecessaryCapability{Capability: "syscall"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valid(tt.args.gocap, tt.args.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valid() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func valid(gocap string, pkg string) (map[string][]error, error) {
	f, err := parse(bytes.NewBufferString(gocap))
	if err != nil {
		return nil, err
	}

	tree, err := NewTree(pkg)
	if err != nil {
		return nil, err
	}

	return check(tree, f), nil
}
