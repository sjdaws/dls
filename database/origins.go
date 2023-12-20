package database

type Origin struct {
    Hostname  string
    Reference string
}

// CreateOrigin creates a new origin record in the database
func (d *Database) CreateOrigin(origin *Origin) error {
    prepared, err := d.connection.Prepare("INSERT INTO origins VALUES (?, ?) ON DUPLICATE KEY UPDATE reference = ?")
    if err != nil {
        return err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    _, err = stmt.prepared.Exec(
        origin.Reference,
        origin.Hostname,
        origin.Reference,
    )

    return err
}

// GetOrigin finds an origin record by reference
func (d *Database) GetOrigin(reference string) (*Origin, error) {
    prepared, err := d.connection.Prepare("SELECT hostname, reference FROM origins WHERE reference = ?")
    if err != nil {
        return nil, err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    origin := &Origin{}
    err = stmt.prepared.QueryRow(reference).Scan(&origin.Hostname, &origin.Reference)
    if err != nil {
        return nil, err
    }

    return origin, nil
}
