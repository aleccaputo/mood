package db

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "golang.org/x/crypto/bcrypt"
    "mood/models"
    "os"
)

type PostgresDb struct {
    db *pgxpool.Pool
}

func NewPostgresDb() (*PostgresDb, error) {
    return Connect()
}

func Connect() (*PostgresDb, error) {
    dbUrl := os.Getenv("DATABASE_URL")
    dbpool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
        return nil, err
    }
    // defer dbpool.Close()

    return &PostgresDb{
        db: dbpool,
    }, nil
}

func (s *PostgresDb) Init() error {
    return s.CreateUserTable()
}

func (s *PostgresDb) CreateUserTable() error {
    createTableStatement := `
	CREATE TABLE IF NOT EXISTS users (
		id uuid PRIMARY KEY,
		first_name text NOT NULL,
		last_name text NOT NULL,
		email text NOT NULL,
		created_at timestamp NOT NULL,
		hashed_password text NOT NULL,
		UNIQUE (email)
	)`

    _, err := s.db.Exec(context.Background(), createTableStatement)
    return err
}

func (s *PostgresDb) Get(id uuid.UUID) (*models.User, error) {
    userByIdQuery := `SELECT (id, first_name, last_name, email, created_at) FROM users WHERE id = $1`
    user := &models.User{}
    stringUuid := id.String()
    err := s.db.QueryRow(context.Background(), userByIdQuery, stringUuid).Scan(&user)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting user: %v\n", err)
        return nil, err
    }

    return user, nil
}

func (s *PostgresDb) CreateUser(user *models.User, password string) (*uuid.UUID, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    var userId uuid.UUID
    insertStatement := `insert into users (first_name, last_name, email, created_at, id, hashed_password) values ($1, $2, $3, $4, $5, $6) returning id`
    dbErr := s.db.QueryRow(context.Background(), insertStatement, user.FirstName, user.LastName, user.Email, user.CreatedAt, user.Id, hashedPassword).Scan(&userId)
    if dbErr != nil {
        return nil, dbErr
    }

    return &userId, nil
}

func (s *PostgresDb) LoginUser(username, password string) (*uuid.UUID, error) {
    var hashedPassword string
    var userId *uuid.UUID
    getUserByEmailQuery := `SELECT hashed_password, id FROM users WHERE email = $1`
    err := s.db.QueryRow(context.Background(), getUserByEmailQuery, username).Scan(&hashedPassword, &userId)
    if err != nil {
        return nil, errors.New("User not found")
    }

    result := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if result == nil {
        return userId, nil
    }

    return nil, errors.New("invalid")
}
