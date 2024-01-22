package xinlogger

import (
	"testing"

	"github.com/gz4z2b/xinlogger/types"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		conf types.Conf
	}
	tests := []struct {
		name    string
		args    args
		want    types.XinLogger
		wantErr error
	}{
		{
			name: "正常",
			args: args{
				conf: types.Conf{
					LogPath:     "./test.log",
					Type:        "zap",
					EnableLevel: types.DebugLevel,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.args.conf)
			logger.Info("test", types.NewField("第一个", "第一个内容"))
			logger.Debug("debug", types.NewField("debug", "debug内容"))
			t.Log(err)
		})
	}
}
