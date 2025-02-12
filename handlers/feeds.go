package handlers

import (
	"blogo/internal/database"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddFeed(app Application, args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: addfeed [name] [url]")
	}

	currentUser := app.GetConfig().CurrentUser
	if currentUser == "" {
		log.Fatal("No user logged in")
	}

	// Get current user ID
	user, err := app.GetDB().GetUser(context.Background(), currentUser)
	if err != nil {
		log.Fatalf("Error getting current user: %v", err)
	}

	feed, err := app.GetDB().CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      args[0],
		Url:       args[1],
		UserID:    user.ID,
	})

	if err != nil {
		log.Fatalf("Error creating feed: %v", err)
	}

	fmt.Printf("Feed created successfully:\nID: %s\nName: %s\nURL: %s\nUser ID: %s\n",
		feed.ID,
		feed.Name,
		feed.Url,
		feed.UserID,
	)
}

func ListFeeds(app Application) {
	feeds, err := app.GetDB().GetFeedsWithUsers(context.Background())
	if err != nil {
		log.Fatal("Error fetching feeds:", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		fmt.Printf("- %s\n  URL: %s\n  Added by: %s\n",
			feed.FeedName,
			feed.Url,
			feed.UserName,
		)
	}
}
