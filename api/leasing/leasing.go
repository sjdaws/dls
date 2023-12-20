package leasing

import (
    db "github.com/sjdaws/dls/database"
)

type Leasing struct {
    database *db.Database
}

// New creates a new leasing instance
func New(database *db.Database) *Leasing {
    return &Leasing{
        database: database,
    }
}
