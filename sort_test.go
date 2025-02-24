package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils/tests"
)

func TestSortScope(t *testing.T) {
	db, _ := gorm.Open(&tests.DummyDialector{}, nil)
	sort := &Sort{Field: "notacolumn", Order: SortAscending}
	modelIdentity := &modelIdentity{
		Columns: map[string]*column{
			"name": {Name: "Name"},
		},
		TableName: "test_models",
	}

	assert.Nil(t, sort.Scope(&Settings{}, modelIdentity))

	sort.Field = "name"

	db = db.Scopes(sort.Scope(&Settings{}, modelIdentity)).Table("table").Find(nil)
	expected := map[string]clause.Clause{
		"ORDER BY": {
			Name: "ORDER BY",
			Expression: clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column: clause.Column{
							Table: "test_models",
							Name:  "name",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expected, db.Statement.Clauses)

	sort.Order = SortDescending
	db, _ = gorm.Open(&tests.DummyDialector{}, nil)
	db = db.Scopes(sort.Scope(&Settings{}, modelIdentity)).Table("table").Find(nil)
	expected["ORDER BY"].Expression.(clause.OrderBy).Columns[0].Desc = true
	assert.Equal(t, expected, db.Statement.Clauses)
}

func TestSortScopeBlacklisted(t *testing.T) {
	sort := &Sort{Field: "name", Order: SortAscending}
	modelIdentity := &modelIdentity{
		Columns: map[string]*column{
			"name": {Name: "Name"},
		},
		TableName: "test_models",
	}
	assert.Nil(t, sort.Scope(&Settings{Blacklist: Blacklist{FieldsBlacklist: []string{"name"}}}, modelIdentity))
}
