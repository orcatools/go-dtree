package dtree

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var treeTest = []byte(`[
	{
		"id": 1,
		"name": "isTest"
	},
	{
		"id": 2,
		"name": "count",
		"parent_id": 1,
		"operator": "eq",
		"value": true
	},
	{
		"id": 3,
		"name": "Never Reach",
		"parent_id": 1,
		"operator": "eq",
		"value": false
	},
	{
		"id": 4,
		"name": "FinalNode 2",
		"parent_id": 2,
		"operator": "gt",
		"value": 10,
		"order":1
	},
	{
		"id": 5,
		"name": "FinalNode 1",
		"parent_id": 2,
		"operator": "lt",
		"value": 10,
		"order":2
	},
	{
		"id": 6,
		"name": "FinalNode 3",
		"parent_id": 2,
		"value": "fallback"
	}
]`)

func TestTree_SimpleTest(t *testing.T) {
	// Arrange

	//Load Tree
	tr, err := LoadTree(treeTest)
	if err != nil {
		t.Fail()
	}

	// Load request
	jsonRequest := []byte(`{
		"isTest":  true,
		"count":   15
	}`)

	//Act
	result, err := tr.ResolveJSON(jsonRequest)

	//Assert
	assert.NoError(t, err, "Resolve should not have errors")
	assert.Equal(t, "FinalNode 2", result.Name)
}

func TestTree_SimpleTest_With_Error_Config(t *testing.T) {
	// Arrange

	//Load Tree
	tr, err := LoadTree(treeTest)
	if err != nil {
		t.Fail()
	}

	// Load request
	jsonRequest := []byte(`{
		"isTest":  true,
		"count":   "15"
	}`)

	f := func(t *TreeOptions) {
		t.StopIfConvertingError = true
	}

	//Act
	result, err := tr.ResolveJSON(jsonRequest, f)

	//Assert
	assert.Error(t, err, "Resolve should not return an error when the type of the request is the not the same as the one defined on tree")
	assert.Equal(t, "count", result.Name)
}

func TestTree_SimpleTest_Without_Error_Config(t *testing.T) {
	// Arrange

	//Load Tree
	tr, err := LoadTree(treeTest)
	if err != nil {
		t.Fail()
	}

	// Load request
	jsonRequest := []byte(`{
		"isTest":  true,
		"count":   "15"
	}`)

	f := func(t *TreeOptions) {
		t.StopIfConvertingError = false
	}

	//Act
	result, err := tr.ResolveJSON(jsonRequest, f)

	//Assert
	assert.NoError(t, err, "Resolve should not return an error even if the type of the request is the not the same as the one defined on tree")
	assert.Equal(t, "FinalNode 3", result.Name)
}

func TestTree_SimpleTest_With_Bad_Json(t *testing.T) {
	//Act
	_, err := LoadTree([]byte("not a json"))

	//Assert
	assert.Error(t, err, "LoadTree should return an error if the json is malformed")
}

func TestTree_SimpleTest_Resolving_Bad_Json(t *testing.T) {
	// Arrange

	//Load Tree
	tr, err := LoadTree(treeTest)
	if err != nil {
		t.Fail()
	}

	// Load request
	jsonRequest := []byte("Obviously not a json")

	//Act
	_, err = tr.ResolveJSON(jsonRequest)

	//Assert
	assert.Error(t, err, "Resolve should return an error  if the jsonrequest is malformed")
}

func ExampleLoadTree() {
	jsonTree := []byte(`[
		{
			"id": 1,
			"name": "sayHello"
		},
		{
			"id": 2,
			"name": "GoodBye",
			"parent_id": 1,
			"operator": "eq",
			"value": false
		},
		{
			"id": 3,
			"name": "gender",
			"parent_id": 1,
			"operator": "eq",
			"value": true
		},
		{
			"id": 4,
			"name": "Hello Miss",
			"parent_id": 3,
			"operator": "eq",
			"value": "F"
		},
		{
			"id": 5,
			"name": "Hello",
			"parent_id": 3,
			"value": "fallback"
		},
		{
			"id": 6,
			"name": "age",
			"parent_id": 3,
			"operator": "eq",
			"value": "M"
		},
		{
			"id": 7,
			"name": "Hello Sir",
			"parent_id": 6,
			"operator": "gt",
			"value": 60
		},
		{
			"id": 8,
			"name": "Hello dude",
			"parent_id": 6,
			"operator": "lte",
			"value": 60
		}
	]`)

	t, err := LoadTree(jsonTree)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	request := make(map[string]interface{})
	request["sayHello"] = true
	request["gender"] = "M"
	request["age"] = 35.0 //does not use int, the engine only support float (if you want  do a PR to include int, it's up to you)

	/*request := []byte(`{
	  		"sayHello": false,
	  		"gender":   "M",
	  		"age": 35
	  	}`)

	  v, _ := t.ResolveJSON(request)
	*/

	v, _ := t.Resolve(request)

	fmt.Println(v.Name)
	// output: Hello dude
}
