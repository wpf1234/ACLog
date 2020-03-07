package utils

import (
	"reflect"
	"testing"
)

func TestDistinct(t *testing.T) {
	type args struct {
		tempItem []string
	}
	tests := []struct {
		name       string
		args       args
		wantNewArr []string
	}{
		// TODO: Add test cases.
		{
			name:"distinct",
			args:args{tempItem: []string{
				"a","b","a","c","b",
			}},
			wantNewArr: []string{
				"a","b","c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewArr := Distinct(tt.args.tempItem); !reflect.DeepEqual(gotNewArr, tt.wantNewArr) {
				t.Errorf("Distinct() = %v, want %v", gotNewArr, tt.wantNewArr)
			}
		})
	}
}