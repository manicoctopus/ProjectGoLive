package mysql

import (
	"database/sql"
	"errors"

	"ProjectGoLive/pkg/models"
)

type ReviewModel struct {
	DB *sql.DB
}

func (m *ReviewModel) Create(c *models.Review) (int, error) {
	stmt := `INSERT INTO 
		review (ReviewText, UserID, ListID, Created) 
		VALUES (?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, c.ReviewText, c.UserID, c.ListID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ReviewModel) Update(c *models.Review) error {
	stmt := `UPDATE review 
			 SET ReviewText=?, UserID=?, ListID=?
			 WHERE ReviewID = ?`

	_, err := m.DB.Exec(stmt, c.ReviewText, c.UserID, c.ListID, c.ReviewID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ReviewModel) Delete(id uint32) error {
	stmt := `DELETE FROM review WHERE 
		ReviewID = ?`
	results, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	} else if rows, _ := results.RowsAffected(); rows == 0 {
		return models.ErrNoRecord
	} else {
		return err
	}
}

func (m *ReviewModel) Retrieve(id uint32) (*models.Review, error) {
	stmt :=
		`SELECT 
			ReviewID, ReviewText, UserID, ListID 
		FROM review
		WHERE ReviewID = ?`
	r := &models.Review{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&r.ReviewID, &r.ReviewText, &r.UserID, &r.ListID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return r, nil
}

func (m *ReviewModel) RetrieveAll() ([]*models.Review, error) {
	stmt :=
		`SELECT 
			ReviewID, ReviewText, UserID, ListID 
		FROM review`

	reviews := []*models.Review{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		r := &models.Review{}
		err := rows.Scan(
			&r.ReviewID, &r.ReviewText, &r.UserID, &r.ListID,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}
