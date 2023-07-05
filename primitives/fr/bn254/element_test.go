package bn254

import (
	"reflect"
	"testing"
)

func TestOne(t *testing.T) {
	var tests []struct {
		name string
		want Element
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := One(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("One() = %v, want %v", got, tt.want)
			}
		})
	}
}
