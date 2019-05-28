package repository



/*
func TestGetUserStore(t *testing.T) {
	Init()
	if _, v := GetUserStore("primary"); v != nil {
		t.Errorf("Expected primary store but got error %v", v)
	}
}

func TestGetUserByUsername(t *testing.T) {
	Init()
	store, _ := GetUserStore("primary")
	AddNewUser(*store, "testuser", "testuser")
	if _, v := GetUserByID("testuser"); v != nil {
		t.Errorf("Expected user in primary store but got error %v", v)
	}
}

func TestGetUserByUsername2(t *testing.T) {
	Init()
	store, _ := GetUserStore("primary")
	AddNewUser(*store, "testuserx", "testuserx")
	if _, v := GetUserByID("userdonotexist"); v != ErrUserNotFound {
		t.Errorf("Expected user not found error but got %v", v)
	}
}

func TestValidateUser(t *testing.T) {
	Init()
	store, _ := GetUserStore("primary")
	if err := ValidateUser(*store, "admin", "admin"); err != nil {
		t.Errorf("Expected nil error but got %v", err)
	}
}
*/
