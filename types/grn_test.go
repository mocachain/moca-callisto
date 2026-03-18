package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types/resource"
)

func TestGRN_ParseFromString(t *testing.T) {
	type fields struct {
		resType    resource.ResourceType
		groupOwner sdk.AccAddress
		name       string
	}
	type args struct {
		res       string
		wildcards bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid bucket",
			fields{
				resType:    resource.RESOURCE_TYPE_GROUP,
				groupOwner: sdk.AccAddress{},
				name:       "test",
			},
			args{
				res:       "grn:g:0x70997970C51812dc3A010C7d01b50e0d17dc79C8:moca",
				wildcards: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GRN{
				resType:    tt.fields.resType,
				groupOwner: tt.fields.groupOwner,
				name:       tt.fields.name,
			}
			if err := r.ParseFromString(tt.args.res, tt.args.wildcards); (err != nil) != tt.wantErr {
				t.Errorf("GRN.ParseFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
