package cmdb_sdk

type RetKey = string

const (
	RetKeyDefault = RetKey("")
	RetKeyID      = RetKey("id")
	RetKeyAlias   = RetKey("alias")
	RetKeyName    = RetKey("name")
)

type NoAttrPolicy = string

const (
	NoAttrPolicyDefault = NoAttrPolicy("")
	NoAttrPolicyReject  = NoAttrPolicy("reject")
)

type ExistPolicy = string

const (
	ExistPolicyDefault = ExistPolicy("")
	ExistPolicyNeed    = ExistPolicy("need")
	ExistPolicyReject  = ExistPolicy("reject")
	ExistPolicyReplace = ExistPolicy("replace")
)

type ResponseError struct {
	Message string `json:"message"`
}

type AddCIResult struct {
	CIID int `json:"ci_id"`
}

type DeleteCIResult struct {
	Message string `json:"message"`
}

type UpdateCIResult struct {
	CIID int `json:"ci_id"`
}

type GetCIResult struct {
	Counter  map[string]int   `json:"counter"`
	Facet    map[string]any   `json:"facet"`
	Numfound int              `json:"numfound"`
	Page     int              `json:"page"`
	Result   []map[string]any `json:"result"`
	Total    int              `json:"total"`
}

type AddRelationResult struct {
	RelationID int `json:"cr_id"`
}

type DeleteRelationResult struct {
	Message string `json:"message"`
}

type UpdateRelationResult struct {
	RelationID int `json:"cr_id"`
}

type GetRelationResult struct {
	Counter  map[string]int   `json:"counter"`
	Facet    map[string]any   `json:"facet"`
	Numfound int              `json:"numfound"`
	Page     int              `json:"page"`
	Result   []map[string]any `json:"result"`
	Total    int              `json:"total"`
}
