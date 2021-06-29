package mysql

import (
	"database/sql"
	"errors"

	"ProjectGoLive/pkg/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Create(c *models.Category) (int, error) {
	stmt := `INSERT INTO 
		category (CatName, ParentCat) 
		VALUES (?, ?)`

	result, err := m.DB.Exec(stmt, c.CatName, c.ParentCat)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *CategoryModel) Update(c *models.Category) error {
	stmt := `UPDATE category 
			 SET CatName=?, ParentCat=?
			 WHERE CatID = ?`

	_, err := m.DB.Exec(stmt, c.CatName, c.ParentCat, c.CatID)
	if err != nil {
		return err
	}

	return nil
}

func (m *CategoryModel) Delete(id int) error {
	stmt := `DELETE FROM category WHERE 
		CatID = ?`
	results, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	} else if rows, _ := results.RowsAffected(); rows == 0 {
		return models.ErrNoRecord
	} else {
		return err
	}
}

func (m *CategoryModel) Retrieve(id int) (*models.Category, error) {
	stmt :=
		`SELECT 
			CatID, CatName, ParentCat
		FROM category
		WHERE CatID = ?`
	c := &models.Category{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&c.CatID, &c.CatName, &c.ParentCat,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return c, nil
}

func (m *CategoryModel) RetrieveAll() ([]*models.Category, error) {
	stmt :=
		`SELECT 
			CatID, CatName, ParentCat
		FROM category`

	categories := []*models.Category{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		c := &models.Category{}
		err := rows.Scan(
			&c.CatID, &c.CatName, &c.ParentCat,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
