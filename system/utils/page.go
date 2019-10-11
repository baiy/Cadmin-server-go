package utils

import (
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Page struct {
	Offset   int `form:"offset"`
	PageSize int `form:"pageSize"`
}

func (p Page) Select(table string, lists interface{}, expressions ...exp.Expression) (total int, err error) {
	count, err := models.Db.From(table).Where(expressions...).Count()
	if err != nil {
		return
	}
	total = int(count)
	err = models.Db.From(table).Where(expressions...).
		Offset(uint(p.Offset)).Limit(uint(p.PageSize)).
		Order(goqu.I("id").Desc()).ScanStructs(lists)
	return
}
