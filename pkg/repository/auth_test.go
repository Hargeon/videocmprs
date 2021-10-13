package repository

// Deprecated
/*
func TestGetUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		email        string
		password     string
		mock         func()
		expectedId   int64
		errorPresent bool
	}{
		{
			name:     "With valid email and password",
			email:    "chech@check.com",
			password: "oiojhoh",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("chech@check.com", "oiojhoh").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedId:   1,
			errorPresent: false,
		},
		{
			name:     "With invalid email",
			email:    "",
			password: "oiojhoh",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("", "oiojhoh").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewAuthRepository(db)
			id, err := repo.GetUser(testCase.email, testCase.password)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}
			if id != testCase.expectedId {
				t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
*/
