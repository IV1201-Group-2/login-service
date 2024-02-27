-- Mock database for testing

CREATE TABLE competence (
	competence_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	i18n_key character varying(255)
);

CREATE TABLE role (
	role_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	name character varying(255)
);

CREATE TABLE application_status (
	person_id integer NOT NULL,
	status character varying,
	application_status_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	FOREIGN KEY (person_id) REFERENCES person(person_id)
);

CREATE TABLE availability (
	availability_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	person_id integer,
	from_date date,
	to_date date,
	FOREIGN KEY (person_id) REFERENCES person(person_id)
);

CREATE TABLE competence_profile (
	competence_profile_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	person_id integer,
	competence_id integer,
	years_of_experience numeric(4,2),
	FOREIGN KEY (person_id) REFERENCES person(person_id),
	FOREIGN KEY (competence_id) REFERENCES competence(competence_id)
);

CREATE TABLE person (
	person_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	name character varying(255),
	surname character varying(255),
	pnr character varying(255),
	email character varying(255),
	password character varying(255),
	role_id integer,
	username character varying(255),
	FOREIGN KEY (role_id) REFERENCES role(role_id)
);

INSERT INTO role VALUES (1, 'recruiter');
INSERT INTO role VALUES (2, 'applicant');

-- Applicant with password (login: mockuser-applicant@example.com, password)
INSERT INTO person VALUES (1, 'Mock', 'Applicant', '200001011111', 'mockuser-applicant@example.com', '$2a$10$c4WCXRkTtYb3fJ7Wpnjok.nhrEcFyxqpJ/mjfAjBDzqW1IWT6EjVi', 2, '');
-- Applicant without password (login: mockuser-applicant2@example.com)
INSERT INTO person VALUES (2, 'Mock', 'Applicant 2', '200001012222', 'mockuser-applicant2@example.com', '', 2, '');
-- Recruiter (login: mockuser-recruiter, password)
INSERT INTO person VALUES (3, 'Mock', 'Recruiter', '200001013333', '', '$2a$10$c4WCXRkTtYb3fJ7Wpnjok.nhrEcFyxqpJ/mjfAjBDzqW1IWT6EjVi', 1, 'mockuser-recruiter');
