package database

import (
    "time"
)

type Lease struct {
    CreatedAt       time.Time
    ExpiresAt       time.Time
    OriginReference string
    Reference       string
}

// CreateLease creates a new lease record in the database
func (d *Database) CreateLease(lease *Lease) error {
    prepared, err := d.connection.Prepare("INSERT INTO leases VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE reference = ? AND origin_reference = ?")
    if err != nil {
        return err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    _, err = stmt.prepared.Exec(
        lease.Reference,
        lease.OriginReference,
        lease.CreatedAt,
        lease.ExpiresAt,
        lease.Reference,
        lease.OriginReference,
    )

    return err
}

// DeleteLease removes a lease from the database
func (d *Database) DeleteLease(lease *Lease) error {
    prepared, err := d.connection.Prepare("DELETE FROM leases WHERE reference = ? AND origin_reference = ?")
    if err != nil {
        return err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    _, err = stmt.prepared.Exec(lease.Reference, lease.OriginReference)
    if err != nil {
        return err
    }

    return nil
}

// DeleteLeases removes lease records by origin reference
func (d *Database) DeleteLeases(originReference string) error {
    prepared, err := d.connection.Prepare("DELETE FROM leases WHERE origin_reference = ?")
    if err != nil {
        return err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    _, err = stmt.prepared.Exec(originReference)
    if err != nil {
        return err
    }

    return nil
}

// GetLease finds a single lease record
func (d *Database) GetLease(reference string, originReference string) (*Lease, error) {
    prepared, err := d.connection.Prepare("SELECT created_at, expires_at, origin_reference, reference FROM leases WHERE reference = ? AND origin_reference = ?")
    if err != nil {
        return nil, err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    lease := &Lease{}
    err = stmt.prepared.QueryRow(reference, originReference).Scan(&lease.CreatedAt, &lease.ExpiresAt, &lease.OriginReference, &lease.Reference)
    if err != nil {
        return nil, err
    }

    return lease, nil
}

// GetLeases finds lease records by origin reference
func (d *Database) GetLeases(originReference string) ([]*Lease, error) {
    prepared, err := d.connection.Query("SELECT created_at, expires_at, origin_reference, reference FROM leases WHERE origin_reference = ?", originReference)
    if err != nil {
        return nil, err
    }

    results := &query{
        rows: prepared,
    }
    defer results.close()

    var leases []*Lease
    for results.rows.Next() {
        var lease Lease
        err = results.rows.Scan(&lease.CreatedAt, &lease.ExpiresAt, &lease.OriginReference, &lease.Reference)
        if err != nil {
            return nil, err
        }
        leases = append(leases, &lease)
    }

    return leases, nil
}

// UpdateLease updates a lease in the database
func (d *Database) UpdateLease(lease *Lease) error {
    prepared, err := d.connection.Prepare("UPDATE leases SET expires_at = ? WHERE reference = ? AND origin_reference = ?")
    if err != nil {
        return err
    }

    stmt := &statement{
        prepared: prepared,
    }
    defer stmt.close()

    _, err = stmt.prepared.Exec(
        lease.ExpiresAt,
        lease.Reference,
        lease.OriginReference,
    )

    return err
}
