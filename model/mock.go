package model

// Plain-text version of password used in tests
const MockPassword = "password"

// Bcrypt encoded version of password used in tests
const MockPasswordBcrypt = "$2a$10$c4WCXRkTtYb3fJ7Wpnjok.nhrEcFyxqpJ/mjfAjBDzqW1IWT6EjVi"

// Mock user with role "applicant".
var MockApplicant = User{
	ID:   0,
	Role: RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant@example.com",
	Password: MockPasswordBcrypt, // password
}

// Mock user with role "recruiter".
var MockRecruiter = User{
	ID:   1,
	Role: RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "",
	Password: MockPasswordBcrypt, // password
}
