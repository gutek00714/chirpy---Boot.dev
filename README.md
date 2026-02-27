# Chirpy

Chirpy is a custom social media backend server built with Go. This project was developed to practice building robust RESTful APIs, handling user authentication, and managing persistent data with PostgreSQL.

## Features

- **User Management**: Create accounts and manage user profiles.
- **Authentication**: Secure login using JWT (JSON Web Tokens) and Refresh Tokens.
- **Chirps**: Users can create "chirps" (posts) that are stored in a database.
- **Content Filtering**: Automatic filtering of profanity from user posts.
- **Admin Tools**: A dedicated admin metrics page to track server health and usage.
- **Database Integration**: Full persistence using PostgreSQL and the `sqlc` library for type-safe SQL queries.

## Prerequisites

To run this project, you will need:

- Go 1.22 or higher
- PostgreSQL
- A `.env` file containing your database connection string and JWT secret
