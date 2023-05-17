package repository

//func TestRulesFromFileRepo_GetAllRules_Success(t *testing.T) {
//	repo := rulesFromFileRepo{}
//
//	// Create a temporary file for the test data
//	file, err := os.CreateTemp("", "tmpfile-")
//	assert.NoError(t, err)
//	defer func(name string) {
//		err := os.Remove(name)
//		if err != nil {
//
//		}
//	}(file.Name())
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//
//		}
//	}(file)
//
//	// Write some test data to the temporary file
//	testData := []byte(`{"condition":"AND","zk_request_type":{"id":"zk_req_type","field":"zk_req_type","type":"string","input":"string","operator":"equal","value":"HTTP"},"rules":[{"id":"req_method","field":"req_method","type":"string","input":"string","operator":"equal","value":"*"},{"id":"resp_status","field":"resp_status","type":"integer","input":"integer","operator":"notequal","value":200}],"valid":true}
//{"condition":"AND","zk_request_type":{"id":"zk_req_type","field":"zk_req_type","type":"string","input":"string","operator":"equal","value":"HTTP"},"rules":[{"id":"category","field":"category","type":"integer","input":"select","operator":"equal","value":1},{"id":"name","field":"name","type":"string","input":"text","operator":"not_in","value":"a,b"},{"condition":"AND","rules":[{"id":"category","field":"category","type":"integer","input":"select","operator":"equal","value":1}]}],"valid":true}`)
//	_, err = file.Write(testData)
//	assert.NoError(t, err)
//
//	// Override the file path in the repo with the temporary file path
//	oldFilePath := filePath
//
//	filePath = file.Name()
//	defer func() { filePath = oldFilePath }()
//
//	// Call GetAllRules and check the returned data
//	rules, err := repo.GetAllRules(nil)
//	assert.Len(t, rules, 2)
//
//	rule1 := rules[0]
//	assert.Equal(t, nil, err)
//	assert.Equal(t, "AND", rule1["condition"])
//	assert.NotNil(t, rule1["zk_request_type"])
//	assert.NotNil(t, rule1["rules"])
//
//	rule2 := rules[1]
//	assert.Equal(t, "AND", rule2["condition"])
//	assert.NotNil(t, rule2["zk_request_type"])
//	assert.NotNil(t, rule2["rules"])
//}
//
//func TestRulesFromFileRepo_GetAllRules_InvalidJson_Failure(t *testing.T) {
//	repo := rulesFromFileRepo{}
//
//	// Create a temporary file for the test data
//	file, err := os.CreateTemp("", "tmpfile-")
//	assert.NoError(t, err)
//	defer func(name string) {
//		err := os.Remove(name)
//		if err != nil {
//
//		}
//	}(file.Name())
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//
//		}
//	}(file)
//
//	// Write some test data to the temporary file
//	testData := []byte(`{"condition":"AND,"zk_request_type":{"id":"zk_req_type","field":"zk_req_type","type":"string","input":"string","operator":"equal","value":"HTTP"},"rules":[{"id":"category","field":"category","type":"integer","input":"select","operator":"equal","value":1},{"id":"name","field":"name","type":"string","input":"text","operator":"not_in","value":"a,b"},{"condition":"AND","rules":[{"id":"category","field":"category","type":"integer","input":"select","operator":"equal","value":1}]}],"valid":true}`)
//	_, err = file.Write(testData)
//	assert.NoError(t, err)
//
//	// Override the file path in the repo with the temporary file path
//	oldFilePath := filePath
//
//	filePath = file.Name()
//	defer func() { filePath = oldFilePath }()
//
//	// Call GetAllRules and check the returned data
//	rules, err := repo.GetAllRules(nil)
//
//	assert.NotNil(t, err)
//	assert.Equal(t, "invalid character 'z' after object key:value pair", err.Error())
//	assert.Len(t, rules, 0)
//}
//
//func TestRulesFromFileRepo_GetAllRules_FileNotFound_Failure(t *testing.T) {
//	repo := rulesFromFileRepo{}
//
//	// Create a temporary file for the test data
//	file, err := os.CreateTemp("", "tmpfile-")
//	assert.NoError(t, err)
//	err = os.Remove(file.Name())
//	if err != nil {
//		return
//	}
//
//	oldFilePath := filePath
//
//	filePath = file.Name()
//	defer func() { filePath = oldFilePath }()
//
//	// Call GetAllRules and check the returned data
//	rules, err := repo.GetAllRules(nil)
//	assert.Len(t, rules, 0)
//	assert.Equal(t, "cannot access given file", err.Error())
//}
//
//func TestNewRulesFromFileRepo(t *testing.T) {
//	x := NewRulesFromFileRepo()
//	assert.NotNil(t, x)
//}
