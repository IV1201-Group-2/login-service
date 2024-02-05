package model

var MockApplicant = User{
	ID:   0,
	Role: RoleApplicant,

	Username: "mockuser_applicant",
	Email:    "mockuser-applicant@example.com",
	Password: "password",
}

var MockRecruiter = User{
	ID:   1,
	Role: RoleRecruiter,

	Username: "mockuser_recruiter",
	Email:    "mockuser-recruiter@example.com",
	Password: "password",
}

var MockJWTSigningKey = []byte("Mock secret")
