package models

import "gorm.io/gorm"

type Repo struct {
	gorm.Model
	GitHubID      int64  `gorm:"uniqueIndex;not null"`
	Owner         string `gorm:"index:idx_repo_owner_name,unique;not null"`
	Name          string `gorm:"index:idx_repo_owner_name,unique;not null"`
	Description   string `gorm:"type:text"`
	Private       bool   `gorm:"not null;default:false"`
	DefaultBranch string
	Status        string `gorm:"index;not null;default:pending"`
	ImportedBy    uint
}
