
# Test Signer Service

The Test Signer is a service written in Go that accepts a set of answers and questions and signs that the user has finished the "test" at a specific point in time. The signatures are stored and can later be verified by a different service.

## Requirements

- Go 1.21.5
- Docker
- HTTP client (e.g., Postman)

## Getting Started

Follow these steps to set up and run the Test Signer service:

## Download Options

You have two options to download and start using the project:

### 1. Unzip the ZIP File

- After downloading the project, unzip the .zip file into your preferred directory.

### 2. Clone from GitHub

- You can also clone the repository directly from GitHub using the following command in your terminal or command line:

  ```bash
  git clone https://github.com/PelicanSL/Signer.git

### 2. Navigate to the Project Directory

- Open your terminal or command prompt and navigate to the project directory:

```bash
cd Signer/
```


### 3. Download Go Dependencies

- In the project directory, run the following command to download the Go dependencies listed in the go.mod file:

```bash
go mod download
```

This will ensure that all required packages and dependencies are installed.

### 4. Build and Run with Docker

- Once you have downloaded the dependencies, you can proceed to build and run the Test Signer service using Docker. Make sure you have Docker installed and running on your system.

```bash
docker-compose up --build
```

This command will start the Test Signer service along with any required dependencies in Docker containers.

### 5. Stopping and Cleaning Up

- If you need to stop the service and remove previous databases:

```bash
docker-compose down -v
```

This command will stop the service and remove the associated Docker volumes, clearing previous database data.

## Database Structure

The Test Signer service utilizes a PostgreSQL database with the following table structure:

### `users`

- `email` (Primary Key): User's email address.
- `password`: User's password.

### `questions`

- `id` (Primary Key): Unique identifier for each question.
- `text`: The text of the question.

### `answers`

- `id` (Primary Key): Unique identifier for each answer.
- `question_id` (Foreign Key): References the `questions` table, indicating the question to which the answer belongs.
- `text`: The text of the answer.
- `user_email` (Foreign Key): References the `users` table, indicating the user who provided the answer. Each user can provide answers to multiple questions. There is a unique constraint on `(question_id, user_email)` to ensure that a user can provide only one answer per question.

### `signatures`

- `id` (Primary Key): Unique identifier for each signature.
- `user_email` (Foreign Key): References the `users` table, indicating the user associated with the signature.
- `timestamp`: The timestamp when the signature was created. It has a default value of the current timestamp.
- `sign`: The signature data.

## Initializing the Database

To set up the database, the project includes an `init.sql` file that is executed when the Docker container is started. This SQL script performs the following actions:

- Creates the `users`, `questions`, `answers`, and `signatures` tables if they do not already exist.
- Inserts 10 sample questions into the `questions` table for testing purposes.

You do not need to manually create the tables or insert the sample questions; the Docker container handles this initialization process for you.

Feel free to explore and use these tables when interacting with the Test Signer service.

## Configuration

The Test Signer service uses configuration values stored in a `.env` file for security reasons. You do not need to create these values manually; Docker will load them automatically when you run the service. Below are the configuration variables used:

### JWT Configuration

- `JWT_SECRET_KEY`: Secret key used for JWT token generation and verification. Example value: `superpasswordforjwt`

### Database Configuration

- `DB_HOST`: Hostname of the PostgreSQL database. Example value: `db`
- `DB_PORT`: Port number of the PostgreSQL database. Example value: `5432`
- `DB_USER`: Username for database authentication. Example value: `user`
- `DB_PASSWORD`: Password for database authentication. Example value: `password`
- `DB_NAME`: Name of the PostgreSQL database. Example value: `mydb`

These configuration values are stored securely in the `.env` file and are automatically loaded when you start the Docker container. There is no need to manually create or configure these values.


## API Endpoints

App is running in the default port 8080. Normally, requests will have the next structure: http://localhost:8080/ENDPOINT

The Test Signer service provides the following API endpoints for interaction:

### 1. Register User

- **Endpoint**: `POST http://localhost:8080/register`
- **Description**: This endpoint allows users to register and create a new account.

**Request body:**
```json
{
  "email": "user@example.com",
  "password": "P@ssw0rd!"
}
```

**Requirements:**

- The "email" field must have a valid email structure (e.g., "user@example.com").
- The password must meet the following criteria:
  - It must be at least 8 characters long.
  - It must include at least one uppercase letter.
  - It must contain at least one symbol (special character).
  - It must have at least one number.
  - It must not exceed 100 characters in length.
- If an attempt is made to register an email that is already registered in the system, an error will be returned.
- If all requirements are met, and user creation is successful, the endpoint will return a success message along with a JSON Web Token (JWT).
- The generated JWT will have a validity period of 1 hour, allowing the user to use it for authentication within that timeframe.


**Response:**
```json
{
    "message": "User registered successfully",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDI3NDQ2ODksInN1YiI6InBlZHJvQGdtYWlsLmNvbSJ9.AQh3qLoHYEjcRLJqPOC967MuGhRv7AmjdKAWFG7f_bc"
}
```

### 2. Sign Answers

- **Endpoint**: `POST http://localhost:8080/sign_answers`
- **Description**: This endpoint allows users to sign their answers to a test.
**Authorization**: Select "Bearer Token" and add the token received during user registration, in this case: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDI3NDQ2ODksInN1YiI6InBlZHJvQGdtYWlsLmNvbSJ9.AQh3qLoHYEjcRLJqPOC967MuGhRv7AmjdKAWFG7f_bc


**Request Body**:

```json
{
  "answers": [
    {
      "question_id": 1,
      "text": "Answer to question 1"
    },
    {
      "question_id": 2,
      "text": "Answer to question 2"
    },
    {
      "question_id": 3,
      "text": "Answer to question 3"
    }
  ]
}

```

**Requirements:**

- The user must provide a valid JSON Web Token (JWT) in the Authorization
- Only questions with valid IDs from the database (up to question number 10) can be answered.
- If a question is answered multiple times with the same ID, only the last answer will be considered.
- If a new set of questions is signed, the previous answers will be replaced in the database.

**Response:**
```json
{
    "signature": "124f2e9ab7427e818f460f9149848064ba3d9e273f8558ef536e5abce098f54c"
}
```

### 3. Verify Signature

- **Endpoint**: `GET http://localhost:8080/verify_signature`
- **Description**: This endpoint allows users to verify the signature of their previously signed answers. 

The user_email param is the one that we signed up in Step 1. 

The signature param is the one generated from Step 2.

**Authorization**: Not needed!

**Request Body**:

```json
{
  "user_email": "user@example.com",
  "signature": "124f2e9ab7427e818f460f9149848064ba3d9e273f8558ef536e5abce098f54c"
}
```

**Requirements:**

- The "user_email" specified in the request body must be registered in the database.
- The "signature" provided must correspond to the specified user.
- It is not required for the JWT to be currently valid or included in the authorization header, as specified in the exercise instructions.

**Response:**
```json
{
    "answers": [
        {
            "ID": 1,
            "QuestionID": 1,
            "Text": "Answer to question 1"
        },
        {
            "ID": 2,
            "QuestionID": 2,
            "Text": "Answer to question 2"
        },
        {
            "ID": 3,
            "QuestionID": 3,
            "Text": "Answer to question 3"
        }
    ],
    "status": "OK",
    "timestamp": "2023-12-16T16:03:19.390091Z"
}
```


