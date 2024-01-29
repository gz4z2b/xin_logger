package xinlogger

import (
	"context"
	"testing"

	_ "embed"

	"github.com/gz4z2b/xinlogger/loggermid/cachelogger"
	"github.com/gz4z2b/xinlogger/loggermid/databaselogger"
	"github.com/gz4z2b/xinlogger/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
					LogPath: "./test.log",
					Type:    "zap",
					//EnableLevel: types.DebugLevel,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.args.conf)
			t.Log(err)
			err = logger.FlushTraceId()
			t.Log(err)
			logger.Info("test", types.NewField("第一个", "第一个内容"))
			logger.Debug("debug", types.NewField("debug", "debug内容"))
		})
	}
}

func TestGormMid(t *testing.T) {
	connection := mysql.Open("root:gz4z2b@tcp(127.0.0.1:3306)/webook?charset=utf8mb4")
	db, err := gorm.Open(connection, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	logger, err := NewLogger(types.Conf{
		LogPath: "./sql.log",
		Type:    "zap",
		//EnableLevel: types.DebugLevel,
	})
	if err != nil {
		t.Error(err)
	}
	err = logger.FlushTraceId()
	if err != nil {
		t.Error(err)
	}
	db.Use(databaselogger.NewSqlLoggerMid(func(sql string, rows, seconds int) {
		logger.Info("数据库操作", types.NewField("sql", sql), types.NewField("影响行数", rows), types.NewField("用时(毫秒)", seconds))
	}))

	var user User
	db.First(&user)
	user.Email = "test"
	db.Save(&user)
	db.Delete(&user)

	newUser := User{
		Email:    "test",
		Password: "123456",
	}
	db.Save(&newUser)
}

//go:embed test/test_redis.lua
var luaSetCode string

func TestRedisLogger(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "gz4z2b", // no password set
		DB:       2,        // use default DB
	})
	logger, err := NewLogger(types.Conf{
		LogPath: "./redis.log",
		Type:    "zap",
		MaxSize: 1,
		//EnableLevel: types.DebugLevel,
	})

	if err != nil {
		t.Error(err)
	}
	r.AddHook(cachelogger.NewRedisLoggerHook(func(cmd string, milliSeconds int) {
		logger.Info("redis操作", types.NewField("命令", cmd), types.NewField("用时(毫秒)", milliSeconds))
	}))
	result, err := r.Get(context.Background(), "string").Bytes()
	t.Log(result, err)

	// for i := 0; i < 100000; i++ {
	// 	res, err := r.Eval(context.Background(), luaSetCode, []string{"test"}, 123).Int()
	// 	t.Log(res, err)
	// }

}

type User struct {
	Id       uint64 `gorm:"primaryKey,not null,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Createtime int64 `gorm:"autoCreateTime:milli"`
	Updatetime int64 `gorm:"autoUpdateTime:milli"`
	Deletetime int64
}

func (u User) TableName() string {
	return "t_user"
}
