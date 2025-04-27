package session

import (
	"bskyway/config"
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

func NewSession(ctx context.Context) (*xrpc.Client, error) {

	client := &xrpc.Client{Host: "https://bsky.social"}
	conf := config.Instance()

	sess, err := atproto.ServerCreateSession(ctx, client, &atproto.ServerCreateSession_Input{
		Identifier: conf.Identifier,
		Password:   conf.PassWord,
	})

	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	client.Auth = &xrpc.AuthInfo{
		AccessJwt: sess.AccessJwt,
		Did:       sess.Did,
		Handle:    sess.Handle,
	}

	return client, nil
}
