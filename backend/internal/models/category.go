package models

type Category struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// TODO: Add methods for CRUD operations on categories