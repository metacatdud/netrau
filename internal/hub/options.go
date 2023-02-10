package hub

import (
	"net"
	"strconv"
)

type Option func(o *Options)

type Options struct {
	ResendLimit int
	BindAddr    string
	BindPort    int
	Join        string
}

func WithResendLimit(resendLimit int) Option {
	return func(o *Options) {
		o.ResendLimit = resendLimit
	}
}

func WithLocalAddr(localAddr string) Option {
	return func(o *Options) {
		h, p, _ := net.SplitHostPort(localAddr)
		o.BindAddr = h
		o.BindPort, _ = strconv.Atoi(p)
	}
}

func WithHubAddr(join string) Option {
	return func(o *Options) {
		o.Join = join
	}
}

func setOptions(opts ...Option) Options {
	defaults := defaultOptions()
	for _, opt := range opts {
		opt(&defaults)
	}

	return defaults
}

func defaultOptions() Options {
	return Options{
		ResendLimit: 3,
		BindPort:    0,
	}
}
