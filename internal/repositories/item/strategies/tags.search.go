package repository

import (
	"strings"

	"gorm.io/gorm"
)

type TagSearch struct {
    Tags []string
}

func (s TagSearch) Apply(db *gorm.DB) *gorm.DB {
    if len(s.Tags) == 0 {
        return db
    }
	for i, tag := range s.Tags {
    	s.Tags[i] = strings.ToLower(tag)
	}
    return db.Joins("JOIN item_tags it ON items.item_id::uuid = it.item_id").
        Joins("JOIN tags t ON it.tag_id::uuid = t.tag_id").
        Where("LOWER(t.tag_name) IN ?", s.Tags).
        Group("items.item_id").
        Having("COUNT(DISTINCT LOWER(t.tag_name)) = ?", len(s.Tags)).
        Select("items.*")
}