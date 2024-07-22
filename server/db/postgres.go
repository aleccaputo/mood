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

//  I don't know why the Go community hates ORMs but they seem to, so i guess i'm writing raw queries. Good practice i guess

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
	err := s.CreateUserTable()
	if err != nil {
		return err
	}
	return s.CreateEntryTable()
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

func (s *PostgresDb) CreateEntryTable() error {
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS entries (
		id uuid PRIMARY KEY,
		overall integer NOT NULL,
		descriptors integer[] NOT NULL,
		good_notes text NOT NULL,
		bad_notes text NOT NULL,
		exercise boolean NOT NULL,
		alcohol boolean NOT NULL,
		created_at timestamp NOT NULL,
		user_id uuid NOT NULL,
	    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	)`

	_, err := s.db.Exec(context.Background(), createTableStatement)
	return err
}

func (s *PostgresDb) GetUserById(id uuid.UUID) (*models.User, error) {
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

func (s *PostgresDb) InsertUser(user *models.User, password string) (*uuid.UUID, error) {
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

func (s *PostgresDb) InsertEntry(entry *models.Entry) (*uuid.UUID, error) {
	var noteId uuid.UUID
	insertStatement := `insert into entries (id, overall, descriptors, good_notes, bad_notes, exercise, alcohol, created_at, user_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	dbErr := s.db.QueryRow(context.Background(), insertStatement, entry.Id, entry.Overall, entry.Descriptors, entry.GoodNotes, entry.BadNotes, entry.Exercise, entry.Alcohol, entry.CreatedAt, entry.UserId).Scan(&noteId)
	if dbErr != nil {
		return nil, dbErr
	}

	return &noteId, nil
}

func (s *PostgresDb) GetEntriesByUserId(id uuid.UUID) ([]models.Entry, error) {
	entriesByIdQuery := `SELECT * FROM entries WHERE user_id = $1`
	stringUuid := id.String()
	rows, err := s.db.Query(context.Background(), entriesByIdQuery, stringUuid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var entries []models.Entry

	for rows.Next() {
		var entry models.Entry
		var descriptors []int
		if err := rows.Scan(&entry.Id, &entry.Overall, &descriptors, &entry.GoodNotes, &entry.BadNotes, &entry.Exercise, &entry.Alcohol, &entry.CreatedAt, &entry.UserId); err != nil {
			return nil, err
		}

		// Convert []int to []MoodDescriptor
		entry.Descriptors = make([]models.MoodDescriptor, len(descriptors))
		for i, d := range descriptors {
			entry.Descriptors[i] = models.MoodDescriptor(d)
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
