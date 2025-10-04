package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vigneshsekar314/gator/internal/config"
	"github.com/vigneshsekar314/gator/internal/database"
	"os"
	"time"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdList map[string]func(*state, command) error
}

func GetNewState(config *config.Config) state {
	return state{
		config: config,
	}
}

func GetNewCommands() commands {
	return commands{
		cmdList: make(map[string]func(*state, command) error),
	}
}

func Run(state *state, cmds *commands) error {
	args := os.Args
	if len(args) < 2 {
		return errors.New("Not enough arguments")
	}
	cmd := command{}
	cmd.name = args[1]
	cmd.arguments = args[2:]
	if err := cmds.run(state, cmd); err != nil {
		return err
	}
	return nil
}

func (cmds *commands) run(s *state, cmd command) error {
	if s == nil || s.config == nil {
		return errors.New("state does not exists")
	}
	fnToRun, ok := cmds.cmdList[cmd.name]
	if !ok {
		return errors.New("command is not registered")
	}
	if err := fnToRun(s, cmd); err != nil {
		return err
	}
	return nil
}

func (cmds *commands) register(name string, fn func(*state, command) error) {
	cmds.cmdList[name] = fn
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Login command expects username")
	}
	usrNm := cmd.arguments[0]
	usr, err := s.db.GetUser(context.Background(), usrNm)
	if err != nil {
		return fmt.Errorf("user is not registered, register the user first.")
	}
	if usr.Name != usrNm {
		return errors.New("User is not found in the database")
	}
	if err := s.config.SetUser(usrNm); err != nil {
		return err
	}
	fmt.Printf("user has been set to %s\n", usrNm)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Register command exxpects username")
	}
	timeNow := time.Now()
	usrNm := cmd.arguments[0]
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      usrNm,
	})
	if err != nil {
		return err
	}
	if err := s.config.SetUser(usrNm); err != nil {
		return err
	}
	fmt.Printf("User is created. User ID: %v, createdAt: %v, updatedAt: %v, Name: %s\n", usr.ID, usr.CreatedAt, usr.UpdatedAt, usr.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.DeleteUser(ctx); err != nil {
		return err
	}
	if err := s.db.DeleteFeeds(ctx); err != nil {
		return err
	}
	fmt.Println("All users and feeds have been deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return errors.New("Atleast two Arguments should be present for name and url")
	}
	ctx := context.Background()
	var feedId uuid.UUID
	var feedName string
	var feedNew database.Feed
	feedExists, err := s.db.GetFeedsByUrl(ctx, cmd.arguments[1])
	if err != nil { // feed is not available, create a new feed
		timeNow := time.Now()
		var err2 error
		feedNew, err2 = s.db.CreateFeed(ctx, database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      cmd.arguments[0],
			Url:       cmd.arguments[1],
		})
		if err2 != nil {
			return err
		}
		feedId = feedNew.ID
		feedName = feedNew.Name
	} else {
		// feed exists
		feedId = feedExists.ID
		feedName = feedExists.Name
	}
	_, err = createFeedFollow(s, user.ID, feedId, ctx)
	if err != nil {
		return err
	}
	fmt.Printf("id: %s\n", feedId)
	fmt.Printf("name: %s\n", feedName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUsername(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("feed name: %s | url: %s | username: %s\n", feed.Name, feed.Url, feed.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return errors.New("follow command expects one argument for url")
	}
	url := cmd.arguments[0]
	ctx := context.Background()
	feeds, err := s.db.GetFeedsByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("Feed does not exists. Add feed url by using addfeed command.")
	}
	feed_follow, err := createFeedFollow(s, user.ID, feeds.ID, ctx)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %v, CreatedAt: %s, UpdatedAt: %s, UserId: %v, FeedId: %v, Username: %s, Feedname: %s\n", feed_follow.FeedFollowsID, feed_follow.CreatedAt, feed_follow.UpdatedAt, feed_follow.UserID, feed_follow.FeedID, feed_follow.Username, feed_follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	userfeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Feeds for %s:\n", user.Name)
	for _, feed := range userfeeds {
		fmt.Printf(" * %s\n", feed.FeedName)
	}
	return nil
}

func createFeedFollow(s *state, userId uuid.UUID, feedId uuid.UUID, ctx context.Context) (database.CreateFeedFollowRow, error) {
	timeNow := time.Now()
	feed_follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		UserID:    userId,
		FeedID:    feedId,
	})
	if err != nil {
		return database.CreateFeedFollowRow{}, err
	}
	return feed_follow, nil
}
