package handlers

import (
	"blogo/internal/database"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func FollowFeed(app Application, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: follow [feed-url]")
	}

	feed, err := app.GetDB().GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	_, err = app.GetDB().CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Printf("Now following: %s\n", feed.Name)
	return nil
}

func ListFollowing(app Application, cmd Command, user database.User) error {
	follows, err := app.GetDB().GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("Not following any feeds")
		return nil
	}

	fmt.Println("Following:")
	for _, follow := range follows {
		fmt.Printf("- %s (%s)\n", follow.FeedName, follow.FeedUrl)
	}
	return nil
}

func getCurrentUser(app Application) database.User {
	currentUser := app.GetConfig().CurrentUser
	user, err := app.GetDB().GetUser(context.Background(), currentUser)
	if err != nil {
		log.Fatal("No logged in user:", err)
	}
	return user
}
