package youtube

import (
	"fmt"
	"net/http"

	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// Search searches YouTube for videos
func search(ctx *exrouter.Context) string {
	developerKey := viper.GetString("youtube.key")
	args := ctx.Args.After(1)

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube
	call := service.Search.List("id,snippet").
		Q(args).
		MaxResults(1) // Only fetch one result
	response, err := call.Do()
	if err != nil {
		log.Info(err)
	}

	videoID := response.Items[0].Id.VideoId
	videoTitle := response.Items[0].Snippet.Title
	var reply string
	if videoID == "" {
		reply = "Search returned nothing."
	} else {
		// Format the reply as <link> - <title>
		reply = fmt.Sprintf("https://youtube.com/watch?v=%s - %s", videoID, videoTitle)
	}
	return reply
}

// Search replies with the search result
func Search(ctx *exrouter.Context) {
	ctx.Reply(search(ctx))
}
