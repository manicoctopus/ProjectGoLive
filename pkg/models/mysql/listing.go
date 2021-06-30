package mysql

import (
	"database/sql"
	"errors"

	"ProjectGoLive/pkg/models"
)

type ListingModel struct {
	DB *sql.DB
}

func (m *ListingModel) Create(c *models.Listing) (int, error) {
	stmt := `INSERT INTO 
		listing (ListName, ListDesc, Ig_url, Fb_url, Website_url, UserID, Created, Modified) 
		VALUES (?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, c.ListName, c.ListDesc, c.Ig_url, c.Fb_url, c.Website_url, c.UserID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ListingModel) Update(c *models.Listing) error {
	stmt := `UPDATE listing 
			 SET ListName=?, ListDesc=?, Ig_url=?, Fb_url=?, Website_url=?, UserID=?, Modified=UTC_TIMESTAMP() 
			 WHERE ListID = ?`

	_, err := m.DB.Exec(stmt, c.ListName, c.ListDesc, c.Ig_url, c.Fb_url, c.Website_url, c.UserID, c.ListID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ListingModel) Delete(id uint32) error {
	stmt := `DELETE FROM listing WHERE 
		ListID = ?`
	results, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	} else if rows, _ := results.RowsAffected(); rows == 0 {
		return models.ErrNoRecord
	} else {
		return err
	}
}

func (m *ListingModel) Retrieve(id uint32) (*models.Listing, error) {
	stmt :=
		`SELECT 
			ListID, ListName, ListDesc, Ig_url, Fb_url, Website_url, 
			UserID,  Created, Modified 
		FROM listing
		WHERE ListID = ?`
	l := &models.Listing{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&l.ListID, &l.ListName, &l.ListDesc, &l.Ig_url,
		&l.Fb_url, &l.Website_url, &l.UserID,
		&l.Created, &l.Modified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return l, nil
}

func (m *ListingModel) RetrieveAll() ([]*models.Listing, error) {
	stmt :=
		`SELECT 
			ListID, ListName, ListDesc, Ig_url, Fb_url, Website_url, 
			UserID,  Created, Modified 
		FROM listing`

	listings := []*models.Listing{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		l := &models.Listing{}
		err := rows.Scan(
			&l.ListID, &l.ListName, &l.ListDesc, &l.Ig_url,
			&l.Fb_url, &l.Website_url, &l.UserID,
			&l.Created, &l.Modified,
		)
		if err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}
	return listings, nil
}
