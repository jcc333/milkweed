package send

import (
	"context"
	"log"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	xrpc "github.com/bluesky-social/indigo/xrpc"

	"github.com/jcc333/milkweed/internal/config"
	"github.com/jcc333/milkweed/internal/poasts"
	"github.com/jcc333/milkweed/internal/state"
)

// One which sends.
type Sender interface {
	// Send skeet.
	Send(*poasts.Poast, *state.State) error
}

// Send - to bluesky.
type BlueskySender struct {
	// Configuration values.
	cfg *config.Config

	// Underlying XRPC Client
	client *xrpc.Client
}

func (c *BlueskySender) Connect(ctx context.Context) error {
	sessionInput := &atproto.ServerCreateSession_Input{
		Identifier: c.cfg.Username,
		Password:   c.cfg.Password,
	}

	session, err := atproto.ServerCreateSession(ctx, c.client, sessionInput)

	if err != nil {
		return err
	}

	// Access Token is used to make authenticated requests
	// Refresh Token allows to generate a new Access Token
	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	return nil
}

// Create a new sender ready to send.
func New(cfg *config.Config) *BlueskySender {
	return &BlueskySender{
		cfg:    cfg,
		client: &xrpc.Client{},
	}

}

// Close any connections to avoid any oddball resource leaks.
func (s *BlueskySender) Close() error {
	return nil
}

// Send skeet.
func (c *BlueskySender) Send(ctx context.Context, p *poasts.Poast, s *state.State) error {
	logger := log.Default()
	post := bsky.FeedPost{}
	post.LexiconTypeID = "app.bsky.feed.post"
	post.CreatedAt = time.Now().Format(time.RFC3339)

	var FeedPost_Embed bsky.FeedPost_Embed
	FeedPost_Embed.EmbedExternal = &bsky.EmbedExternal{
		LexiconTypeID: "app.bsky.embed.external",
		External: &bsky.EmbedExternal_External{
			Title:       p.Title,
			Uri:         p.GUID,
			Description: p.Description,
			Thumb:       nil,
		},
	}
	post.Embed = &FeedPost_Embed

	input := &atproto.RepoCreateRecord_Input{
		// collection: The NSID of the record collection.
		Collection: "app.bsky.feed.post",
		// repo: The handle or DID of the repo (aka, current account).
		Repo: c.client.Auth.Did,
		// record: The record itself. Must contain a $type field.
		Record: &lexutil.LexiconTypeDecoder{Val: &post},
	}

	response, err := atproto.RepoCreateRecord(ctx, c.client, input)
	if err != nil {
		return err
	}
	if response != nil {
		logger.Println("Sent skeet to '", response.Uri, "' with CID '", response.Commit.Cid, "'")
	}

	return nil
}
