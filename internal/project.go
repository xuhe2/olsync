package internal

import (
	"encoding/json"
	"time"
)

// User 定义了JSON中嵌套的用户信息结构。
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Project 定义了项目数组中每个项目的结构。
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Archived      bool      `json:"archived"`
	Trashed       bool      `json:"trashed"`
	AccessLevel   string    `json:"accessLevel"`
	Source        string    `json:"source"`
	LastUpdated   time.Time `json:"lastUpdated"`
	LastUpdatedBy User      `json:"lastUpdatedBy"`
	Owner         User      `json:"owner"`
}

// PrefetchedProjectsBlob 定义了从 meta 标签中解析出的 JSON 顶级结构。
type PrefetchedProjectsBlob struct {
	TotalSize int       `json:"totalSize"`
	Projects  []Project `json:"projects"`
}

func ParsePrefetchedProjectsBlob(jsonBlob string) (PrefetchedProjectsBlob, error) {
	var blob PrefetchedProjectsBlob
	err := json.Unmarshal([]byte(jsonBlob), &blob)
	return blob, err
}
