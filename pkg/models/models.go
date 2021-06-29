package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//var Category = []string{"Frozen Food", "Staples", "Meat and Seafood", "Beverages", "Fruit and Vegetables"}
//var SortBy = []string{"Popular", "Highly Rated", "Price (Asc.)", "Price (Desc.)"}

var MapOTP = make(map[string]string)

type User struct {
	UserID         uint32    `json:"userID"`
	UserName       string    `json:"userName"`
	UserEmail      string    `json:"userEmail"`
	HashedPassword []byte    `json:"hashedPassword"`
	UserContact    string    `json:"userContact"`
	IsBOwner       uint8     `json:"isBOwner"`
	IsVerified     uint8     `json:"isVerified"`
	Created        time.Time `json:"created"`
}

type Category struct {
	CatID     uint32 `json:"catID"`
	CatName   string `json:"catName"`
	ParentCat uint32 `json:"ParentCat"`
}

type Listing struct {
	ListID      uint32    `json:"listID"`
	ListName    string    `json:"listName"`
	ListDesc    string    `json:"listDesc"`
	Ig_url      string    `json:"ig_url"`
	Fb_url      string    `json:"fb_url"`
	Website_url string    `json:"website_url"`
	UserID      uint32    `json:"id"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

type Pdtsvc struct {
	PdtsvcID    uint32    `json:"pdtsvcID"`
	PdtsvcName  string    `json:"pdtsvcName"`
	PdtsvcPrice float64   `json:"pdtsvcPrice"`
	PdtsvcDesc  string    `json:"pdtsvcDesc"`
	CatID       uint32    `json:"catID"`
	ListID      uint32    `json:"listID"`
	Views       uint32    `json:"views"`
	Likes       uint32    `json:"likes"`
	Keyword     string    `json:"keyword"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

type Review struct {
	ReviewID   uint32    `json:"reviewID"`
	ReviewText string    `json:"reviewText"`
	UserID     uint32    `json:"id"`
	ListID     uint32    `json:"listID"`
	Created    time.Time `json:"created"`
}

type ErrorMsg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
