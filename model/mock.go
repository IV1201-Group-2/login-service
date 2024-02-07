package model

// MockPassword is a plain-text version of password used in tests.
const MockPassword = "password"

// MockPasswordBcrypt is a Bcrypt encoded version of password used in tests.
const MockPasswordBcrypt = "$2a$10$c4WCXRkTtYb3fJ7Wpnjok.nhrEcFyxqpJ/mjfAjBDzqW1IWT6EjVi" // #nosec G101

// MockApplicant is an example user with role "applicant".
var MockApplicant = User{
	ID:   0,
	Role: RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant@example.com",
	Password: MockPasswordBcrypt, // password
}

// MockRecruiter is an example user with role "recruiter".
var MockRecruiter = User{
	ID:   1,
	Role: RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "",
	Password: MockPasswordBcrypt, // password
}
