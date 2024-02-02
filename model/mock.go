package model

var MockApplicant = User{
	Username: "mockuser_applicant",
	Email:    "mockuser-applicant@example.com",

	Name:           "Mock",
	Surname:        "Applicant",
	PersonalNumber: "20000101-1111",

	Role: RoleApplicant,
}

var MockRecruiter = User{
	Username: "mockuser_recruiter",
	Email:    "mockuser-recruiter@example.com",

	Name:           "Mock",
	Surname:        "Recruiter",
	PersonalNumber: "20000101-9999",

	Role: RoleRecruiter,
}

var MockJWTSigningKey = []byte("Mock secret")
