<h1 style="text-align: center;">Atur Dana BE</h1>
<br />

**Atur Dana** is a personal financial management app designed to help you manage your personal finances effectively. This repository serves as the backend for **Atur Dana**, written in [Golang](https://go.dev/).

## Features

- Track income and expenses
- Budgeting tools
- Financial reports
- Secure user authentication
- Real-time data synchronization

## Technologies Used

- **Programming Language:** Golang
- **Database:** PostgreSQL
- **API:** RESTful API
- **Authentication:** JWT
- **Frameworks and Libraries:** GORM

## Getting Started

### Prerequisites

- Go (version 1.20.1 or higher)
- PostgreSQL
- Git

## Project Structure
```bash
├── build               # CI/CD configuration files
├── cmd                 # Main application entry point
│   └── main.go          # Main file
└── internal            # Internal packages (shhh... it's internal 🤫)
    ├── auth            # Authentication logic
    ├── common          # Common utilities and helpers
    ├── db              # Database interactions and migrations
    ├── handlers        # HTTP request handlers
    ├── middleware      # Middleware components
    ├── models          # Database entities
    ├── requests        # Request payload definitions
    ├── responses       # Response payload definitions
    └── routes          # API route definitions

```

### Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/paulmhrdka/aturdana-backend.git
    cd aturdana-backend
    ```

2. **Set up environment variables:**

    copy `.env.example` file into `.env` file, and insert your credential for DB connection & JWT Secret

    ```sh
    cp .env.example .env
    ```

3. **Install dependencies:**

    ```sh
    go mod tidy
    ```

4. **Start the server:**

    ```sh
    go run cmd/main.go # it automatically run migration
    ```

## Contact

For any questions or feedback, please reach out to:

- **Email:** [mahardikapaul@gmail.com](mailto:mahardikapaul@gmail.com)
- **GitHub Issues:** [GitHub Issues Page]
