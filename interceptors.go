package colorlog

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"

	callTypeUnary     = "unary call"
	callTypeStreaming = "streaming call"
)

type logPayload struct {
	startTime  time.Time
	callType   string
	fullMethod string
	statusCode codes.Code
	err        error
}

func (l logPayload) statusCodeColor() string {
	switch {
	case l.statusCode >= codes.OK && l.statusCode < codes.Canceled:
		return green
	case l.statusCode >= codes.Canceled && l.statusCode < codes.Unknown:
		return white
	case l.statusCode >= codes.Unknown && l.statusCode < codes.Unimplemented || l.statusCode == codes.Unauthenticated:
		return yellow
	default:
		return red
	}
}

func (l logPayload) callTypeColor() string {
	switch l.callType {
	case callTypeUnary:
		return blue
	case callTypeStreaming:
		return cyan
	default:
		return reset
	}
}

func (l logPayload) formatter(isOutputColor bool) string {
	var statusColor, callTypeColor, resetColor string
	if isOutputColor {
		statusColor = l.statusCodeColor()
		callTypeColor = l.callTypeColor()
		resetColor = reset
	}
	var errMsg string
	if l.err != nil {
		errMsg = "\nError: " + l.err.Error()
	}
	return fmt.Sprintf("[gRPC-Go] %v |%s %18s %s| %13v |%s %-14s %s %s%s\n",
		l.startTime.Format("2006/01/02 - 15:04:05"),
		statusColor, l.statusCode.String(), resetColor,
		time.Since(l.startTime),
		callTypeColor, l.callType, resetColor,
		l.fullMethod,
		errMsg,
	)
}

// UnaryServerInterceptor returns a new unary server interceptors.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logPayload{
			startTime:  time.Now(),
			callType:   callTypeUnary,
			fullMethod: info.FullMethod,
		}
		res, err := handler(ctx, req)
		if !o.shouldLog(l.fullMethod, err) {
			return res, err
		}
		l.err = err
		l.statusCode = o.codeFunc(err)
		fmt.Fprint(o.output, l.formatter(o.isOutputColor))
		return res, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOpt(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		l := logPayload{
			startTime:  time.Now(),
			callType:   callTypeUnary,
			fullMethod: info.FullMethod,
		}
		err := handler(srv, ss)
		if !o.shouldLog(l.fullMethod, err) {
			return err
		}
		l.err = err
		l.statusCode = o.codeFunc(err)
		fmt.Fprint(o.output, l.formatter(o.isOutputColor))
		return err
	}
}
