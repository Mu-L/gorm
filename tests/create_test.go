package tests_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	. "gorm.io/gorm/utils/tests"
)

func TestCreate(t *testing.T) {
	u1 := *GetUser("create", Config{})

	if results := DB.Create(&u1); results.Error != nil {
		t.Fatalf("errors happened when create: %v", results.Error)
	} else if results.RowsAffected != 1 {
		t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
	}

	if u1.ID == 0 {
		t.Errorf("user's primary key should has value after create, got : %v", u1.ID)
	}

	if u1.CreatedAt.IsZero() {
		t.Errorf("user's created at should be not zero")
	}

	if u1.UpdatedAt.IsZero() {
		t.Errorf("user's updated at should be not zero")
	}

	var newUser User
	if err := DB.Where("id = ?", u1.ID).First(&newUser).Error; err != nil {
		t.Fatalf("errors happened when query: %v", err)
	} else {
		CheckUser(t, newUser, u1)
	}

	type user struct {
		ID   int `gorm:"primaryKey;->:false"`
		Name string
		Age  int
	}

	var u2 user
	if results := DB.Create(&u2); results.Error != nil {
		t.Fatalf("errors happened when create: %v", results.Error)
	} else if results.RowsAffected != 1 {
		t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
	}

	if u2.ID != 0 {
		t.Errorf("don't have the permission to read primary key from db, but got %v", u2.ID)
	}
}

func TestCreateInBatches(t *testing.T) {
	users := []User{
		*GetUser("create_in_batches_1", Config{Account: true, Pets: 2, Toys: 3, Company: true, Manager: true, Team: 0, Languages: 1, Friends: 1}),
		*GetUser("create_in_batches_2", Config{Account: false, Pets: 2, Toys: 4, Company: false, Manager: false, Team: 1, Languages: 3, Friends: 5}),
		*GetUser("create_in_batches_3", Config{Account: true, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 4, Languages: 0, Friends: 1}),
		*GetUser("create_in_batches_4", Config{Account: true, Pets: 3, Toys: 0, Company: false, Manager: true, Team: 0, Languages: 3, Friends: 0}),
		*GetUser("create_in_batches_5", Config{Account: false, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 1, Languages: 3, Friends: 1}),
		*GetUser("create_in_batches_6", Config{Account: true, Pets: 4, Toys: 3, Company: false, Manager: true, Team: 1, Languages: 3, Friends: 0}),
	}

	result := DB.CreateInBatches(&users, 2)
	if result.RowsAffected != int64(len(users)) {
		t.Errorf("affected rows should be %v, but got %v", len(users), result.RowsAffected)
	}

	for _, user := range users {
		if user.ID == 0 {
			t.Fatalf("failed to fill user's ID, got %v", user.ID)
		} else {
			var newUser User
			if err := DB.Where("id = ?", user.ID).Preload(clause.Associations).First(&newUser).Error; err != nil {
				t.Fatalf("errors happened when query: %v", err)
			} else {
				CheckUser(t, newUser, user)
			}
		}
	}
}

func TestCreateInBatchesWithDefaultSize(t *testing.T) {
	users := []User{
		*GetUser("create_with_default_batch_size_1", Config{Account: true, Pets: 2, Toys: 3, Company: true, Manager: true, Team: 0, Languages: 1, Friends: 1}),
		*GetUser("create_with_default_batch_sizs_2", Config{Account: false, Pets: 2, Toys: 4, Company: false, Manager: false, Team: 1, Languages: 3, Friends: 5}),
		*GetUser("create_with_default_batch_sizs_3", Config{Account: true, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 4, Languages: 0, Friends: 1}),
		*GetUser("create_with_default_batch_sizs_4", Config{Account: true, Pets: 3, Toys: 0, Company: false, Manager: true, Team: 0, Languages: 3, Friends: 0}),
		*GetUser("create_with_default_batch_sizs_5", Config{Account: false, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 1, Languages: 3, Friends: 1}),
		*GetUser("create_with_default_batch_sizs_6", Config{Account: true, Pets: 4, Toys: 3, Company: false, Manager: true, Team: 1, Languages: 3, Friends: 0}),
	}

	result := DB.Session(&gorm.Session{CreateBatchSize: 2}).Create(&users)
	if result.RowsAffected != int64(len(users)) {
		t.Errorf("affected rows should be %v, but got %v", len(users), result.RowsAffected)
	}

	for _, user := range users {
		if user.ID == 0 {
			t.Fatalf("failed to fill user's ID, got %v", user.ID)
		} else {
			var newUser User
			if err := DB.Where("id = ?", user.ID).Preload(clause.Associations).First(&newUser).Error; err != nil {
				t.Fatalf("errors happened when query: %v", err)
			} else {
				CheckUser(t, newUser, user)
			}
		}
	}
}

func TestCreateFromMap(t *testing.T) {
	if err := DB.Model(&User{}).Create(map[string]interface{}{"Name": "create_from_map", "Age": 18}).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}

	var result User
	if err := DB.Where("name = ?", "create_from_map").First(&result).Error; err != nil || result.Age != 18 {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	if err := DB.Model(&User{}).Create(map[string]interface{}{"name": "create_from_map_1", "age": 18}).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}

	var result1 User
	if err := DB.Where("name = ?", "create_from_map_1").First(&result1).Error; err != nil || result1.Age != 18 {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	datas := []map[string]interface{}{
		{"Name": "create_from_map_2", "Age": 19},
		{"name": "create_from_map_3", "Age": 20},
	}

	if err := DB.Model(&User{}).Create(&datas).Error; err != nil {
		t.Fatalf("failed to create data from slice of map, got error: %v", err)
	}

	var result2 User
	if err := DB.Where("name = ?", "create_from_map_2").First(&result2).Error; err != nil || result2.Age != 19 {
		t.Fatalf("failed to query data after create from slice of map, got error %v", err)
	}

	var result3 User
	if err := DB.Where("name = ?", "create_from_map_3").First(&result3).Error; err != nil || result3.Age != 20 {
		t.Fatalf("failed to query data after create from slice of map, got error %v", err)
	}
}

func TestCreateWithAssociations(t *testing.T) {
	user := *GetUser("create_with_associations", Config{
		Account:   true,
		Pets:      2,
		Toys:      3,
		Company:   true,
		Manager:   true,
		Team:      4,
		Languages: 3,
		Friends:   1,
	})

	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("errors happened when create: %v", err)
	}

	CheckUser(t, user, user)

	var user2 User
	DB.Preload("Account").Preload("Pets").Preload("Toys").Preload("Company").Preload("Manager").Preload("Team").Preload("Languages").Preload("Friends").Find(&user2, "id = ?", user.ID)
	CheckUser(t, user2, user)
}

func TestBulkCreateWithAssociations(t *testing.T) {
	users := []User{
		*GetUser("bulk_1", Config{Account: true, Pets: 2, Toys: 3, Company: true, Manager: true, Team: 0, Languages: 1, Friends: 1}),
		*GetUser("bulk_2", Config{Account: false, Pets: 2, Toys: 4, Company: false, Manager: false, Team: 1, Languages: 3, Friends: 5}),
		*GetUser("bulk_3", Config{Account: true, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 4, Languages: 0, Friends: 1}),
		*GetUser("bulk_4", Config{Account: true, Pets: 3, Toys: 0, Company: false, Manager: true, Team: 0, Languages: 3, Friends: 0}),
		*GetUser("bulk_5", Config{Account: false, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 1, Languages: 3, Friends: 1}),
		*GetUser("bulk_6", Config{Account: true, Pets: 4, Toys: 3, Company: false, Manager: true, Team: 1, Languages: 3, Friends: 0}),
		*GetUser("bulk_7", Config{Account: true, Pets: 1, Toys: 3, Company: true, Manager: true, Team: 4, Languages: 3, Friends: 1}),
		*GetUser("bulk_8", Config{Account: false, Pets: 0, Toys: 0, Company: false, Manager: false, Team: 0, Languages: 0, Friends: 0}),
	}

	if results := DB.Create(&users); results.Error != nil {
		t.Fatalf("errors happened when create: %v", results.Error)
	} else if results.RowsAffected != int64(len(users)) {
		t.Fatalf("rows affected expects: %v, got %v", len(users), results.RowsAffected)
	}

	var userIDs []uint
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
		CheckUser(t, user, user)
	}

	var users2 []User
	DB.Preload("Account").Preload("Pets").Preload("Toys").Preload("Company").Preload("Manager").Preload("Team").Preload("Languages").Preload("Friends").Find(&users2, "id IN ?", userIDs)
	for idx, user := range users2 {
		CheckUser(t, user, users[idx])
	}
}

func TestBulkCreatePtrDataWithAssociations(t *testing.T) {
	users := []*User{
		GetUser("bulk_ptr_1", Config{Account: true, Pets: 2, Toys: 3, Company: true, Manager: true, Team: 0, Languages: 1, Friends: 1}),
		GetUser("bulk_ptr_2", Config{Account: false, Pets: 2, Toys: 4, Company: false, Manager: false, Team: 1, Languages: 3, Friends: 5}),
		GetUser("bulk_ptr_3", Config{Account: true, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 4, Languages: 0, Friends: 1}),
		GetUser("bulk_ptr_4", Config{Account: true, Pets: 3, Toys: 0, Company: false, Manager: true, Team: 0, Languages: 3, Friends: 0}),
		GetUser("bulk_ptr_5", Config{Account: false, Pets: 0, Toys: 3, Company: true, Manager: false, Team: 1, Languages: 3, Friends: 1}),
		GetUser("bulk_ptr_6", Config{Account: true, Pets: 4, Toys: 3, Company: false, Manager: true, Team: 1, Languages: 3, Friends: 0}),
		GetUser("bulk_ptr_7", Config{Account: true, Pets: 1, Toys: 3, Company: true, Manager: true, Team: 4, Languages: 3, Friends: 1}),
		GetUser("bulk_ptr_8", Config{Account: false, Pets: 0, Toys: 0, Company: false, Manager: false, Team: 0, Languages: 0, Friends: 0}),
	}

	if err := DB.Create(&users).Error; err != nil {
		t.Fatalf("errors happened when create: %v", err)
	}

	var userIDs []uint
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
		CheckUser(t, *user, *user)
	}

	var users2 []User
	DB.Preload("Account").Preload("Pets").Preload("Toys").Preload("Company").Preload("Manager").Preload("Team").Preload("Languages").Preload("Friends").Find(&users2, "id IN ?", userIDs)
	for idx, user := range users2 {
		CheckUser(t, user, *users[idx])
	}
}

func TestPolymorphicHasOne(t *testing.T) {
	t.Run("Struct", func(t *testing.T) {
		pet := Pet{
			Name: "PolymorphicHasOne",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne"},
		}

		if err := DB.Create(&pet).Error; err != nil {
			t.Fatalf("errors happened when create: %v", err)
		}

		CheckPet(t, pet, pet)

		var pet2 Pet
		DB.Preload("Toy").Find(&pet2, "id = ?", pet.ID)
		CheckPet(t, pet2, pet)
	})

	t.Run("Slice", func(t *testing.T) {
		pets := []Pet{{
			Name: "PolymorphicHasOne-Slice-1",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-1"},
		}, {
			Name: "PolymorphicHasOne-Slice-2",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-2"},
		}, {
			Name: "PolymorphicHasOne-Slice-3",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-3"},
		}}

		if err := DB.Create(&pets).Error; err != nil {
			t.Fatalf("errors happened when create: %v", err)
		}

		var petIDs []uint
		for _, pet := range pets {
			petIDs = append(petIDs, pet.ID)
			CheckPet(t, pet, pet)
		}

		var pets2 []Pet
		DB.Preload("Toy").Find(&pets2, "id IN ?", petIDs)
		for idx, pet := range pets2 {
			CheckPet(t, pet, pets[idx])
		}
	})

	t.Run("SliceOfPtr", func(t *testing.T) {
		pets := []*Pet{{
			Name: "PolymorphicHasOne-Slice-1",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-1"},
		}, {
			Name: "PolymorphicHasOne-Slice-2",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-2"},
		}, {
			Name: "PolymorphicHasOne-Slice-3",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Slice-3"},
		}}

		if err := DB.Create(&pets).Error; err != nil {
			t.Fatalf("errors happened when create: %v", err)
		}

		for _, pet := range pets {
			CheckPet(t, *pet, *pet)
		}
	})

	t.Run("Array", func(t *testing.T) {
		pets := [...]Pet{{
			Name: "PolymorphicHasOne-Array-1",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-1"},
		}, {
			Name: "PolymorphicHasOne-Array-2",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-2"},
		}, {
			Name: "PolymorphicHasOne-Array-3",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-3"},
		}}

		if err := DB.Create(&pets).Error; err != nil {
			t.Fatalf("errors happened when create: %v", err)
		}

		for _, pet := range pets {
			CheckPet(t, pet, pet)
		}
	})

	t.Run("ArrayPtr", func(t *testing.T) {
		pets := [...]*Pet{{
			Name: "PolymorphicHasOne-Array-1",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-1"},
		}, {
			Name: "PolymorphicHasOne-Array-2",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-2"},
		}, {
			Name: "PolymorphicHasOne-Array-3",
			Toy:  Toy{Name: "Toy-PolymorphicHasOne-Array-3"},
		}}

		if err := DB.Create(&pets).Error; err != nil {
			t.Fatalf("errors happened when create: %v", err)
		}

		for _, pet := range pets {
			CheckPet(t, *pet, *pet)
		}
	})
}

func TestCreateEmptyStruct(t *testing.T) {
	type EmptyStruct struct {
		ID uint
	}
	DB.Migrator().DropTable(&EmptyStruct{})

	if err := DB.AutoMigrate(&EmptyStruct{}); err != nil {
		t.Errorf("no error should happen when auto migrate, but got %v", err)
	}

	if err := DB.Create(&EmptyStruct{}).Error; err != nil {
		t.Errorf("No error should happen when creating user, but got %v", err)
	}
}

func TestCreateEmptySlice(t *testing.T) {
	data := []User{}
	if err := DB.Create(&data).Error; err != gorm.ErrEmptySlice {
		t.Errorf("no data should be created, got %v", err)
	}

	sliceMap := []map[string]interface{}{}
	if err := DB.Model(&User{}).Create(&sliceMap).Error; err != gorm.ErrEmptySlice {
		t.Errorf("no data should be created, got %v", err)
	}
}

func TestCreateInvalidSlice(t *testing.T) {
	users := []*User{
		GetUser("invalid_slice_1", Config{}),
		GetUser("invalid_slice_2", Config{}),
		nil,
	}

	if err := DB.Create(&users).Error; !errors.Is(err, gorm.ErrInvalidData) {
		t.Errorf("should returns error invalid data when creating from slice that contains invalid data")
	}
}

func TestCreateWithExistingTimestamp(t *testing.T) {
	user := User{Name: "CreateUserExistingTimestamp"}
	curTime := now.MustParse("2016-01-01")
	user.CreatedAt = curTime
	user.UpdatedAt = curTime
	DB.Save(&user)

	AssertEqual(t, user.CreatedAt, curTime)
	AssertEqual(t, user.UpdatedAt, curTime)

	var newUser User
	DB.First(&newUser, user.ID)

	AssertEqual(t, newUser.CreatedAt, curTime)
	AssertEqual(t, newUser.UpdatedAt, curTime)
}

func TestCreateWithNowFuncOverride(t *testing.T) {
	user := User{Name: "CreateUserTimestampOverride"}
	curTime := now.MustParse("2016-01-01")

	NEW := DB.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return curTime
		},
	})

	NEW.Save(&user)

	AssertEqual(t, user.CreatedAt, curTime)
	AssertEqual(t, user.UpdatedAt, curTime)

	var newUser User
	NEW.First(&newUser, user.ID)

	AssertEqual(t, newUser.CreatedAt, curTime)
	AssertEqual(t, newUser.UpdatedAt, curTime)
}

func TestCreateWithNoGORMPrimaryKey(t *testing.T) {
	type JoinTable struct {
		UserID   uint
		FriendID uint
	}

	DB.Migrator().DropTable(&JoinTable{})
	if err := DB.AutoMigrate(&JoinTable{}); err != nil {
		t.Errorf("no error should happen when auto migrate, but got %v", err)
	}

	jt := JoinTable{UserID: 1, FriendID: 2}
	err := DB.Create(&jt).Error
	if err != nil {
		t.Errorf("No error should happen when create a record without a GORM primary key. But in the database this primary key exists and is the union of 2 or more fields\n But got: %s", err)
	}
}

func TestSelectWithCreate(t *testing.T) {
	user := *GetUser("select_create", Config{Account: true, Pets: 3, Toys: 3, Company: true, Manager: true, Team: 3, Languages: 3, Friends: 4})
	DB.Select("Account", "Toys", "Manager", "ManagerID", "Languages", "Name", "CreatedAt", "Age", "Active").Create(&user)

	var user2 User
	DB.Preload("Account").Preload("Pets").Preload("Toys").Preload("Company").Preload("Manager").Preload("Team").Preload("Languages").Preload("Friends").First(&user2, user.ID)

	user.Birthday = nil
	user.Pets = nil
	user.Company = Company{}
	user.Team = nil
	user.Friends = nil

	CheckUser(t, user2, user)
}

func TestOmitWithCreate(t *testing.T) {
	user := *GetUser("omit_create", Config{Account: true, Pets: 3, Toys: 3, Company: true, Manager: true, Team: 3, Languages: 3, Friends: 4})
	DB.Omit("Account", "Toys", "Manager", "Birthday").Create(&user)

	var result User
	DB.Preload("Account").Preload("Pets").Preload("Toys").Preload("Company").Preload("Manager").Preload("Team").Preload("Languages").Preload("Friends").First(&result, user.ID)

	user.Birthday = nil
	user.Account = Account{}
	user.Toys = nil
	user.Manager = nil

	CheckUser(t, result, user)

	user2 := *GetUser("omit_create", Config{Account: true, Pets: 3, Toys: 3, Company: true, Manager: true, Team: 3, Languages: 3, Friends: 4})
	DB.Omit(clause.Associations).Create(&user2)

	var result2 User
	DB.Preload(clause.Associations).First(&result2, user2.ID)

	user2.Account = Account{}
	user2.Toys = nil
	user2.Manager = nil
	user2.Company = Company{}
	user2.Pets = nil
	user2.Team = nil
	user2.Languages = nil
	user2.Friends = nil

	CheckUser(t, result2, user2)
}

func TestFirstOrCreateNotExistsTable(t *testing.T) {
	company := Company{Name: "first_or_create_if_not_exists_table"}
	if err := DB.Table("not_exists").FirstOrCreate(&company).Error; err == nil {
		t.Errorf("not exists table, but err is nil")
	}
}

func TestFirstOrCreateWithPrimaryKey(t *testing.T) {
	company := Company{ID: 100, Name: "company100_with_primarykey"}
	DB.FirstOrCreate(&company)

	if company.ID != 100 {
		t.Errorf("invalid primary key after creating, got %v", company.ID)
	}

	companies := []Company{
		{ID: 101, Name: "company101_with_primarykey"},
		{ID: 102, Name: "company102_with_primarykey"},
	}
	DB.Create(&companies)

	if companies[0].ID != 101 || companies[1].ID != 102 {
		t.Errorf("invalid primary key after creating, got %v, %v", companies[0].ID, companies[1].ID)
	}
}

func TestCreateFromSubQuery(t *testing.T) {
	user := User{Name: "jinzhu"}

	DB.Create(&user)

	subQuery := DB.Table("users").Where("name=?", user.Name).Select("id")

	result := DB.Session(&gorm.Session{DryRun: true}).Model(&Pet{}).Create([]map[string]interface{}{
		{
			"name":    "cat",
			"user_id": gorm.Expr("(?)", DB.Table("(?) as tmp", subQuery).Select("@uid:=id")),
		},
		{
			"name":    "dog",
			"user_id": gorm.Expr("@uid"),
		},
	})

	if !regexp.MustCompile(`INSERT INTO .pets. \(.name.,.user_id.\) .*VALUES \(.+,\(SELECT @uid:=id FROM \(SELECT id FROM .users. WHERE name=.+\) as tmp\)\),\(.+,@uid\)`).MatchString(result.Statement.SQL.String()) {
		t.Errorf("invalid insert SQL, got %v", result.Statement.SQL.String())
	}
}

func TestCreateNilPointer(t *testing.T) {
	var user *User

	err := DB.Create(user).Error
	if err == nil || err != gorm.ErrInvalidValue {
		t.Fatalf("it is not ErrInvalidValue")
	}
}

func TestFirstOrCreateRowsAffected(t *testing.T) {
	user := User{Name: "TestFirstOrCreateRowsAffected"}

	res := DB.FirstOrCreate(&user, "name = ?", user.Name)
	if res.Error != nil || res.RowsAffected != 1 {
		t.Fatalf("first or create rows affect err:%v rows:%d", res.Error, res.RowsAffected)
	}

	res = DB.FirstOrCreate(&user, "name = ?", user.Name)
	if res.Error != nil || res.RowsAffected != 0 {
		t.Fatalf("first or create rows affect err:%v rows:%d", res.Error, res.RowsAffected)
	}
}

func TestCreateWithAutoIncrementCompositeKey(t *testing.T) {
	type CompositeKeyProduct struct {
		ProductID    int `gorm:"primaryKey;autoIncrement:true;"` // primary key
		LanguageCode int `gorm:"primaryKey;"`                    // primary key
		Code         string
		Name         string
	}

	if err := DB.Migrator().DropTable(&CompositeKeyProduct{}); err != nil {
		t.Fatalf("failed to migrate, got error %v", err)
	}
	if err := DB.AutoMigrate(&CompositeKeyProduct{}); err != nil {
		t.Fatalf("failed to migrate, got error %v", err)
	}

	prod := &CompositeKeyProduct{
		LanguageCode: 56,
		Code:         "Code56",
		Name:         "ProductName56",
	}
	if err := DB.Create(&prod).Error; err != nil {
		t.Fatalf("failed to create, got error %v", err)
	}

	newProd := &CompositeKeyProduct{}
	if err := DB.First(&newProd).Error; err != nil {
		t.Fatalf("errors happened when query: %v", err)
	} else {
		AssertObjEqual(t, newProd, prod, "ProductID", "LanguageCode", "Code", "Name")
	}
}

func TestCreateOnConflictWithDefaultNull(t *testing.T) {
	type OnConflictUser struct {
		ID     string
		Name   string `gorm:"default:null"`
		Email  string
		Mobile string `gorm:"default:'133xxxx'"`
	}

	err := DB.Migrator().DropTable(&OnConflictUser{})
	AssertEqual(t, err, nil)
	err = DB.AutoMigrate(&OnConflictUser{})
	AssertEqual(t, err, nil)

	u := OnConflictUser{
		ID:     "on-conflict-user-id",
		Name:   "on-conflict-user-name",
		Email:  "on-conflict-user-email",
		Mobile: "on-conflict-user-mobile",
	}
	err = DB.Create(&u).Error
	AssertEqual(t, err, nil)

	u.Name = "on-conflict-user-name-2"
	u.Email = "on-conflict-user-email-2"
	u.Mobile = ""
	err = DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&u).Error
	AssertEqual(t, err, nil)

	var u2 OnConflictUser
	err = DB.Where("id = ?", u.ID).First(&u2).Error
	AssertEqual(t, err, nil)
	AssertEqual(t, u2.Name, "on-conflict-user-name-2")
	AssertEqual(t, u2.Email, "on-conflict-user-email-2")
	AssertEqual(t, u2.Mobile, "133xxxx")
}

func TestCreateFromMapWithoutPK(t *testing.T) {
	if !isMysql() {
		t.Skipf("This test case skipped, because of only supporting for mysql")
	}

	// case 1: one record, create from map[string]interface{}
	mapValue1 := map[string]interface{}{"name": "create_from_map_with_schema1", "age": 1}
	if err := DB.Model(&User{}).Create(mapValue1).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}

	if _, ok := mapValue1["id"]; !ok {
		t.Fatal("failed to create data from map with table, returning map has no primary key")
	}

	var result1 User
	if err := DB.Where("name = ?", "create_from_map_with_schema1").First(&result1).Error; err != nil || result1.Age != 1 {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	var idVal int64
	_, ok := mapValue1["id"].(uint)
	if ok {
		t.Skipf("This test case skipped, because the db supports returning")
	}

	idVal, ok = mapValue1["id"].(int64)
	if !ok {
		t.Fatal("ret result missing id")
	}

	if int64(result1.ID) != idVal {
		t.Fatal("failed to create data from map with table, @id != id")
	}

	// case2: one record, create from *map[string]interface{}
	mapValue2 := map[string]interface{}{"name": "create_from_map_with_schema2", "age": 1}
	if err := DB.Model(&User{}).Create(&mapValue2).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}

	if _, ok := mapValue2["id"]; !ok {
		t.Fatal("failed to create data from map with table, returning map has no primary key")
	}

	var result2 User
	if err := DB.Where("name = ?", "create_from_map_with_schema2").First(&result2).Error; err != nil || result2.Age != 1 {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	_, ok = mapValue2["id"].(uint)
	if ok {
		t.Skipf("This test case skipped, because the db supports returning")
	}

	idVal, ok = mapValue2["id"].(int64)
	if !ok {
		t.Fatal("ret result missing id")
	}

	if int64(result2.ID) != idVal {
		t.Fatal("failed to create data from map with table, @id != id")
	}

	// case 3: records
	values := []map[string]interface{}{
		{"name": "create_from_map_with_schema11", "age": 1}, {"name": "create_from_map_with_schema12", "age": 1},
	}

	beforeLen := len(values)
	if err := DB.Model(&User{}).Create(&values).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}

	// mariadb with returning, values will be appended with id map
	if len(values) == beforeLen*2 {
		t.Skipf("This test case skipped, because the db supports returning")
	}

	for i := range values {
		v, ok := values[i]["id"]
		if !ok {
			t.Fatal("failed to create data from map with table, returning map has no primary key")
		}

		var result User
		if err := DB.Where("name = ?", fmt.Sprintf("create_from_map_with_schema1%d", i+1)).First(&result).Error; err != nil || result.Age != 1 {
			t.Fatalf("failed to create from map, got error %v", err)
		}
		if int64(result.ID) != v.(int64) {
			t.Fatal("failed to create data from map with table, @id != id")
		}
	}
}

func TestCreateFromMapWithTable(t *testing.T) {
	tableDB := DB.Table("users")
	supportLastInsertID := isMysql() || isSqlite()

	// case 1: create from map[string]interface{}
	record := map[string]interface{}{"name": "create_from_map_with_table", "age": 18}
	if err := tableDB.Create(record).Error; err != nil {
		t.Fatalf("failed to create data from map with table, got error: %v", err)
	}

	if _, ok := record["@id"]; !ok && supportLastInsertID {
		t.Fatal("failed to create data from map with table, returning map has no key '@id'")
	}

	var res map[string]interface{}
	if err := tableDB.Select([]string{"id", "name", "age"}).Where("name = ?", "create_from_map_with_table").Find(&res).Error; err != nil || res["age"] != int64(18) {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	if _, ok := record["@id"]; ok && fmt.Sprint(res["id"]) != fmt.Sprint(record["@id"]) {
		t.Fatalf("failed to create data from map with table, @id != id, got %v, expect %v", res["id"], record["@id"])
	}

	// case 2: create from *map[string]interface{}
	record1 := map[string]interface{}{"name": "create_from_map_with_table_1", "age": 18}
	tableDB2 := DB.Table("users")
	if err := tableDB2.Create(&record1).Error; err != nil {
		t.Fatalf("failed to create data from map, got error: %v", err)
	}
	if _, ok := record1["@id"]; !ok && supportLastInsertID {
		t.Fatal("failed to create data from map with table, returning map has no key '@id'")
	}

	var res1 map[string]interface{}
	if err := tableDB2.Select([]string{"id", "name", "age"}).Where("name = ?", "create_from_map_with_table_1").Find(&res1).Error; err != nil || res1["age"] != int64(18) {
		t.Fatalf("failed to create from map, got error %v", err)
	}

	if _, ok := record1["@id"]; ok && fmt.Sprint(res1["id"]) != fmt.Sprint(record1["@id"]) {
		t.Fatal("failed to create data from map with table, @id != id")
	}

	// case 3: create from []map[string]interface{}
	records := []map[string]interface{}{
		{"name": "create_from_map_with_table_2", "age": 19},
		{"name": "create_from_map_with_table_3", "age": 20},
	}

	tableDB = DB.Table("users")
	if err := tableDB.Create(&records).Error; err != nil {
		t.Fatalf("failed to create data from slice of map, got error: %v", err)
	}

	if _, ok := records[0]["@id"]; !ok && supportLastInsertID {
		t.Fatal("failed to create data from map with table, returning map has no key '@id'")
	}

	if _, ok := records[1]["@id"]; !ok && supportLastInsertID {
		t.Fatal("failed to create data from map with table, returning map has no key '@id'")
	}

	var res2 map[string]interface{}
	if err := tableDB.Select([]string{"id", "name", "age"}).Where("name = ?", "create_from_map_with_table_2").Find(&res2).Error; err != nil || res2["age"] != int64(19) {
		t.Fatalf("failed to query data after create from slice of map, got error %v", err)
	}

	var res3 map[string]interface{}
	if err := DB.Table("users").Select([]string{"id", "name", "age"}).Where("name = ?", "create_from_map_with_table_3").Find(&res3).Error; err != nil || res3["age"] != int64(20) {
		t.Fatalf("failed to query data after create from slice of map, got error %v", err)
	}

	if _, ok := records[0]["@id"]; ok && fmt.Sprint(res2["id"]) != fmt.Sprint(records[0]["@id"]) {
		t.Errorf("failed to create data from map with table, @id != id, got %v, expect %v", res2["id"], records[0]["@id"])
	}

	if _, ok := records[1]["id"]; ok && fmt.Sprint(res3["id"]) != fmt.Sprint(records[1]["@id"]) {
		t.Errorf("failed to create data from map with table, @id != id")
	}
}
