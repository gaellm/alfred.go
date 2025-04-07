package function

import (
	"alfred/internal/helper"
	"alfred/internal/mock"
	"alfred/pkg/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlfredFunc_MissingFunction(t *testing.T) {
	// Mock JavaScript content without an alfred function
	jsContent := `
		function test() {}
	`

	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call AlfredFunc
	req := request.Req{}
	res := request.Res{}
	mock := mock.Mock{}
	_, err := f.AlfredFunc(mock, nil, req, res)
	assert.False(t, f.HasFuncAlfred)
	assert.Error(t, err)
	assert.Equal(t, "function file test.js not contains alfred function", err.Error())
}

func TestSetupFunc_MissingSetup(t *testing.T) {
	// Mock JavaScript content without a setup function
	jsContent := `
		function test2() {}
	`
	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call SetupFunc
	err := f.SetupFunc()
	//assert.Error(t, err)
	assert.Equal(t, "function file test.js not contains setup function", err.Error())
}

func TestUpdateHelpersListener_MissingFunction(t *testing.T) {
	// Mock JavaScript content without an updateHelpers function
	jsContent := `
		function test3() {}
	`

	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call UpdateHelpersListener
	helpers := []helper.Helper{{Name: "existingHelper"}}
	_, err := f.UpdateHelpersListener(helpers)
	assert.Error(t, err)
	assert.Equal(t, "function file test.js not contains updateHelpers function", err.Error())
}

func TestCreateFunction(t *testing.T) {
	// Mock JavaScript content
	jsContent := `
		function setup() {}
		function alfred() {}
		function updateHelpers() {}
	`

	// Create a function
	f, err := CreateFunction("test.js", []byte(jsContent))
	assert.NoError(t, err)
	assert.Equal(t, "test.js", f.FileName)
	assert.True(t, f.HasFuncSetup)
	assert.True(t, f.HasFuncAlfred)
	assert.True(t, f.HasFuncUpdateHelpers)
}

func TestSetupFunc(t *testing.T) {
	// Mock JavaScript content with a setup function
	jsContent := `
		function setup() {
			console.log("Setup completed");
		}
	`

	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call SetupFunc
	err := f.SetupFunc()
	assert.NoError(t, err)
}

func TestUpdateHelpersListener(t *testing.T) {
	// Mock JavaScript content with an updateHelpers function
	jsContent := `
		function updateHelpers(helpers) {
			helpers.push({ name: "newHelper" });
			return helpers;
		}
	`

	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call UpdateHelpersListener
	helpers := []helper.Helper{{Name: "existingHelper"}}
	updatedHelpers, err := f.UpdateHelpersListener(helpers)
	assert.NoError(t, err)
	assert.Len(t, updatedHelpers, 2)
	assert.Equal(t, "newHelper", updatedHelpers[1].Name)
}

func TestAlfredFunc(t *testing.T) {
	// Mock JavaScript content with an alfred function
	jsContent := `
		function alfred(mock, helpers, req, res) {
			res.body = "Hello, Alfred!";
			return res;
		}
	`

	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Call AlfredFunc
	req := request.Req{}
	res := request.Res{}
	mock := mock.Mock{}
	updatedRes, err := f.AlfredFunc(mock, []helper.Helper{}, req, res)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Alfred!", updatedRes.Body)

}

func TestCheckIfFuncExists(t *testing.T) {
	// Mock JavaScript content with a setup function
	jsContent := `
		function setup() {}
	`
	// Create a function
	f, _ := CreateFunction("test.js", []byte(jsContent))

	// Check if setup function exists
	exists, err := f.CheckIfFuncExists("setup")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check if a non-existent function exists
	exists, err = f.CheckIfFuncExists("nonExistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}
