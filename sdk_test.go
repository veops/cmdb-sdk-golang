package cmdb_sdk

import (
	"reflect"
	"testing"
)

const (
	urlPrefix = "http://YOUR HOST/api/v0.1"
	key       = "YOUR KEY"
	secret    = "YOUR SECRET"
)

var (
	helper = NewHelper(urlPrefix, key, secret)
)

func TestHelper_buildAPIKey(t *testing.T) {
	type args struct {
		u      string
		params map[string]any
	}
	tests := []struct {
		name string
		h    *Helper
		args args
		want string
	}{
		{
			name: "Test-1",
			h:    helper,
			args: args{
				u:      urlPrefix + "/ci/s",
				params: map[string]any{},
			},
			want: "a9b7bfe568d3ffcb335528650b0c77a54f9c0414",
		},
		{
			name: "Test-2",
			h:    helper,
			args: args{
				u: urlPrefix + "/ci/s",
				params: map[string]any{
					"ci_type": "server",
					"custom":  123,
				},
			},
			want: "769aa8ae3687b7930a84f24a9c0698b6657dd774",
		},
		{
			name: "Test-3",
			h:    helper,
			args: args{
				u: urlPrefix + "/ci/s",
				params: map[string]any{
					"ci_type":             "server",
					"no_attribute_policy": NoAttrPolicyReject,
				},
			},
			want: "9d3ba60b3ca4019808e63d8e7892d32002c285d7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.buildAPIKey(tt.args.u, tt.args.params); got != tt.want {
				t.Errorf("Helper.buildAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelper_AddCI(t *testing.T) {
	type args struct {
		ciType       string
		noAttrPolicy NoAttrPolicy
		existPolicy  ExistPolicy
		attrs        map[string]any
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *AddCIResult
		wantErr bool
	}{
		{
			name: "Test-1",
			h:    helper,
			args: args{
				ciType:       "mycitype",
				noAttrPolicy: NoAttrPolicyDefault,
				existPolicy:  ExistPolicyReject,
				attrs: map[string]any{
					"server_name": "test-1",
					"ip":          "192.168.0.1",
					"custom_attr": 123,
				},
			},
			wantRes: &AddCIResult{},
			wantErr: false,
		},
		{
			name: "Test-2",
			h:    helper,
			args: args{
				ciType:       "mycitype",
				noAttrPolicy: NoAttrPolicyDefault,
				existPolicy:  ExistPolicyReject,
				attrs: map[string]any{
					"server_name": "test-2",
					"ip":          "172.0.0.1",
					"custom_attr": 999,
				},
			},
			wantRes: &AddCIResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.AddCI(tt.args.ciType, tt.args.noAttrPolicy, tt.args.existPolicy, tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.AddCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes.CIID <= 0 {
				t.Errorf("Helper.AddCI() = %v, want gotRes.CIID > 0", gotRes)
			}
		})
	}
}

func TestHelper_DeleteCI(t *testing.T) {
	type args struct {
		ciID int
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *DeleteCIResult
		wantErr bool
	}{
		{
			name: "Test-1",
			h:    helper,
			args: args{
				ciID: 9723,
			},
			wantRes: &DeleteCIResult{
				Message: "ok",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.DeleteCI(tt.args.ciID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.DeleteCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Helper.DeleteCI() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestHelper_UpdateCI(t *testing.T) {
	type args struct {
		ciID         int
		ciType       string
		noAttrPolicy NoAttrPolicy
		attrs        map[string]any
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *UpdateCIResult
		wantErr bool
	}{
		{
			name: "Test-1",
			h:    helper,
			args: args{
				ciID:         9723,
				ciType:       "mycitype",
				noAttrPolicy: NoAttrPolicyDefault,
				attrs: map[string]any{
					"ip":          "192.168.0.1",
					"custom_attr": 123,
				},
			},
			wantRes: &UpdateCIResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.UpdateCI(tt.args.ciID, tt.args.ciType, tt.args.noAttrPolicy, tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.UpdateCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes.CIID <= 0 {
				t.Errorf("Helper.UpdateCI() = %v, want gotRes.CIID > 0", gotRes)
			}
		})
	}
}

func TestHelper_GetCI(t *testing.T) {
	type args struct {
		q      string
		fl     string
		facet  string
		sort   string
		page   int
		count  int
		retKey RetKey
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *GetCIResult
		wantErr bool
	}{
		{
			name: "",
			h:    helper,
			args: args{
				q:      "_type:mycitype",
				fl:     "",
				facet:  "",
				sort:   "",
				page:   0,
				count:  0,
				retKey: key,
			},
			wantRes: &GetCIResult{
				Counter: map[string]int{
					"mycitype": 1,
				},
				Facet:    map[string]any{},
				Numfound: 1,
				Page:     1,
				Result: []map[string]any{
					{
						"_id":           float64(9723),
						"_type":         float64(165),
						"ci_type":       "mycitype",
						"ci_type_alias": "mycitype",
						"custom_attr":   float64(123),
						"ip":            "192.168.0.1",
						"server_name":   "test-1",
						"unique":        "server_name",
						"unique_alias":  "Server名",
					},
				},
				Total: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.GetCI(tt.args.q, tt.args.fl, tt.args.facet, tt.args.sort, tt.args.page, tt.args.count, tt.args.retKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.GetCI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Helper.GetCI() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestHelper_AddRelation(t *testing.T) {
	type args struct {
		srcCIID int
		dstCIID int
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *AddRelationResult
		wantErr bool
	}{
		{
			name: "",
			h:    helper,
			args: args{
				srcCIID: 9723,
				dstCIID: 9727,
			},
			wantRes: &AddRelationResult{
				RelationID: 978,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.AddRelation(tt.args.srcCIID, tt.args.dstCIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.AddRelation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Helper.AddRelation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestHelper_DeleteRelation(t *testing.T) {
	type args struct {
		relationID int
		firstCIID  int
		secondCIID int
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *DeleteRelationResult
		wantErr bool
	}{
		{
			name: "Test-1",
			h:    helper,
			args: args{
				relationID: 979,
			},
			wantRes: &DeleteRelationResult{
				Message: "CIType relation deleted",
			},
			wantErr: false,
		},
		{
			name: "Test-2",
			h:    helper,
			args: args{
				firstCIID:  9723,
				secondCIID: 9727,
			},
			wantRes: &DeleteRelationResult{
				Message: "CIType relation deleted",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.DeleteRelation(tt.args.relationID, tt.args.firstCIID, tt.args.secondCIID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.DeleteRelation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Helper.DeleteRelation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestHelper_GetRelation(t *testing.T) {
	type args struct {
		rootId  int
		reverse int
		level   string
		q       string
		fl      string
		facet   string
		sort    string
		page    int
		count   int
		retKey  RetKey
	}
	tests := []struct {
		name    string
		h       *Helper
		args    args
		wantRes *GetRelationResult
		wantErr bool
	}{
		{
			name: "",
			h:    helper,
			args: args{
				rootId:  9723,
				reverse: 0,
				level:   "1",
				q:       "",
				fl:      "",
				facet:   "",
				sort:    "",
				page:    0,
				count:   0,
				retKey:  RetKeyDefault,
			},
			wantRes: &GetRelationResult{
				Counter: map[string]int{
					"mycitype": 1,
				},
				Facet:    map[string]any{},
				Numfound: 1,
				Page:     1,
				Result: []map[string]any{
					{
						"_id":           float64(9727),
						"_type":         float64(165),
						"ci_type":       "mycitype",
						"ci_type_alias": "mycitype",
						"custom_attr":   float64(999),
						"ip":            "172.0.0.1",
						"server_name":   "test-2",
						"unique":        "server_name",
						"unique_alias":  "Server名",
					},
				},
				Total: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.h.GetRelation(tt.args.rootId, tt.args.reverse, tt.args.level, tt.args.q, tt.args.fl, tt.args.facet, tt.args.sort, tt.args.page, tt.args.count, tt.args.retKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Helper.GetRelation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Helper.GetRelation() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
