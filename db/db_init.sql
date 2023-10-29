BEGIN;

SET
	client_encoding = 'LATIN1';

CREATE TABLE mc_merchant (
	id VARCHAR(100) PRIMARY KEY,
	code VARCHAR(100) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mc_team_member (
	id VARCHAR(100) PRIMARY KEY,
	merchant_id VARCHAR(100) NOT NULL REFERENCES mc_merchant (id),
	email VARCHAR(100) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sample data
INSERT INTO
	mc_merchant (id, code)
VALUES
	('6d33b3dc-70d1-453d-a7c2-91308c09e0a7', 'MC01'),
	('8a51ab21-5d66-4a19-b907-999531b8d692', 'MC02');

INSERT INTO
	mc_team_member(id, merchant_id, email)
VALUES
	(
		'54ea366b-c49c-4d8b-a71a-db763bb98a73',
		'6d33b3dc-70d1-453d-a7c2-91308c09e0a7',
		'member1@merchant1.com'
	),
	(
		'9d81bad1-7371-44bb-b0dc-ff30857464d2',
		'6d33b3dc-70d1-453d-a7c2-91308c09e0a7',
		'member2@merchant1.com'
	),
	(
		'64dee758-7409-4099-b14e-fb226c02b7b3',
		'8a51ab21-5d66-4a19-b907-999531b8d692',
		'teammember1@merchant2.com'
	),
	(
		'd2503449-fe5c-4b1e-badd-e9fda8d5d5e2',
		'8a51ab21-5d66-4a19-b907-999531b8d692',
		'teammember2@merchant2.com'
	);

COMMIT;