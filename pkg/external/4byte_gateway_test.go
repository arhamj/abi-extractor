package external

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestFourByteGateway_GetEventTextSignature(t *testing.T) {
	type args struct {
		eventSign string
	}
	tests := []struct {
		name    string
		args    args
		want    *FourBytesResp
		wantErr bool
	}{
		{
			name: "Fetch text signature for event success - Transfer",
			args: args{
				eventSign: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			},
			want: &FourBytesResp{
				Count: 1,
				Results: []TextSignResult{
					{
						TextSignature: "Transfer(address,address,uint256)",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFourByteGateway()
			got, err := f.GetEventTextSignature(tt.args.eventSign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventTextSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventTextSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFourByteGateway_GetFunctionTextSignature(t *testing.T) {
	type args struct {
		eventSign string
	}
	tests := []struct {
		name    string
		args    args
		want    *FourBytesResp
		wantErr bool
	}{
		{
			name: "Fetch text signature for function success - swapTokensForExactTokens",
			args: args{
				eventSign: "0x8803dbee",
			},
			want: &FourBytesResp{
				Count: 1,
				Results: []TextSignResult{
					{
						TextSignature: "swapTokensForExactTokens(uint256,uint256,address[],address,uint256)",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFourByteGateway()
			got, err := f.GetFunctionTextSignature(tt.args.eventSign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFunctionTextSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFunctionTextSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}
