package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"social/internal/store"
)

var userNames = []string{
	"Alice", "Bob", "Charlie", "David", "Eve",
	"Frank", "Grace", "Heidi", "Ivan", "Judy",
	"Kathy", "Leo", "Mallory", "Niaj", "Olivia",
	"Peggy", "Quentin", "Rupert", "Sybil", "Trent",
	"Uma", "Victor", "Wendy", "Xander", "Yvonne", "Zack",
}

var titles = []string{
	"Exploring the Go Programming Language",
	"Understanding RESTful APIs",
	"Introduction to Microservices Architecture",
	"Building Scalable Web Applications",
	"Database Design Best Practices",
	"Getting Started with Docker",
	"An Overview of Kubernetes",
	"Implementing Authentication and Authorization",
	"Optimizing Web Application Performance",
	"Deploying Applications to the Cloud",
	"Building Microservices with Go",
	"Mastering Docker for Containers",
	"Building RESTful APIs with Go",
	"Managing Kubernetes Clusters",
	"Securing Microservices with JWT",
	"Optimizing Database Performance",
	"Deploying Microservices to the Cloud",
	"Building Microservices with Go",
}

var contents = []string{
	"This is a sample post content about Go programming.",
	"Learn how to build RESTful APIs in this post.",
	"An introduction to microservices architecture.",
	"Tips for building scalable web applications.",
	"Best practices for database design.",
	"Getting started with Docker containers.",
	"An overview of Kubernetes.",
	"Implementing authentication and authorization.",
	"Optimizing web application performance.",
	"Deploying applications to the cloud.",
	"Building microservices with Go.",
	"Mastering Docker for containers.",
	"Building RESTful APIs with Go.",
	"Managing Kubernetes clusters.",
	"Securing microservices with JWT.",
	"Optimizing database performance.",
	"Deploying microservices to the cloud.",
	"Building microservices with Go.",
}

var tags = []string{
	"golang", "api", "microservices", "web", "database",
	"docker", "kubernetes", "auth", "performance", "cloud",
	"jwt", "scalability", "design", "development", "deployment",
}

var comments = []string{
	"Great post! Very informative.",
	"Thanks for sharing this.",
	"I found this post very helpful.",
	"Can you provide more details on this topic?",
	"This is exactly what I was looking for.",
	"Interesting perspective!",
	"I disagree with some points made here.",
	"Looking forward to more posts like this.",
	"Could you share some code examples?",
	"This helped me understand the concept better.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		err := store.Users.Create(ctx, user)
		if err != nil {
			log.Println("Error seeding user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Println("Error seeding post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Println("Error seeding comment:", err)
			return
		}
	}

	log.Println("Database seeding completed successfully.")


}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)

	for i := range users {
		userName := userNames[i%len(userNames)] + fmt.Sprintf("%d", i)
		users[i] = &store.User{
			Username: userName,
			Email:    userName + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)

	for i := range posts {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			UserID:  user.ID,
			Tags:    []string{tags[rand.Intn(len(tags))], tags[rand.Intn(len(tags))]},
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	commentsList := make([]*store.Comment, n)

	for i := range commentsList {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		commentsList[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return commentsList
}
