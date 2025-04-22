# SyncSpend Backend

**SyncSpend** is a backend service designed to manage and synchronize spending data across various platforms. Built with **Go**, it follows a modular structure for scalability and maintainability.

## Features

- User authentication and management
- Expense tracking and synchronization
- RESTful API design
- Middleware support
- Modular architecture (services, repositories, controllers)

## Tech Stack

- **Go**: Main programming language
- **Gorilla Mux**: HTTP router
- **GORM**: ORM for database interaction
- **JWT**: Authentication with JSON Web Tokens

## Project Structure

    ├── config          # Configuration files
    ├── controllers     # Handlers for HTTP requests
    ├── dto             # Data Transfer Objects for requests/responses
    ├── helpers         # Utility functions
    ├── middleware      # Middleware for authentication, logging, etc.
    ├── models          # Data models for GORM
    ├── repositories    # Database interaction logic
    ├── routes          # API route definitions
    ├── services        # Business logic
    ├── main.go         # Application entry point

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/daudoscar/syncspend-backend.git
   ```

2. Navigate into the project directory:
   ```bash
   cd syncspend-backend
   ```

3. Install Dependencies:
   ```bash
   go mod download
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## Credits

Developed and maintained by [Oscar Daud](https://github.com/daudoscar).
