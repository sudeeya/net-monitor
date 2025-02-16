package postgresql

// SQL queries to create tables needed to store snapshots.
const (
	createTableSnapshotsQuery = `
CREATE TABLE IF NOT EXISTS snapshots (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMPTZ UNIQUE NOT NULL
);
`

	createTableVendorsQuery = `
CREATE TABLE IF NOT EXISTS vendors (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);
`

	createTableOperatingSystemsQuery = `
CREATE TABLE IF NOT EXISTS operating_systems (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	version TEXT,
	UNIQUE (name, version)
);
`

	createTableDevicesQuery = `
CREATE TABLE IF NOT EXISTS devices (
	id SERIAL PRIMARY KEY,
	vendor_id INT REFERENCES vendors(id) ON DELETE RESTRICT,
	operating_system_id INT REFERENCES operating_systems(id) ON DELETE RESTRICT,
	hostname TEXT UNIQUE NOT NULL,
	serial_number TEXT
);
`

	createTableDeviceStatesQuery = `
CREATE TABLE IF NOT EXISTS device_states (
	id SERIAL PRIMARY KEY,
	snapshot_id INT REFERENCES snapshots(id) ON DELETE CASCADE,
	device_id INT REFERENCES devices(id) ON DELETE RESTRICT,
	is_snapshot_successful BOOLEAN NOT NULL
);
`

	createTableInterfacesQuery = `
CREATE TABLE IF NOT EXISTS interfaces (
	id SERIAL PRIMARY KEY,
	device_id INT REFERENCES devices(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	UNIQUE (device_id, name)
);
`

	createTableInterfaceStatesQuery = `
CREATE TABLE IF NOT EXISTS interface_states (
	id SERIAL PRIMARY KEY,
	interface_id INT REFERENCES interfaces(id) ON DELETE CASCADE,
	device_state_id INT REFERENCES device_states(id) ON DELETE CASCADE,
	is_up BOOLEAN NOT NULL,
	ip INET,
	mtu INT
);
`
)

// SQL queries for inserting a snapshot.
const (
	insertSnapshotQuery = `
INSERT INTO snapshots (timestamp)
VALUES (@timestamp)
RETURNING id;
`

	insertVendorQuery = `
WITH insert_vendor AS (
	INSERT INTO vendors (name)
	VALUES (@vendor)
	ON CONFLICT (name) DO NOTHING
	RETURNING id
)
SELECT id 
FROM insert_vendor
UNION
SELECT id
FROM vendors 
WHERE name = @vendor;
`

	insertOperatingSystemQuery = `
WITH insert_operating_system AS (
	INSERT INTO operating_systems (name, version)
	VALUES (@os, @version)
	ON CONFLICT (name, version) DO NOTHING
	RETURNING id
)
SELECT id 
FROM insert_operating_system
UNION
SELECT id
FROM operating_systems 
WHERE name = @os AND version = @version;
`

	insertDeviceQuery = `
WITH insert_device AS (
	INSERT INTO devices (vendor_id, operating_system_id, hostname, serial_number)
	VALUES (@vendor_id, @operating_system_id, @hostname, @serial_number)
	ON CONFLICT (hostname) DO NOTHING
	RETURNING id
)
SELECT id 
FROM insert_device
UNION
SELECT id
FROM devices 
WHERE hostname = @hostname;
`

	insertDeviceStateQuery = `
INSERT INTO device_states (snapshot_id, device_id, is_snapshot_successful)
VALUES (@snapshot_id, @device_id, @is_snapshot_successful)
RETURNING id;
`

	insertInterfaceQuery = `
WITH insert_interface AS (
	INSERT INTO interfaces (device_id, name)
	VALUES (@device_id, @name)
	ON CONFLICT (device_id, name) DO NOTHING
	RETURNING id
)
SELECT id 
FROM insert_interface
UNION
SELECT id
FROM interfaces 
WHERE device_id = @device_id AND name = @name;
`

	insertInterfaceStateQuery = `
INSERT INTO interface_states (interface_id, device_state_id, is_up, ip, mtu)
VALUES (@interface_id, @device_state_id, @is_up, @ip, @mtu);
`
)

// SQL query to get snapshot ids and timestamps.
const (
	selectTimestampsQuery = `
SELECT id, timestamp
FROM snapshots
ORDER BY timestamp DESC
LIMIT @limit;
`
)

// SQL query to get a snapshot.
const (
	selectSnapshotQuery = `
SELECT
	s.timestamp,
	v.name AS vendor_name,
	o.name AS os_name,
	o.version AS os_version,
	d.id AS device_id,
	d.hostname,
	d.serial_number,
	ds.is_snapshot_successful,
	i.name AS interface_name,
	is.is_up AS interface_is_up,
	is.ip,
	is.mtu
FROM
	devices AS d
	JOIN vendors AS v ON v.id = d.vendor_id
	JOIN operating_systems AS o ON o.id = d.operating_system_id
	JOIN device_states AS ds ON d.id = ds.device_id
	JOIN snapshots AS s ON s.id = ds.snapshot_id
	JOIN interfaces AS i ON d.id = i.device_id
	JOIN interface_states AS is ON i.id = is.interface_id AND ds.id = is.device_state_id
WHERE
	s.id = @id
ORDER BY device_id ASC;
`
)

// SQL query to delete a snapshot.
const (
	deleteSnapshotQuery = `
DELETE FROM snapshots
WHERE id = @id;
`
)
