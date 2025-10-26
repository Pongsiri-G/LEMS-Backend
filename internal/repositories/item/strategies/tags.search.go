package repository

import "gorm.io/gorm"

type TagSearch struct {
    Tags []string
}

func (s TagSearch) Apply(db *gorm.DB) *gorm.DB {
    if len(s.Tags) == 0 {
        return db
    }
    return db.Joins("JOIN item_tags it ON items.item_id = it.item_id").
        Joins("JOIN tags t ON it.tag_id = t.tag_id").
        Where("t.tag_name IN ?", s.Tags).
        Group("items.item_id").
        Having("COUNT(DISTINCT t.tag_name) = ?", len(s.Tags)).
        Select("items.*")
}