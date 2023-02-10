package hub

import (
	"fmt"
	"net"

	"gitlab.com/macroscope-lab/atomika/log"
	"gitlab.com/macroscope-lab/atomika/runtime"
)

type Runtime struct {
	ResendLimit int    `mapstructure:"resendLimit"`
	LocalAddr   string `mapstructure:"localAddr"`
	Join        string `mapstructure:"join"`
}

func (c *Runtime) Configure(key ...string) error {
	return runtime.Get(c, key...)
}

func (c *Runtime) Bind() {
	runtime.BindKeyToEnv("resendLimit", "NETRAU_RESEND_LIMIT")
	runtime.BindKeyToEnv("localAddr", "NETRAU_LOCAL_ADDR")
	runtime.BindKeyToEnv("join", "NETRAU_JOIN_ADDR")
}

func (c *Runtime) Validate() error {
	if c.ResendLimit == 0 {
		e := fmt.Errorf("%w:[%d]", ErrOptResendLimitInvalid, c.ResendLimit)
		log.Error(e.Error(), "service", "Hub", "runtime", "Validate")

		return ErrOptResendLimitInvalid
	}

	if _, _, err := net.SplitHostPort(c.LocalAddr); err != nil {
		e := fmt.Errorf("%w:[%s]", ErrOptLocalAddrInvalid, c.LocalAddr)
		log.Error(e.Error(), "service", "Hub", "runtime", "Validate")

		return ErrOptLocalAddrInvalid
	}

	return nil
}
