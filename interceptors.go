package colorlog

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
)

type logPayload struct {
	startTime     time.Time
	fullMethod    string
	statusCode    codes.Code
	err           error
	isOutputColor bool
}

func (l *logPayload) detectStatusCodeFromError(err error) {
	l.err = err
	l.statusCode = status.Code(l.err)
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

func (l logPayload) print() {
	var statusColor, methodColor, resetColor string
	if l.isOutputColor {
		statusColor = l.statusCodeColor()
		methodColor = blue
		resetColor = reset
	}
	var errMsg string
	if l.err != nil {
		errMsg = l.err.Error()
	}
	fmt.Printf("[gRPC-Go] %v |%s %19s %s| %13v |%s %s %s\n%s",
		l.startTime.Format("2006/01/02 - 15:04:05"),
		statusColor, l.statusCode.String(), resetColor,
		time.Since(l.startTime),
		methodColor, l.fullMethod, resetColor,
		errMsg,
	)
}

// UnaryServerInterceptor returns a new unary server interceptors.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	isTerm := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := logPayload{
			startTime:     time.Now(),
			fullMethod:    info.FullMethod,
			isOutputColor: isTerm,
		}
		res, err := handler(ctx, req)
		l.detectStatusCodeFromError(err)
		l.print()
		return res, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	isTerm := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		l := logPayload{
			startTime:     time.Now(),
			fullMethod:    info.FullMethod,
			isOutputColor: isTerm,
		}
		err := handler(srv, ss)
		l.detectStatusCodeFromError(err)
		l.print()
		return err
	}
}
