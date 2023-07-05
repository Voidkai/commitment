package pedersen_commitment

import (
	"crypto/rand"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"reflect"
	"testing"
)

func Test_pedersen_commiter_Commit(t *testing.T) {
	r := rand.Reader

	_, g, _ := bn256.RandomG1(r)
	_, h, _ := bn256.RandomG1(r)
	type fields struct {
		G *bn256.G1
		H *bn256.G1
	}
	type args struct {
		message     []byte
		blingFactor []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{name: "test1", fields: fields{g, h}, args: args{
			message:     []byte("hello world"),
			blingFactor: []byte("10"),
		}, want: pedersen_commiter{g, h}.Commit([]byte("hello world"), []byte("10"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := pedersen_commiter{
				G: tt.fields.G,
				H: tt.fields.H,
			}
			if got := pc.Commit(tt.args.message, tt.args.blingFactor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Commit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pedersen_commiter_Verify(t *testing.T) {
	r := rand.Reader
	_, g, _ := bn256.RandomG1(r)
	_, h, _ := bn256.RandomG1(r)

	type fields struct {
		G *bn256.G1
		H *bn256.G1
	}
	type args struct {
		commitment  []byte
		message     []byte
		blingFactor []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{name: "test1", fields: fields{g, h}, args: args{
			message:     []byte("hello world"),
			blingFactor: []byte("10"),
		}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := pedersen_commiter{
				G: tt.fields.G,
				H: tt.fields.H,
			}
			if got := pc.Verify(tt.args.commitment, tt.args.message, tt.args.blingFactor); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
