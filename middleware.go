package main

import (
	"context"
	"fmt"
	"github.com/vigneshsekar314/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(st *state, cm command) error {
		ctx := context.Background()
		user, err := st.db.GetUser(ctx, st.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user is not registered, register the user first.")
		}
		if err := handler(st, cm, user); err != nil {
			return err
		}
		return nil
	}
}
