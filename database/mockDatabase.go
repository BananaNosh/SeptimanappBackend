package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
)

func GetMockRepository() (Repository, error) {
	//, &gorm.Config{
	//	DisableForeignKeyConstraintWhenMigrating: true,
	//	//Logger: newLogger,
	//}
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	return Repository{Db: db}, err
}

func InitMockRepository(repository Repository) {
	dataPath := "../data/testData"
	repository.InitDatabaseFromPath(path.Join(dataPath, "horariumHelper"))
	repository.InitDatabaseFromPath(path.Join(dataPath, "locationHelper"))
}
