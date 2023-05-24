package model

import (
	"time"
)

type GitHubUser struct {
	Id           uint32    `json:"id" gorm:"primaryKey"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	NodeId       string    `json:"node_id"`
	AvatarURL    string    `json:"avatar_url"`
	GravatarId   string    `json:"gravatar_id"`
	Url          string    `json:"url"`
	HtmlUrl      string    `json:"html_url"`
	FollowersUrl string    `json:"followers_url"`
	Type         string    `json:"type"`
	SiteAdmin    bool      `json:"site_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Sample JSON
// {
// 	"login":"",
// 	"id":0,
// 	"node_id":"",
// 	"avatar_url":"",
// 	"gravatar_id":"",
// 	"url":"",
// 	"html_url":"",
// 	"followers_url":"",
// 	"following_url":"",
// 	"gists_url":"",
// 	"starred_url":"",
// 	"subscriptions_url":"",
// 	"organizations_url":"",
// 	"repos_url":"",
// 	"events_url":"",
// 	"received_events_url":"",
// 	"type":"",
// 	"site_admin":false,
// 	"name":null,
// 	"company":null,
// 	"blog":"",
// 	"location":null,
// 	"email":null,
// 	"hireable":null,
// 	"bio":null,
// 	"twitter_username":null,
// 	"public_repos":0,
// 	"public_gists":0,
// 	"followers":0,
// 	"following":0,
// 	"created_at":"0000-00-00T00:00:00Z",
// 	"updated_at":"0000-00-00T00:00:00Z"
// }
