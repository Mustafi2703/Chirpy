# Chirpy

A Twitter-like social media platform that allows users to share short messages and interact with other users.

## Overview

Chirpy is a lightweight social media application that enables users to post short messages ("chirps"), follow other users, and engage with content through likes and replies.

## Features

- User authentication (sign up, login, logout)
- Create, read, update, and delete chirps
- Follow/unfollow other users
- Like and reply to chirps
- User profiles with bio and avatar
- Timeline of chirps from followed users

## Installation

### Prerequisites
- Go (version 1.22.5)
- PostgreSQL database

### Setup
1. Clone the repository from https://github.com/Mustafi2703/Chirpy2.git
2. Install dependencies
   go mod download

3. Set up environment variables
   cp .env.example .env

4. Run the application
   go run main.go


## Usage

Once running, access the application at `http://localhost:8080`

- Create an account or login
- Post new chirps from your dashboard
- Browse the public timeline or your personalized feed
- Click on usernames to view profiles

## Technologies Used

- Go
- HTML/CSS/JavaScript
- SQLite/PostgreSQL
- [Any other technologies you used]

## Future Enhancements

- Direct messaging between users
- Hashtag support
- Media uploads
- Mobile application

## License

This is a personal portfolio project created for educational purposes.
