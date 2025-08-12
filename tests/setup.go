package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Ulpio/reservas-cipt/database"
	"github.com/Ulpio/reservas-cipt/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	os.Setenv("JWT_SECRET", "testsecret")

	dsn := fmt.Sprintf("file:%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	database.DB = db

	return db
}
