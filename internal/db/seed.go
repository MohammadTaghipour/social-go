package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/MohammadTaghipour/social/internal/store"
)

var usernames = []string{
	"tiago", "bob", "jack", "alice", "maria",
	"kevin", "lucas", "sara", "david", "julia",
	"max", "oliver", "emma", "liam", "sophia",
	"noah", "ava", "ethan", "mia", "daniel",
	"zoe", "chris", "isabella", "mason", "ella",
	"leo", "grace", "ben", "chloe", "nathan",
	"ryan", "lily", "sebastian", "nora", "matt",
	"alex", "hannah", "sam", "victor", "elena",
	"george", "carla", "henry", "amelia", "peter",
	"felix", "clara", "thomas", "sofia", "marco",
}

var titles = []string{
	"Learning Go Basics", "Concurrency in Go", "Understanding Interfaces",
	"REST API with Go", "Error Handling in Go", "Building Microservices",
}

var contents = []string{
	"Go is a statically typed, compiled language designed for simplicity and performance.",
	"Concurrency in Go is powerful thanks to goroutines and channels.",
	"This post explains how to build REST APIs using Go's net/http package.",
	"Error handling in Go uses explicit error returns rather than exceptions.",
	"Go interfaces provide a way to define behavior without inheritance.",
}

var allTags = []string{"go", "programming", "backend", "http", "api", "database", "testing", "cloud"}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "email.com",
			RoleID:   1,
		}
		users[i].Password.Set("123123")
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[i%len(titles)],
			Content: contents[i%len(contents)],
			Tags: []string{
				allTags[i%len(allTags)],
				allTags[i%len(allTags)],
			},
		}
	}

	return posts
}

func generateComments(num int, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  post.UserID,
			Content: post.Title + post.Content,
		}
	}

	return comments
}
