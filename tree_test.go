package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTree(t *testing.T) {
	type args struct {
		start string
	}
	tests := []struct {
		name    string
		args    args
		want    *Tree
		wantErr bool
	}{
		{
			name: "github.com/cugu/gocap/test/pkg",
			args: args{
				start: "github.com/cugu/gocap/test/pkg",
			},
			want: &Tree{
				root: "github.com/cugu/gocap/test/pkg",
				nodes: map[string]*Node{
					"github.com/cugu/gocap/test/pkg":         {capabilities: NewSet(), imports: NewSet("github.com/cugu/gocap/test/pkg/aws", "github.com/cugu/gocap/test/pkg/logging", "github.com/cugu/gocap/test/pkg/sqlite")},
					"github.com/cugu/gocap/test/pkg/aws":     {capabilities: NewSet("network"), imports: NewSet()},
					"github.com/cugu/gocap/test/pkg/logging": {capabilities: NewSet("execute"), imports: NewSet()},
					"github.com/cugu/gocap/test/pkg/sqlite":  {capabilities: NewSet("file"), imports: NewSet()},
				},
			},
			wantErr: false,
		},
		{
			name: "github.com/cugu/gocap/test/pkg/aws",
			args: args{
				start: "github.com/cugu/gocap/test/pkg/aws",
			},
			want: &Tree{
				root: "github.com/cugu/gocap/test/pkg/aws",
				nodes: map[string]*Node{
					"github.com/cugu/gocap/test/pkg/aws": {capabilities: NewSet("network"), imports: NewSet()},
				},
			},
			wantErr: false,
		},
		{
			name: "not existing",
			args: args{
				start: "nope",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTree(tt.args.start)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
