package cafe

import (
	"database/sql"
	"errors"
	"time"
)

type CafeDomain struct {
	db *sql.DB
}

func NewCafeDomain(db *sql.DB) (*CafeDomain, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}

	return &CafeDomain{db: db}, nil
}

type Cafe struct {
	ID                     int
	Name                   string
	GmapsLink              string
	CreatedByDiscordUserID string
	CreatedAt              string
	UpdatedAt              string
}

type ListParam struct {
	Limit  int
	Offset int
}

func (c *CafeDomain) List(param ListParam) ([]Cafe, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	query := `SELECT id, name, gmaps_link, created_by_discord_user_id, created_at, updated_at FROM cafe LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, param.Limit, param.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cafes []Cafe
	for rows.Next() {
		var cafe Cafe
		err = rows.Scan(&cafe.ID, &cafe.Name, &cafe.GmapsLink, &cafe.CreatedByDiscordUserID, &cafe.CreatedAt, &cafe.UpdatedAt)
		if err != nil {
			return nil, err
		}

		cafes = append(cafes, cafe)
	}

	return cafes, nil
}

func (c *CafeDomain) validateInsert(cafe *Cafe) error {
	if cafe.Name == "" {
		return errors.New("name is required")
	}

	if cafe.GmapsLink == "" {
		return errors.New("gmaps link is required")
	}

	if cafe.CreatedByDiscordUserID == "" {
		return errors.New("discord user id is required")
	}

	return nil
}

func (c *CafeDomain) Insert(cafe *Cafe) error {
	err := c.validateInsert(cafe)
	if err != nil {
		return err
	}

	createdAt := time.Now().String()
	query := `INSERT INTO cafe (name, gmaps_link, created_by_discord_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	res, err := c.db.Exec(query, cafe.Name, cafe.GmapsLink, cafe.CreatedByDiscordUserID, createdAt, createdAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	cafe.ID = int(id)
	cafe.CreatedAt = createdAt
	cafe.UpdatedAt = createdAt

	return nil
}
