package sentry

import (
	"encoding/json"
	"errors"
	"time"
)

type SentryProject struct {
	DateCreated  time.Time `json:"dateCreated"`
	HasAccess    bool      `json:"hasAccess"`
	ID           string    `json:"id"`
	IsBookmarked bool      `json:"isBookmarked"`
	IsMember     bool      `json:"isMember"`
	Environments []string  `json:"environments"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Team         struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"team"`
	Teams []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"teams"`
}

// 為 struct 內的 number 字段添加自定義反序列化
func (p *SentryProject) UnmarshalJSON(data []byte) error {
	// Create alias to avoid recursion
	type TeamAlias struct {
		ID   json.Number `json:"id"`
		Name string      `json:"name"`
		Slug string      `json:"slug"`
	}

	type Alias SentryProject
	aux := &struct {
		ID    json.Number `json:"id"`
		Team  TeamAlias   `json:"team"`
		Teams []TeamAlias `json:"teams"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert json.Number to string
	p.ID = aux.ID.String()
	p.Team.ID = aux.Team.ID.String()

	// Initialize Teams slice and convert IDs
	p.Teams = make([]struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}, len(aux.Teams))

	for i, team := range aux.Teams {
		p.Teams[i].ID = team.ID.String()
		p.Teams[i].Name = team.Name
		p.Teams[i].Slug = team.Slug
	}

	return nil
}

func (sc *SentryClient) GetProjects(organizationSlug string, withPagination bool) ([]SentryProject, error) {
	projects := []SentryProject{}
	if organizationSlug == "" {
		organizationSlug = sc.OrgSlug
	}
	url := "/api/0/organizations/" + organizationSlug + "/projects/"

	if withPagination {
		for url != "" {
			batch := []SentryProject{}
			nextURL, err := sc.FetchWithPagination(url, &batch)
			if err != nil {
				return nil, err
			}

			projects = append(projects, batch...)
			url = nextURL
		}
		return projects, nil
	} else {
		err := sc.Fetch(url, &projects)
		return projects, err
	}
}

func (sc *SentryClient) GetTeamsProjects(organizationSlug string, teamSlug string) ([]SentryProject, error) {
	out := []SentryProject{}
	if organizationSlug == "" {
		organizationSlug = sc.OrgSlug
	}
	if teamSlug == "" {
		return out, errors.New("invalid team slug")
	}
	err := sc.Fetch("/api/0/teams/"+organizationSlug+"/"+teamSlug+"/projects/", &out)
	return out, err
}
