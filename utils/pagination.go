package utils

import (
	"strconv"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"gorm.io/gorm"
)

type paginate struct {
   limit int
   page  int
}

func newPaginate(limit int, page int) *paginate {
   return &paginate{limit: limit, page: page}
}

func (p *paginate) paginatedResult(db *gorm.DB) *gorm.DB {
   offset := (p.page - 1) * p.limit

   return db.Offset(offset).
      Limit(p.limit)
}

func GetPaginatedQuery[T interface{}](model T, pageInput string, limitInput string) (*gorm.DB, int,int) {
	page, err := strconv.Atoi(pageInput)
    if err != nil || page < 1 {
        page = 1
    }
    limit, err := strconv.Atoi(limitInput)
    if err != nil || limit < 1 {
        limit = 10
    }
	return database.Database.Db.Model(model).Scopes(newPaginate(limit, page).paginatedResult), page, limit
}