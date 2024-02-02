package model

var MockApplicant = User{
	Username: "mockuser_applicant",
	Email:    "mockuser-applicant@example.com",

	Name:           "Mock",
	Surname:        "Applicant",
	PersonalNumber: "200001011111",

	Role: RoleApplicant,
}

var MockRecruiter = User{
	Username: "mockuser_recruiter",
	Email:    "mockuser-recruiter@example.com",

	Name:           "Mock",
	Surname:        "Recruiter",
	PersonalNumber: "200001019999",

	Role: RoleRecruiter,
}

var MockJWTSigningKey = []byte("Mock secret")
