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

func TestSamczsunGateway_GetEventTextSignature(t *testing.T) {
	type args struct {
		eventSign string
	}
	tests := []struct {
		name    string
		args    args
		want    *SamczsunResp
		wantErr bool
	}{
		{
			name: "Fetch text signature for event success - Transfer",
			args: args{
				eventSign: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			},
			want: &SamczsunResp{
				Ok: true,
				Result: struct {
					Event    map[string][]SamczsunTextResult `json:"event"`
					Function map[string][]SamczsunTextResult `json:"function"`
				}{
					Event: map[string][]SamczsunTextResult{
						"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": {{
							Filtered: false,
							Name:     "Transfer(address,address,uint256)",
						}},
					},
					Function: map[string][]SamczsunTextResult{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewSamczsunGateway()
			got, err := g.GetEventTextSignature(tt.args.eventSign)
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

func TestSamczsunGateway_GetFunctionTextSignature(t *testing.T) {
	type args struct {
		functionSign string
	}
	tests := []struct {
		name    string
		args    args
		want    *SamczsunResp
		wantErr bool
	}{
		{
			name: "Fetch text signature for function success - Transfer",
			args: args{
				functionSign: "0xa9059cbb",
			},
			want: &SamczsunResp{
				Ok: true,
				Result: struct {
					Event    map[string][]SamczsunTextResult `json:"event"`
					Function map[string][]SamczsunTextResult `json:"function"`
				}{
					Event: map[string][]SamczsunTextResult{},
					Function: map[string][]SamczsunTextResult{
						"0xa9059cbb": {{
							Filtered: false,
							Name:     "transfer(address,uint256)",
						}},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewSamczsunGateway()
			got, err := g.GetFunctionTextSignature(tt.args.functionSign)
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
