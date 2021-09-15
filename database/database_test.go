package database

import (
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"github.com/jinzhu/copier"
	"testing"
)

func databaseTestSetup() {
	setupTestEventVariables()
}

func setupDatabaseMock(t *testing.T) Repository {
	repository, err := GetMockRepository()
	if err != nil {
		t.Error(err)
	}
	InitMockRepository(repository)
	return repository
}

func TestRepository_GetEvent(t *testing.T) {
	databaseTestSetup()
	var wantedEvents types.Events
	_ = copier.Copy(&wantedEvents, &events)
	wantedEvents[0].ID = 2
	wantedEvents[1].ID = 3
	tests := []struct {
		name    string
		eventId int
		want    *types.Event
		wantErr bool
	}{
		{"unknown id", 4, nil, false},
		{"totalEvent", 1, nil, false}, // Do not return TotalSeptimana event
		{"e0", 2, &wantedEvents[0], false},
		{"ev1", 3, &wantedEvents[1], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			got, err := rep.GetEvent(tt.eventId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualEvent(got, tt.want) {
				t.Errorf("\nGetEvent() got = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetEvents(t *testing.T) {
	databaseTestSetup()
	tests := []struct {
		name    string
		year    *int
		want    []types.Event
		wantErr bool
	}{
		{"all", nil, types.Events{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := setupDatabaseMock(t)
			got, err := rep.GetEvents(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !util.EqualEvents(got, tt.want) {
				t.Errorf("GetEvents() got = %v,\n want %v", got, tt.want)
			}
		})
	}
}

//
//func TestRepository_GetLocation(t *testing.T) {
//	type fields struct {
//		Db *gorm.DB
//	}
//	type args struct {
//		id string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *types.Location
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rep := Repository{
//				Db: tt.fields.Db,
//			}
//			got, err := rep.GetLocation(tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetLocation() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetLocation() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepository_GetLocations(t *testing.T) {
//	type fields struct {
//		Db *gorm.DB
//	}
//	type args struct {
//		overallLocation *types.OverallLocation
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []types.Location
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rep := Repository{
//				Db: tt.fields.Db,
//			}
//			got, err := rep.GetLocations(tt.args.overallLocation)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetLocations() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetLocations() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepository_HasApiKeyInfo(t *testing.T) {
//	type fields struct {
//		Db *gorm.DB
//	}
//	type args struct {
//		info types.ApiKeyInfo
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    bool
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rep := Repository{
//				Db: tt.fields.Db,
//			}
//			got, err := rep.HasApiKeyInfo(tt.args.info)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("HasApiKeyInfo() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("HasApiKeyInfo() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestRepository_InitDatabase(t *testing.T) {
//	type fields struct {
//		Db *gorm.DB
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rep := Repository{
//				Db: tt.fields.Db,
//			}
//		})
//	}
//}
//
//func TestRepository_StoreSecurityInfo(t *testing.T) {
//	type fields struct {
//		Db *gorm.DB
//	}
//	type args struct {
//		info types.ApiKeyInfo
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			rep := Repository{
//				Db: tt.fields.Db,
//			}
//		})
//	}
//}
//
//func Test_insertStartEnd(t *testing.T) {
//	type args struct {
//		db *gorm.DB
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//		})
//	}
//}
//
//func Test_openDB(t *testing.T) {
//	tests := []struct {
//		name    string
//		want    *gorm.DB
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := openDB()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("openDB() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("openDB() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
