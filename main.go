package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gpanaretou/Gator/internal/config"
	"github.com/gpanaretou/Gator/internal/database"
	_ "github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	available map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	switch name {
	case "login":
		c.available[name] = f
	case "register":
		c.available[name] = f
	case "reset":
		c.available[name] = f
	case "users":
		c.available[name] = f
	case "agg":
		c.available[name] = f
	case "addfeed":
		c.available[name] = f
	case "feeds":
		c.available[name] = f
	case "follow":
		c.available[name] = f
	case "following":
		c.available[name] = f
	case "unfollow":
		c.available[name] = f
	}

}

func (c *commands) run(s *state, cmd command) error {
	_, ok := c.available[cmd.name]
	if !ok {
		err := fmt.Errorf("%v is not a command", cmd.name)
		return err
	}

	switch cmd.name {
	case "login":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "register":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "reset":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "users":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "agg":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "addfeed":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "feeds":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "follow":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "following":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "unfollow":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("-- '%v' command not implemented", cmd.name)
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		err := fmt.Errorf("login requires exactly 1 arguement: login <name>")
		return err
	}

	user_name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), user_name)
	if err != nil {
		return fmt.Errorf("%v user does not exist", user_name)
	}

	err = s.cfg.SetUser(user_name)
	if err != nil {
		return err
	}

	fmt.Printf("> User is set to: %v\n", user_name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register requires 1 arguement: register <name>")
	}

	user_name := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      user_name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("user with name %v already exists", user_name)
	}

	s.cfg.SetUser(user_name)
	fmt.Printf("*SYSTEM: USER %v was created successfully\n", user_name)
	fmt.Printf("*SYSTEM: USER: \n%v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset requires no arguements")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("*SYSTEM: SUCCESSFULLY RESET DATABASE")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users requires not arguements")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		str := fmt.Sprintf("* %v", user.Name)
		if s.cfg.CurrentUserName == user.Name {
			str = str + " (current)"
		}
		fmt.Println(str)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg expects one arguement: agg <time>\n i.e agg 1m30s")
	}

	time_between_reqs := cmd.args[0]
	timeBetweenReqs, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("could not parse time between requests: try something like 1m30s")
	}

	lowest_time_for_update := 5 * time.Second
	if timeBetweenReqs < lowest_time_for_update {
		fmt.Printf("> tried to set a time interval every %v, lowset is %v\n", timeBetweenReqs.Seconds(), lowest_time_for_update)
		timeBetweenReqs = lowest_time_for_update
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed expects one arguement: addfeed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		feed, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("could not find feed in DB, tried to create it but failed: url given %v", url)
		}
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("could not create feed follow")
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %v, url: %v, user: %v\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow expects one arguement: follow <url>")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not get feed for %v", url)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("%v already follows this source", user.Name)
	}

	fmt.Printf("User %v is now following %v\n", user.Name, feed.Url)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("was expecting a signle arguement: unfollow <url>")
	}

	url := cmd.args[0]

	// check if user is subscribed to the feed
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("user us not subscribed to: %v", url)
	}

	err = s.db.DeleteFeedFollow(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("something when trying unfollow %v", url)
	}

	fmt.Printf("> successfully unfollowed %v\n", feed.Url)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following does not expect any arguements")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feeds for user")
	}

	for i, feed := range feeds {
		fmt.Printf("%v. feed name: %v\n", i+1, feed.FeedName)
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var RSSFeed *RSSFeed

	err = xml.Unmarshal(body, &RSSFeed)
	if err != nil {
		return nil, fmt.Errorf("something went wrong when trying to parse the XML feed")
	}

	RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
	RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)

	for _, item := range RSSFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return RSSFeed, nil
}

func scrapeFeeds(s *state) error {
	num_of_feeds, err := s.db.GetTotalNumberOfFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get total number of feeds")
	}

	feeds_not_updated := 0
	for i := 0; i < int(num_of_feeds); i++ {
		next_feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return fmt.Errorf("could not fetch feed from db")
		}

		rssfeed, err := fetchFeed(context.Background(), next_feed.Url)
		if err != nil {
			feeds_not_updated += 1
			return fmt.Errorf("%v could not be updated", next_feed.Url)
		}

		feed, err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			ID: next_feed.ID,
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("could not update feed")
		}
		fmt.Printf("%v - %v was updated\n", feed.UpdatedAt, rssfeed.Channel.Title)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not get user")
		}

		return handler(s, cmd, current_user)
	}
}

func main() {
	var s state
	cfg := config.Read()
	s.cfg = &cfg

	cmds := commands{
		available: make(map[string]func(*state, command) error),
	}
	cmds.available["login"] = nil
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	args := os.Args

	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		log.Fatal("error while trying to connect to the Database", err)
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	if len(args) < 2 {
		log.Fatal("was expecting at least one arguement")
	}

	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
