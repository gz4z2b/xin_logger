/*
 * @Author: p_hanxichen
 * @Date: 2024-01-23 11:09:29
 * @LastEditors: p_hanxichen
 * @FilePath: /xinlogger/middleware/database/gormmidlogger/plugin.go
 * @Description: gorm中间件sql日志
 *
 * Copyright (c) 2024 by gdtengnan, All Rights Reserved.
 */

package gormmidlogger

import (
	"time"

	"gorm.io/gorm"
)

type SqlLoggerMid struct {
	// 日志记录，sql，影响行数，执行时间(毫秒)
	LogFunc func(sql string, rows int, milliSeconds int)
}

func NewSqlLoggerMid(logFunc func(sql string, rows int, seconds int)) gorm.Plugin {
	return &SqlLoggerMid{
		LogFunc: logFunc,
	}
}

func (l *SqlLoggerMid) Name() string {
	return "SqlLogger"
}

func (l *SqlLoggerMid) Initialize(db *gorm.DB) error {
	// 注册所有操作的前后操作
	db.Callback().Create().Before("gorm:before_create").Register("sql_logger_before", l.Before())
	db.Callback().Delete().Before("gorm:before_delete").Register("sql_logger_before", l.Before())
	db.Callback().Query().Before("gorm:query").Register("sql_logger_before", l.Before())
	db.Callback().Update().Before("gorm:setup_reflect_value").Register("sql_logger_before", l.Before())
	db.Callback().Row().Before("gorm:row").Register("sql_logger_before", l.Before())

	db.Callback().Create().After("gorm:after_create").Register("sql_logger_after", l.After())
	db.Callback().Delete().After("gorm:after_delete").Register("sql_logger_after", l.After())
	db.Callback().Query().After("gorm:after_query").Register("sql_logger_after", l.After())
	db.Callback().Update().After("gorm:after_update").Register("sql_logger_after", l.After())
	db.Callback().Row().After("gorm:row").Register("sql_logger_after", l.After())
	return nil
}

func NewSqlLoggerMidBuilder(logFunc func(sql string, rows int, seconds int)) *SqlLoggerMid {
	return &SqlLoggerMid{
		LogFunc: logFunc,
	}
}

// Before 记录开始操作的时间
func (l *SqlLoggerMid) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		db.InstanceSet("startTime", time.Now())
	}
}

// After 计算操作时间以及记录日志
func (l *SqlLoggerMid) After() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		rows := int(db.Statement.RowsAffected)
		startTime, isExist := db.InstanceGet("startTime")
		if !isExist {
			return
		}
		cost := time.Since(startTime.(time.Time)).Microseconds()
		l.LogFunc(sql, rows, int(cost))
	}
}
