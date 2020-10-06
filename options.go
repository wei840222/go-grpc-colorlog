package colorlog

import (
	"io"
	"os"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/mattn/go-isatty"
)

type options struct {
	output        io.Writer
	isOutputColor bool
	shouldLog     grpc_logging.Decider
	codeFunc      grpc_logging.ErrorToCode
}

var defaultOptions = &options{
	output:        os.Stdout,
	isOutputColor: true,
	shouldLog:     grpc_logging.DefaultDeciderMethod,
	codeFunc:      grpc_logging.DefaultErrorToCode,
}

func evaluateOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	if w, ok := optCopy.output.(*os.File); !ok || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		optCopy.isOutputColor = false
	}
	return optCopy
}

// Option function for set options
type Option func(*options)

// WithDecider customizes the function for deciding if the gRPC interceptor logs should log.
func WithDecider(f grpc_logging.Decider) Option {
	return func(o *options) {
		o.shouldLog = f
	}
}

// WithErrorToCode customizes the function for mapping errors to error codes.
func WithErrorToCode(f grpc_logging.ErrorToCode) Option {
	return func(o *options) {
		o.codeFunc = f
	}
}
