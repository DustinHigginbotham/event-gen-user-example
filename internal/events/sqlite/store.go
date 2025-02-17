package sqlite

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver

	"context"
	"database/sql"

	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

type Store struct {
	db       *sql.DB
	registry *app.Registry
}

func (s *Store) Migrate() error {

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("could not get migration data: %w", err)
	}

	db, err := sqlite.WithInstance(s.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver instance for migration: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", db)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

func (s *Store) Connect(location string) error {
	db, err := sql.Open("sqlite3", location)
	if err != nil {
		return err
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil

}

// Get implements app.Store.
func (s *Store) Get(ctx context.Context, aggregateID string) ([]app.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		WITH latest_snapshot AS (
			SELECT version
			FROM events
			WHERE aggregate_id = ? AND is_snapshot = 1
			ORDER BY version DESC
			LIMIT 1
		)
		SELECT id, type, created, metadata, payload, version
		FROM events
		WHERE aggregate_id = ?
		AND (version >= (SELECT version FROM latest_snapshot)
			OR NOT EXISTS (SELECT 1 FROM latest_snapshot))
		ORDER BY version ASC`, aggregateID, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventList []app.Event
	for rows.Next() {
		var e app.Event
		var payload []byte
		var metadata any
		if err := rows.Scan(&e.ID, &e.Type, &e.Created, &metadata, &payload, &e.Version); err != nil {
			return nil, err
		}

		event, err := s.registry.Deserialize(e.Type, payload)
		if err != nil {
			return nil, err
		}

		e.Data = event

		// Unmarshal the payload into the correct event type
		// (use a registry or type switch if necessary)
		e.Metadata = make(map[string]interface{})
		eventList = append(eventList, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return eventList, nil
}

// Save implements app.Store.
func (s *Store) Save(ctx context.Context, aggregateID string, event app.Eventer) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	user, err := app.GetUserFromContext[app.AuthUser](ctx)
	if err != nil {
		fmt.Println("Failed to get user auth", err)
	}

	var count int
	row := s.db.QueryRow(`SELECT COUNT(aggregate_id) FROM events WHERE aggregate_id = ?`, aggregateID)
	if err := row.Scan(&count); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO events (aggregate_id, id, user_id, user_name, type, created, metadata, payload, version, is_snapshot)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	payload, err := json.Marshal(event) // Serialize event payload
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		aggregateID,
		uuid.New().String(),
		user.ID,
		user.Name,
		event.GetType(),
		time.Now().UTC(),
		nil,
		payload,
		count+1,
		false,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Store) Close() error {
	return s.db.Close()
}

var _ app.Store = new(Store)

func New(a *app.App) *Store {
	return &Store{
		registry: a.Registry(),
	}
}
