/*
 * @Author: p_hanxichen
 * @Date: 2023-11-27 15:05:55
 * @LastEditors: p_hanxichen
 * @FilePath: /xinlogger/loggermidd/http/gin_mid.go
 * @Description: gin的访问日志中间件
 *
 * Copyright (c) 2023 by gdtengnan, All Rights Reserved.
 */

package httplogger

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
)

// 日志内容
type AccessLog struct {
	Method   string
	Header   map[string][]string
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}

type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody *atomic.Bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		allowReqBody:  atomic.NewBool(false),
		allowRespBody: atomic.NewBool(false),
		loggerFunc:    fn,
	}
}

// AllowReqBody 是否允许记录请求体
func (b *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	b.allowReqBody.Store(ok)
	return b
}

// AllowRespBody 是否允许记录响应体
func (b *MiddlewareBuilder) AllowRespBody(ok bool) *MiddlewareBuilder {
	b.allowRespBody.Store(ok)
	return b
}

// Build 构建中间件
func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Header: ctx.Request.Header,
			Url:    url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			if len(body) > 1024 {
				al.ReqBody = string(body[:1024])
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody.Load() {
			ctx.Writer = &responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		ctx.Next()

		defer func() {
			al.Duration = time.Since(start).String()
			al.Status = ctx.Writer.Status()
			if len(al.RespBody) > 1024 {
				al.RespBody = al.RespBody[:1024]
			}

			b.loggerFunc(ctx, al)
		}()

	}
}

// 重写ctx的writer
type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

// 记录header
func (w *responseWriter) WriteHeader(code int) {
	w.al.Status = code
	w.ResponseWriter.WriteHeader(code)
}

// 记录返回体
func (w *responseWriter) Write(data []byte) (n int, err error) {
	dataString := string(data)
	if len(dataString) > 1024 {
		dataString = dataString[:1024]
	}
	w.al.RespBody = dataString
	return w.ResponseWriter.Write(data)
}

// 返回体是string的情况
func (w *responseWriter) WriteString(s string) (n int, err error) {
	if len(s) > 1024 {
		s = s[:1024]
	}
	w.al.RespBody = s
	return w.ResponseWriter.WriteString(s)
}
