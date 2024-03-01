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
// This user should not be modified during the tests.
var MockApplicant2 = model.User{
	ID:   1,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant2@example.com",
	Password: "",
}

// MockApplicant3 is an example user with role "applicant" and a missing password.
// This user should be modified during the API tests.
var MockApplicant3 = model.User{
	ID:   2,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant3@example.com",
	Password: "",
}

// MockApplicant3 is an example user with role "applicant" and a missing password.
// This user should be modified during the service tests.
var MockApplicant4 = model.User{
	ID:   3,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant4@example.com",
	Password: "",
}

// MockApplicant4 is an example user with role "applicant" and a missing password.
// This user should be modified during the database tests.
var MockApplicant5 = model.User{
	ID:   4,
	Role: model.RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant5@example.com",
	Password: "",
}

// MockRecruiter is an example user with role "recruiter".
var MockRecruiter = model.User{
	ID:   5,
	Role: model.RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "",
	Password: MockPasswordBcrypt, // password
}
