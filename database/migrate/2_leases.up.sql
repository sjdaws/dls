CREATE TABLE leases (
    reference varchar(36) NOT NULL,
    origin_reference varchar(36) NOT NULL,
    created_at datetime NOT NULL,
    expires_at datetime NOT NULL,
    PRIMARY KEY (reference),
    FOREIGN KEY origin_reference (origin_reference) REFERENCES origins(reference) ON DELETE CASCADE
);
