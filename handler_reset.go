package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to delete users:", err)
		os.Exit(1) // explicitly required
	}
	fmt.Println("All users deleted successfully.")
	return nil
}
