package model

// Mock user with role "applicant".
var MockApplicant = User{
	ID:   0,
	Role: RoleApplicant,

	Username: "",
	Email:    "mockuser-applicant@example.com",
	Password: "password",
}

// Mock user with role "recruiter".
var MockRecruiter = User{
	ID:   1,
	Role: RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "",
	Password: "password",
}
