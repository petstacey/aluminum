package data

type JobTitle struct {
	ID    int64  `json:"jobTitleId"`
	Title string `json:"title"`
}

type Workgroup struct {
	ID   int64  `json:"workgroupId"`
	Name string `json:"name"`
}

type Location struct {
	ID   int64  `json:"locationId"`
	Name string `json:"name"`
}

type Resource struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
	JobTitle  string `json:"jobTitle"`
	Workgroup string `json:"workgroup"`
	Location  string `json:"location"`
	Manager   string `json:"manager"`
	Active    bool   `json:"active"`
}
