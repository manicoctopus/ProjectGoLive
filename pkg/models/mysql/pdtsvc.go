package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ProjectGoLive/pkg/models"
)

type PdtsvcModel struct {
	DB *sql.DB
}

func (m *PdtsvcModel) Create(c *models.Pdtsvc) (int, error) {
	stmt := `INSERT INTO pdtsvc 
			 (PdtsvcName, PdtsvcPrice, PdtsvcDesc, CatID, ListID, Views, Likes, Keyword, Created, Modified) 
    		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, c.PdtsvcName, c.PdtsvcPrice, c.PdtsvcDesc, c.CatID,
		c.ListID, c.Views, c.Likes, c.Keyword)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PdtsvcModel) Update(c *models.Pdtsvc) error {
	stmt := `UPDATE pdtsvc 
			 SET PdtsvcName=?, PdtsvcPrice=?, PdtsvcDesc=?, CatID=?, ListID=?, Views=?, Likes=?, Keyword=?, Modified=UTC_TIMESTAMP() 
			 WHERE PdtsvcID = ?`

	_, err := m.DB.Exec(stmt, c.PdtsvcName, c.PdtsvcPrice, c.PdtsvcDesc, c.CatID, c.ListID, c.Views, c.Likes, c.Keyword, c.PdtsvcID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *PdtsvcModel) Delete(id uint32) error {
	stmt := `DELETE FROM pdtsvc WHERE PdtsvcID=?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (m *PdtsvcModel) Retrieve(id uint32) (*models.Pdtsvc, error) {
	stmt :=
		`SELECT 
			PdtsvcID, PdtsvcName, PdtsvcPrice, PdtsvcDesc, CatID, ListID, 
			Views, Likes, Keyword, Created, Modified 
		FROM pdtsvc
		WHERE PdtsvcID = ?`

	p := &models.Pdtsvc{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&p.PdtsvcID, &p.PdtsvcName, &p.PdtsvcPrice, &p.PdtsvcDesc,
		&p.CatID, &p.ListID, &p.Views, &p.Likes, &p.Keyword,
		&p.Created, &p.Modified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PdtsvcModel) RetrieveAll() ([]*models.Pdtsvc, error) {
	stmt :=
		`SELECT 
			PdtsvcID, PdtsvcName, PdtsvcPrice, PdtsvcDesc, CatID, ListID, 
			Views, Likes, Keyword, Created, Modified 
		FROM pdtsvc`

	products := []*models.Pdtsvc{}

	rows, _ := m.DB.Query(stmt)

	for rows.Next() {
		p := &models.Pdtsvc{}
		err := rows.Scan(
			&p.PdtsvcID, &p.PdtsvcName, &p.PdtsvcPrice, &p.PdtsvcDesc,
			&p.CatID, &p.ListID, &p.Views, &p.Likes, &p.Keyword,
			&p.Created, &p.Modified,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (m *PdtsvcModel) GetSellerPdtsvcs(sellerID string) ([]*models.Pdtsvc, error) {
	stmt := `SELECT 
				 PdtsvcID, Name, Description, Price, CategoryID, 
				 Inventory, Created, UserID, Rating, RatingNum, UnitSold
			FROM pdtsvc
			WHERE UserID=?`

	rows, err := m.DB.Query(stmt, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Pdtsvc{}

	for rows.Next() {
		p := &models.Pdtsvc{}
		err := rows.Scan(
			&p.PdtsvcID, &p.PdtsvcName, &p.PdtsvcPrice, &p.PdtsvcDesc,
			&p.CatID, &p.ListID, &p.Views, &p.Likes, &p.Keyword,
			&p.Created, &p.Modified,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (m *PdtsvcModel) GetSearchPdtsvcs() ([]*models.Pdtsvc, error) {
	results, err := m.DB.Query("Select PdtsvcID, PdtsvcName, PdtsvcDesc, Keyword FROM Pdtsvc")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := []*models.Pdtsvc{}
	for results.Next() {
		p := &models.Pdtsvc{}
		err = results.Scan(&p.PdtsvcID, &p.PdtsvcName, &p.PdtsvcDesc, &p.Keyword)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, p)

	}
	return products, nil
}

func (m *PdtsvcModel) GetSearchResults(rankedIndex []int) ([]*models.Pdtsvc, error) {
	if len(rankedIndex) == 0 {
		return nil, nil
	}

	var pIDs string
	for _, ID := range rankedIndex {
		i := strconv.Itoa(ID)
		pIDs = pIDs + "," + i
	}
	pIDs = "(" + strings.TrimLeft(pIDs, ",") + ")"
	stmt := `SELECT PdtsvcID, Name, Description, Price, CategoryID, Inventory, Created, SellerID, Rating, RatingNum, UnitSold
	FROM Pdtsvc WHERE PdtsvcID IN ` + pIDs

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Pdtsvc{}
	for rows.Next() {
		p := &models.Pdtsvc{}
		err = rows.Scan(
			&p.PdtsvcID,
			&p.PdtsvcName,
			&p.PdtsvcPrice,
			&p.PdtsvcDesc,
			&p.CatID,
			&p.ListID,
			&p.Created,
			&p.Views,
			&p.Likes,
			&p.Keyword,
			&p.Created,
			&p.Modified,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
