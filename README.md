# Mood

Mood is an application designed to help you log your daily mood and analyze the data to find the factors that make your days the best. The app is built with Go for the backend and Remix for the frontend. This project is structured into two main folders: `server` (containing the Go code) and `client` (containing the Remix code).

## Features

- Log daily mood with overall ratings and specific mood descriptors.
- Analyze mood data to identify patterns and factors that influence your mood.
- User-friendly interface for entering and viewing mood data.

## Project Structure

- **`server`**: Contains the Go code for the backend API.
- **`client`**: Contains the Remix code for the frontend application.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose
- Node.js 14 or higher
- PostgreSQL (if running outside of Docker)

### Installation

Follow these steps to set up the project:

1. **Clone the repository**

2. **Set up the backend with Docker:**

   - Navigate to the `server` directory:

     ```bash
     cd server
     ```

   - Build and start the Docker containers:

     ```bash
     docker-compose up --build
     ```

   - This will set up the Go server and PostgreSQL database. Ensure the connection settings in your `docker-compose.yml` file are configured correctly.

3. **Set up the frontend:**

   - Navigate to the `client` directory:

     ```bash
     cd client
     ```

   - Install Node.js dependencies:

     ```bash
     npm install
     ```

   - Run the Remix development server:

     ```bash
     npm run dev
     ```

4. **Access the application:**

   Open your web browser and navigate to `http://localhost:3000` to use the Mood app.

## Contributing

We welcome contributions! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Commit your changes.
4. Push your branch to your fork.
5. Create a pull request with a description of your changes.
