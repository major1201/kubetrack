package kube

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// Ping returns if the API server is online or not
func (c *ClientImpl) Ping() error {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second)
	defer cancel()

	content, err := c.GetDiscoveryClient().RESTClient().Get().AbsPath("/healthz").Timeout(time.Second).DoRaw(ctx)
	if err != nil {
		return errors.Wrap(err, "ping failed")
	}

	contentStr := string(content)
	if contentStr != "ok" {
		return errors.Errorf("ping response is not ok: %s", contentStr)
	}

	return nil
}
