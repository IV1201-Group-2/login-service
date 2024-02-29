package tests

import "github.com/IV1201-Group-2/login-service/model"

// MockPassword is a plain-text version of the password used in tests.
const MockPassword = "password"

// MockPasswordBcrypt is a Bcrypt encoded version of the password used in tests.
const MockPasswordBcrypt = "$2a$10$c4WCXRkTtYb3fJ7Wpnjok.nhrEcFyxqpJ/mjfAjBDzqW1IWT6EjVi" // #nosec G101

// MockSecret is the JWT secret used in tests.
const MockSecret = "mocksecret"

// MockApplicant is an example user with role "applicant".
var MockApplicant = model.User{
	ID:   0,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant@example.com",
	Password: MockPasswordBcrypt, // password
}

// MockApplicant2 is an example user with role "applicant" and a missing password.
var MockApplicant2 = model.User{
	ID:   1,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant2@example.com",
	Password: "",
}

// MockRecruiter is an example user with role "recruiter".
var MockRecruiter = model.User{
	ID:   2,
	Role: model.RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "",
	Password: MockPasswordBcrypt, // password
}
