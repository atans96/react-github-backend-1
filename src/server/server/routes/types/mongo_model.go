package types

import (
	"time"
)

type languagePreference struct {
	Language string `json:"language"`
	Checked  bool   `json:"checked"`
}
type renderImages struct {
	ID    int      `json:"id"`
	Value []string `json:"value"`
}
type starRanking struct {
	Id     int `json:"id"`
	Trends struct {
		Daily     int `json:"daily"`
		Weekly    int `json:"weekly"`
		Monthly   int `json:"monthly"`
		Quarterly int `json:"quarterly"`
		Yearly    int `json:"yearly"`
	} `json:"trends"`
	TimeSeries struct {
		Daily   []int `json:"daily"`
		Monthly struct {
			Year     int `json:"year"`
			Months   int `json:"months"`
			FirstDay int `json:"firstDay"`
			LastDay  int `json:"lastDay"`
			Delta    int `json:"delta"`
		}
	} `json:"timeSeries"`
}
type repoInfo struct {
	FullName      string   `json:"full_name"`
	Description   string   `json:"description"`
	Stars         int      `json:"stars"`
	Forks         int      `json:"forks"`
	UpdatedAt     string   `json:"updatedAt"`
	Language      string   `json:"language"`
	Topics        []string `json:"topics"`
	DefaultBranch string   `json:"default_branch"`
	HtmlUrl       string   `json:"html_url"`
	Readme        string   `json:"readme"`
}
type repoContributions struct {
	FullName     int `json:"full_name"`
	Contributors []struct {
		Login         string `json:"login"`
		AvatarUrl     string `json:"avatar_url"`
		Contributions int    `json:"contributions"`
	} `json:"contributors"`
}
type seenCards struct {
	StargazersCount int    `json:"stargazers_count"`
	FullName        string `json:"full_name"`
	Owner           struct {
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
	} `json:"owner"`
	Description   string   `json:"description"`
	Language      string   `json:"language"`
	Topics        []string `json:"topics"`
	HtmlUrl       []string `json:"html_url"`
	Name          string   `json:"name"`
	Id            int      `json:"id"`
	DefaultBranch string   `json:"default_branch"`
	IsQueried     bool     `json:"is_queried"`
}
type searches struct {
	Search    string `json:"search"`
	Count     int    `json:"count"`
	UpdatedAt string `json:"updatedAt"`
}
type repoInfoSuggested struct {
	From            string `json:"from"`
	IsSeen          bool   `json:"is_seen"`
	StargazersCount int    `json:"stargazers_count"`
	FullName        string `json:"full_name"`
	DefaultBranch   string `json:"default_branch"`
	Owner           struct {
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
	} `json:"owner"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
	HtmlUrl     string   `json:"html_url"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
}
type clicked struct {
	IsQueried   bool      `json:"is_queried"`
	FullName    string    `json:"full_name"`
	Count       int       `json:"count"`
	DateClicked time.Time `json:"dateClicked"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}
type starred struct {
	FullName  string `json:"full_name"`
	IsQueried bool   `json:"is_queried"`
}
type User struct {
	UserName           string               `json:"userName"`
	Avatar             string               `json:"avatar"`
	JoinDate           time.Time            `json:"joinDate"`
	Token              string               `json:"token"`
	LanguagePreference []languagePreference `json:"languagePreference"`
	Starred            []starred            `json:"starred"`
	RenderImages       []renderImages
	Languages          []string            `json:"languages"`
	RepoContributions  []repoContributions `json:"repoContributions"`
	RepoInfo           []repoInfo          `json:"repoInfo"`
	StarRanking        []starRanking       `json:"starRanking"`
	SeenCards          []seenCards         `json:"seenCards"`
	Searches           []searches          `json:"searches"`
	RepoInfoSuggested  []repoInfoSuggested `json:"repoInfoSuggested"`
	Clicked            []clicked           `json:"clicked"`
	Rss                []string            `json:"rss"`
	LastSeen           []string            `json:"lastSeen"`
}
