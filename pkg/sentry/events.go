package sentry

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type SentryEvent struct {
	ID              string                 `json:"eventID"`
	Title           string                 `json:"title"`
	Count           int64                  `json:"count()"`
	EventsPerMinute float64                `json:"epm()"`
	Level           string                 `json:"level"`
	EventType       string                 `json:"event.type"`
	Platform        string                 `json:"platform"`
	IDHex           string                 `json:"id.hex"`
	ProjectId       int64                  `json:"issue.project_id"`
	GroupID         string                 `json:"groupID"`
	Timestamp       time.Time              `json:"timestamp"`
	Received        time.Time              `json:"received"`
	Dist            string                 `json:"dist"`
	Transaction     string                 `json:"transaction"`
	DataModules     map[string]interface{} `json:"data.modules"`
	GetTypeDisplay  string                 `json:"get_type_display"`
	Message         string                 `json:"message"`
	Metadata        map[string]interface{} `json:"metadata"`
	Tags            []interface{}          `json:"tags"`
	Entries         []Entry                `json:"entries"`
	DataContexts    map[string]interface{} `json:"data.contexts"`
	DataExtra       map[string]interface{} `json:"data.extra"`
	DataUser        string                 `json:"data.user"`
}

type Entry struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type GetEventsInput struct {
	OrganizationSlug string
	ProjectIds       []string
	Environments     []string
	Query            string
	From             time.Time
	To               time.Time
	Sort             string
	Limit            int64
}

// getRequiredFields returns the list of fields that are required to be fetched
// from the sentry API. This is used to build the query string.
func getRequiredFields() []string {
	return []string{
		"id",
		"title",
		"project",
		"project.id",
		"release",
		"count()",
		"epm()",
		"last_seen()",
		"level",
		"event.type",
		"platform",
	}
}

func (gei *GetEventsInput) ToQuery() string {
	urlPath := fmt.Sprintf("/api/0/projects/%s/%s/events/?", gei.OrganizationSlug, gei.ProjectIds[0])
	if gei.Limit < 1 || gei.Limit > 100 {
		gei.Limit = 100
	}
	params := url.Values{}
	params.Set("query", gei.Query)
	params.Set("start", gei.From.Format("2006-01-02T15:04:05"))
	params.Set("end", gei.To.Format("2006-01-02T15:04:05"))
	if gei.Sort != "" {
		params.Set("sort", gei.Sort)
	}
	params.Set("per_page", strconv.FormatInt(gei.Limit, 10))
	for _, field := range getRequiredFields() {
		params.Add("field", field)
	}
	for _, projectId := range gei.ProjectIds {
		params.Add("project", projectId)
	}
	for _, environment := range gei.Environments {
		params.Add("environment", environment)
	}
	return urlPath + params.Encode()
}

func (sc *SentryClient) GetEvents(gei GetEventsInput) ([]SentryEvent, string, error) {
	var out []SentryEvent
	executedQueryString := gei.ToQuery()
	err := sc.Fetch(executedQueryString, &out)
	return out, sc.BaseURL + executedQueryString, err
}
