# cmdb-sdk-golang
golang sdk of cmdb operations

## install
```shell
go get "github.com/veops/cmdb-sdk-golang"
```

## Operation of CI
```golang
// create a helper firstly
// urlPrefix is used to combine http request eg. the final query url with urlPrefix https://demo.veops.cn/api/v0.1 is https://demo.veops.cn/api/v0.1/ci/s
// key and secret is obtained from ACL
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#%E5%9B%9Bapi%E9%89%B4%E6%9D%83%E6%96%B9%E6%B3%95
helper := cmdb_sdk.NewHelper("urlprefix", "your key", "your secret")

// suppose you have created a ci type called mycitype
// with three attributes server_name, ip and custom_attr
// you can add a ci instance as following
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#2-%E6%96%B0%E5%A2%9Eci%E6%8E%A5%E5%8F%A3
attrs := map[string]any{
  "server_name": "test-1",
  "ip":          "192.168.0.1",
  "custom_attr": 123,
}
addCIRes, err := helper.AddCI("mycitype", cmdb_sdk.NoAttrPolicyDefault, cmdb_sdk.ExistPolicyDefault, attrs)

// you are able to get the ci instance created above by query now
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#1-ci%E6%9F%A5%E8%AF%A2%E6%8E%A5%E5%8F%A3
getCIRes, err := helper.GetCI("_type:mycitype", "", "", "", 0, 0, cmdb_sdk.RetKeyDefault)

// now let's update some attribute of a ci instance
// reminder that you can get ci_id from the return result of AddCI or instance info of GetCI
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#3-%E4%BF%AE%E6%94%B9ci%E6%8E%A5%E5%8F%A3
updates := map[string]any{
  "ip":          "127.0.0.1",
  "custom_attr": 321,
}
updateCIRes, err := helper.UpdateCI(666, "mycitype", cmdb_sdk.NoAttrPolicyDefault, updates)

// finally, delete this ci if you want
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#4-%E5%88%A0%E9%99%A4ci%E6%8E%A5%E5%8F%A3
deleteCIRes, err := helper.DeleteCI(666)
```

## Operation of CI Relation
```golang
// create a helper firstly
helper := cmdb_sdk.NewHelper("urlprefix", "your key", "your secret")

// assuming you now have create a relation between two ci types
// and two ci instances with ci_id=666 and 777 sperately of these types
// you can add a relation to these instances
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#2-%E5%A2%9E%E5%8A%A0ci%E5%85%B3%E7%B3%BB%E6%8E%A5%E5%8F%A3
AddCIRRes, err := helper.AddRelation(666, 777)

// you will get a cr_id after the creation which can be used to delete the relation
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#3-%E5%88%A0%E9%99%A4ci%E5%85%B3%E7%B3%BB%E6%8E%A5%E5%8F%A3
deleteCIRRes, err := helper.DeleteRelation(666777, 0, 0)
// in the case that you don't keep the ci_id, you can delete relation by two ci_id
deleteCIRRes, err := helper.DeleteRelation(0, 666, 777)

// except for root_id,reverse and level, relation query is much similar to ci query
//
// doc https://github.com/veops/cmdb/blob/master/docs/cmdb_api.md#1-ci%E5%85%B3%E7%B3%BB%E6%9F%A5%E8%AF%A2%E6%8E%A5%E5%8F%A3
getCIRRes, err := helper.GetRelation(666, 0, "1", "", "", "", "", 0, 0, cmdb_sdk.RetKeyDefault)
```
