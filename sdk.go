package cmdb_sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// Helper
type Helper struct {
	urlPrefix string
	key       string
	secret    string
	c         *resty.Client
}

// NewHelper creates a helper instance for cmdb operation
//
// urlPrefix is used to combine http request eg. the final query url with urlPrefix https://demo.veops.cn/api/v0.1 is https://demo.veops.cn/api/v0.1/ci/s
func NewHelper(urlPrefix, key, secret string) *Helper {
	urlPrefix = strings.TrimRight(urlPrefix, "/")
	return &Helper{
		urlPrefix: urlPrefix,
		key:       key,
		secret:    secret,
		c:         resty.New(),
	}
}

// buildAPIKey
func (h *Helper) buildAPIKey(u string, params map[string]any) string {
	pu, _ := url.Parse(u)
	keys := lo.Keys(params)
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	vals := strings.Join(
		lo.Map(keys, func(k string, _ int) string {
			return lo.Ternary(strings.HasPrefix(k, "-"), "", cast.ToString(params[k]))
		}),
		"")
	sha := sha1.New()
	sha.Write([]byte(strings.Join([]string{pu.Path, h.secret, vals}, "")))
	return hex.EncodeToString(sha.Sum(nil))
}

// checkRespError
func (h *Helper) checkRespError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		re := &ResponseError{}
		json.Unmarshal(resp.Body(), re)
		return fmt.Errorf("httpCode=%d message=%s", resp.StatusCode(), re.Message)
	}
	return nil
}

// AddCI adds a CI model
//
// you need to add ci_types firstly
func (h *Helper) AddCI(ciType string, noAttrPolicy NoAttrPolicy, existPolicy ExistPolicy, attrs map[string]any) (res *AddCIResult, err error) {
	u := fmt.Sprintf("%s/ci", h.urlPrefix)
	fmt.Println(u)
	payload := make(map[string]any)
	copier.Copy(&payload, attrs)
	payload["ci_type"] = ciType
	payload["no_attrbute_policy"] = noAttrPolicy
	payload["exist_policy"] = existPolicy
	payload["_secret"] = h.buildAPIKey(u, payload)
	payload["_key"] = h.key

	res = &AddCIResult{}
	resp, err := h.c.R().
		SetBody(payload).
		SetResult(res).
		Post(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// DeleteCI deletes a CI by id
func (h *Helper) DeleteCI(ciID int) (res *DeleteCIResult, err error) {
	u := fmt.Sprintf("%s/ci/%d", h.urlPrefix, ciID)
	payload := make(map[string]any)
	payload["_secret"] = h.buildAPIKey(u, payload)
	payload["_key"] = h.key

	res = &DeleteCIResult{}
	resp, err := h.c.R().
		SetBody(payload).
		SetResult(res).
		Delete(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// UpdateCI updates a CI model
func (h *Helper) UpdateCI(ciID int, ciType string, noAttrPolicy NoAttrPolicy, attrs map[string]any) (res *UpdateCIResult, err error) {
	u := fmt.Sprintf("%s/ci/%d", h.urlPrefix, ciID)
	payload := make(map[string]any)
	copier.Copy(&payload, attrs)
	payload["ci_type"] = ciType
	payload["_secret"] = h.buildAPIKey(u, payload)
	payload["_key"] = h.key

	res = &UpdateCIResult{}
	resp, err := h.c.R().
		SetBody(payload).
		SetResult(res).
		Put(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// GetCI queries CIs
func (h *Helper) GetCI(q, fl, facet, sort string, page, count int, retKey RetKey) (res *GetCIResult, err error) {
	params := map[string]any{
		"q":       q,
		"fl":      fl,
		"sort":    sort,
		"page":    page,
		"count":   count,
		"ret_key": retKey,
	}
	u := fmt.Sprintf("%s/ci/s", h.urlPrefix)
	params["_secret"] = h.buildAPIKey(u, params)
	params["_key"] = h.key
	ps := make(map[string]string)
	for k, v := range params {
		ps[k] = cast.ToString(v)
	}

	res = &GetCIResult{}
	resp, err := h.c.R().
		SetQueryParams(ps).
		SetResult(res).
		Get(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// AddRelation adds a relation between two CIs
//
// you need to add ci_relations firstly
func (h *Helper) AddRelation(srcCIID, dstCIID int) (res *AddRelationResult, err error) {
	u := fmt.Sprintf("%s/ci_relations/%d/%d", h.urlPrefix, srcCIID, dstCIID)
	payload := make(map[string]any)
	payload["_secret"] = h.buildAPIKey(u, payload)
	payload["_key"] = h.key

	res = &AddRelationResult{}
	resp, err := h.c.R().
		SetBody(payload).
		SetResult(res).
		Post(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// DeleteRelation deletes a CI relation
//
// use relationID or firstCIID and secondCIID
// firstCIID and secondCIID will be used if ciID is 0
func (h *Helper) DeleteRelation(relationID, firstCIID, secondCIID int) (res *DeleteRelationResult, err error) {
	u := lo.Ternary(relationID != 0,
		fmt.Sprintf("%s/ci_relations/%d", h.urlPrefix, relationID),
		fmt.Sprintf("%s/ci_relations/%d/%d", h.urlPrefix, firstCIID, secondCIID),
	)
	payload := make(map[string]any)
	payload["_secret"] = h.buildAPIKey(u, payload)
	payload["_key"] = h.key

	res = &DeleteRelationResult{}
	resp, err := h.c.R().
		SetBody(payload).
		SetResult(res).
		Delete(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}

// GetRelation queries the CI relation
func (h *Helper) GetRelation(rootId, reverse int, level, q, fl, facet, sort string, page, count int, retKey RetKey) (res *GetRelationResult, err error) {
	u := fmt.Sprintf("%s/ci_relations/s", h.urlPrefix)
	params := map[string]any{
		"root_id": fmt.Sprint(rootId),
		"level":   level,
		"reverse": fmt.Sprint(reverse),
		"q":       q,
		"fl":      fl,
		"sort":    sort,
		"page":    fmt.Sprint(page),
		"count":   fmt.Sprint(count),
		"ret_key": string(retKey),
	}
	params["_secret"] = h.buildAPIKey(u, params)
	params["_key"] = h.key
	ps := make(map[string]string)
	for k, v := range params {
		ps[k] = cast.ToString(v)
	}

	res = &GetRelationResult{}
	resp, err := h.c.R().
		SetQueryParams(ps).
		SetResult(res).
		Get(u)
	if err = h.checkRespError(resp, err); err != nil {
		return nil, err
	}

	return
}
