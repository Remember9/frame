package gorm

import (
	"fmt"
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/util/xcast"
	"github.com/Remember9/frame/xlog"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"sync"
	"testing"
	"time"
)

type DictTest struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Money     float64    `json:"money"`
	Status    int8       `json:"status"`
	Type      int8       `json:"type"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (p *DictTest) TableName() string {
	return "dict_test"
}

func init() {
	err := config.InitTest()
	if err != nil {
		panic(err)
	}
	err = xlog.Build()
	if err != nil {
		panic(err)
	}
}

func TestBuild(t *testing.T) {
	db := Build("mysql")

	var dictTest []DictTest
	err := db.Clauses(dbresolver.Read).
		Select("name", "status").
		Where("name = ?", "test2").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "status"}, Desc: true}).
		// Order("status DESC").
		Limit(2).
		Find(&dictTest).Error
	if err != nil {
		t.Error(err)
	}
	t.Logf("show info %+v", dictTest)
}

func TestOnConflict(t *testing.T) {
	db := Build("mysql")

	var dictTest = []DictTest{{ID: 3, Name: "test11"}, {ID: 2, Name: "test12", Status: 1}, {ID: 1, Name: "test13"}}

	// err := db.Create(&dictTest).Error

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&dictTest).Error

	if err != nil {
		t.Error(err)
	}

	t.Logf("show info %+v", dictTest)

}

func TestOnConflict2(t *testing.T) {
	db := Build("mysql")

	var dictTest = []DictTest{{Money: 245, Name: "test11", Type: 2, Status: 1}}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "money"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(&dictTest).Error

	if err != nil {
		t.Error(err)
	}

	t.Logf("show info %+v", dictTest)

}

func TestExpr(t *testing.T) {
	db := Build("mysql")

	var dictTest DictTest
	err := db.Model(&dictTest).
		Where("name = ?", "test13").
		Update("money", gorm.Expr("money * ? + ?", 2, 100)).Error
	if err != nil {
		t.Error(err)
	}
	t.Log("update success")
}

func TestTx(t *testing.T) {
	db := Build("mysql")

	// rand.Seed(time.Now().UnixNano())

	s := time.Now()
	num := 10

	wg := &sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(t int) {
			defer wg.Done()
			// err := CreateAnimals(db, "test"+cast.ToString(t), cast.ToFloat64(rand.Intn(1<<10)))
			// err := CreateAnimals2(db, "test"+xcast.ToString(t), xcast.ToFloat64(rand.Intn(1<<10)))

			err := SelectAnimals3(db, "test"+xcast.ToString(t))
			if err != nil {
				fmt.Println("err-", t, err)
			}
		}(i)
	}
	wg.Wait()

	// wg := &sync.WaitGroup{}
	// wg.Add(num)
	// for i := 0; i < num; i++ {
	//	go func(n int) {
	//		defer wg.Done()
	//		ss := time.Now()
	//		xlog.Debug("test concurrent", xlog.Duration("cost", time.Since(ss)), xlog.Int("No.", n))
	//	}(i)
	// }
	// wg.Wait()

	// for i := 0; i < num; i++ {
	//	err := SelectAnimals3(db, "test"+xcast.ToString(t))
	//	if err != nil {
	//		fmt.Println("err-", t, err)
	//	}
	// }

	t.Log(time.Since(s))
}

func CreateAnimals(db *gorm.DB, name string, money float64) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&DictTest{Name: name, Money: money}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func CreateAnimals2(db *gorm.DB, name string, money float64) error {
	return db.Create(&DictTest{Name: name, Money: money}).Error
}

func SelectAnimals3(db *gorm.DB, name string) (err error) {
	var dictTest []DictTest
	ctx := context.Background()
	err = db.WithContext(ctx).Clauses(dbresolver.Read).
		Select("name", "status").
		Where("name = ?", name).
		Limit(2).
		Find(&dictTest).Error
	return err
}
