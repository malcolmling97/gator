package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) >= 1 {
		if n, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = n
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	fmt.Println("Posts:")
	for _, post := range posts {
		fmt.Printf("- %s (%s)\n", post.Title, post.Url)
	}
	return nil
}
