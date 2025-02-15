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
	mac MACADDR UNIQUE NOT NULL 
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

	createTableSubinterfacesQuery = `
CREATE TABLE IF NOT EXISTS subinterfaces (
	id SERIAL PRIMARY KEY,
	interface_id INT REFERENCES interfaces(id) ON DELETE CASCADE,
	name TEXT NOT NULL
);
`

	createTableSubinterfaceStatesQuery = `
CREATE TABLE IF NOT EXISTS subinterface_states (
	id SERIAL PRIMARY KEY,
	subinterface_id INT REFERENCES subinterfaces(id) ON DELETE CASCADE,
	interface_state_id INT REFERENCES interface_states(id) ON DELETE CASCADE,
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
VALUES (@timestamp);
`

	insertVendorQuery = `
INSERT INTO vendors (name)
VALUES (@vendor)
ON CONFLICT (name) DO NOTHING;
`

	insertOperatingSystemQuery = `
INSERT INTO operating_systems (name, version)
VALUES (@os, @version)
ON CONFLICT (name, version) DO NOTHING;
`

	insertDeviceQuery = `
WITH vendor_id AS (
	SELECT id
	FROM vendors 
	WHERE name = @vendor
), operating_system_id AS (
	SELECT id
	FROM operating_systems
	WHERE name = @os AND version = @version
) 
INSERT INTO devices (vendor_id, operating_system_id, hostname, serial_number)
SELECT v.id, o.id, @hostname, @serial
FROM vendor_id AS v, operating_system_id AS o;
`

	insertDeviceStateQuery = `
WITH snapshot_id AS (
	SELECT id
	FROM snapshots 
	WHERE timestamp = @timestamp
), device_id AS (
	SELECT id
	FROM devices
	WHERE hostname = @hostname
) 
INSERT INTO device_states (snapshot_id, device_id, is_snapshot_successful)
SELECT s.id, d.id, @is_snapshot_successful
FROM snapshot_id AS s, device_id AS d;
`

	insertInterfaceQuery = `
WITH device_id AS (
	SELECT id
	FROM devices 
	WHERE hostname = @hostname
)
INSERT INTO interfaces (device_id, name, mac)
SELECT d.id, @name, @mac
FROM device_id AS d;
`

	insertInterfaceStateQuery = `
WITH interface_id AS (
	SELECT id
	FROM interfaces 
	WHERE mac = @mac
), device_state_id AS (
	SELECT id
	FROM device_states
	WHERE snapshot_id IN (
		SELECT id
		FROM snapshots 
		WHERE timestamp = @timestamp
	) AND device_id IN (
		SELECT id
		FROM devices
		WHERE hostname = @hostname
	)
) 
INSERT INTO interface_states (interface_id, device_state_id, is_up, ip, mtu)
SELECT i.id, d.id, @is_up, @ip, @mtu
FROM interface_id AS i, device_state_id AS d;
`

	insertSubinterfaceQuery = `
WITH interface_id AS (
	SELECT id
	FROM interfaces 
	WHERE mac = @mac
)
INSERT INTO subinterfaces (interface_id, name)
SELECT i.id, @name
FROM interface_id AS i;
`

	insertSubinterfaceStateQuery = `
WITH subinterface_id AS (
	SELECT id
	FROM subinterfaces 
	WHERE interface_id IN (
		SELECT id
		FROM interfaces 
		WHERE mac = @mac
	)
), interface_state_id AS (
	SELECT id
	FROM interface_states
	WHERE interface_id IN (
		SELECT id
		FROM interfaces 
		WHERE mac = @mac
	) AND device_state_id IN (
		SELECT id
		FROM device_states
		WHERE snapshot_id IN (
			SELECT id
			FROM snapshots 
			WHERE timestamp = @timestamp
		) AND device_id IN (
			SELECT id
			FROM devices
			WHERE hostname = @hostname
		)
	)
) 
INSERT INTO subinterface_states (subinterface_id, interface_state_id, is_up, ip, mtu)
SELECT s.id, i.id, @is_up, @ip, @mtu
FROM subinterface_id AS s, interface_state_id AS i;
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
	d.id AS device_id,
	d.hostname,
	d.os_name,
	d.os_version,
	d.serial_number,
	d.management_ip,
	i.name AS interface_name,
	i.mac,
	i.ip,
	i.mtu,
	i.bandwidth
FROM
	snapshots s
	JOIN devices d ON s.id = d.snapshot_id
	JOIN vendors v ON v.id = d.vendor_id
	JOIN interfaces i ON d.id = i.device_id
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
