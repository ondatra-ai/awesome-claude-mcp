package github

type commentNode struct {
	Path     *string `json:"path"`
	Line     *int    `json:"line"`
	Body     string  `json:"body"`
	Outdated bool    `json:"outdated"`
	URL      string  `json:"url"`
}

type pageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}

type threadsPageResponse struct {
	Data struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					PageInfo pageInfo `json:"pageInfo"`
					Nodes    []struct {
						ID         string `json:"id"`
						IsResolved bool   `json:"isResolved"`
						Comments   struct {
							Nodes []commentNode `json:"nodes"`
						} `json:"comments"`
					} `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}
