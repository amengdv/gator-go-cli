package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amengdv/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("Login expect only 1 arguments")
    }

    user, err := s.db.GetUser(context.Background(), cmd.args[0])
    if err != nil {
        return err
    }

    err = s.cfg.SetUser(user.Name)
    if err != nil {
        return err
    }

    fmt.Printf("The user %v has been set\n", s.cfg.CurrentUserName)
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("Register expect only 1 arguments")
    }

    user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: cmd.args[0],
    })

    if err != nil {
        return errors.New("Failed to create users")
    }

    if err = s.cfg.SetUser(user.Name); err != nil {
        return err
    }

    fmt.Printf("The user %v was created\n", user.Name)
    log.Println(user)
    return nil
}

func handlerReset(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("Reset does not take any args")
    }

    if err := s.db.ResetDB(context.Background()); err != nil {
        return err
    }

    log.Println("Successfully resetting users db")

    return nil
}

func handlerUsers(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("Users does not take any args")
    }

    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return err
    }

    for _, user := range users {
        fmt.Print("* ")
        if user.Name == s.cfg.CurrentUserName {
            fmt.Printf("%v (current)\n", user.Name)
            continue
        }
        fmt.Println(user.Name)
    }

    return nil
}

func handlerAgg(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("agg does not take any args")
    }

    feed, err := fetchFeed("https://www.wagslane.dev/index.xml")
    if err != nil {
        return err
    }

    fmt.Println(feed)
    return nil
}

func handlerAddFeed(s *state, cmd command, currUser database.User) error {
    if len(cmd.args) != 2 {
        return errors.New("addfeed takes only 2 arguments")
    }

    name := cmd.args[0]
    url := cmd.args[1]

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: name,
        Url: url,
        UserID: currUser.ID,
    })

    if err != nil {
        return err
    }

    _, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: currUser.ID,
        FeedID: feed.ID,
    })
    if err != nil {
        log.Println("Create feed follow err")
        return err
    }

    fmt.Println(feed)
    return nil
}

func handlerFeeds(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return errors.New("feeds does not take any args")
    }

    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return err
    }

    fmt.Println("-----------------------------")
    for _, feed := range feeds {
        fmt.Printf("Name: %v\n", feed.Name)
        fmt.Printf("Url: %v\n", feed.Url)
        userName, err := s.db.GetUserNameByID(context.Background(), feed.UserID)
        if err != nil {
            return err
        }
        fmt.Printf("Created By: %v\n", userName)
        fmt.Println("-----------------------------")
    }

    return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return errors.New("follow only takes 1 args")
    }

    feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
    if err != nil {
        return err
    }

    feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID,
        FeedID: feed.ID,
    })
    if err != nil {
        return err
    }

    fmt.Printf("Follower: %v\n", feedFollow.UserName)
    fmt.Printf("Feed: %v\n", feedFollow.FeedName)

    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 0 {
        return errors.New("following does not take any args")
    }

    feedFollowByUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return err
    }

    fmt.Printf("Feeds followed by: %v\n", user.Name)

    for _, feeds := range feedFollowByUser {
        fmt.Printf("Name: %v\n", feeds.FeedName)
    }
    
    return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return errors.New("unfollow takes 1 arguments")
    }

    feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
    if err != nil {
        return err
    }

    err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
        UserID: user.ID,
        FeedID: feed.ID,
    })
    if err != nil {
        return err
    }

    log.Println("Successfully unfollow", feed.Name)

    return nil
}
