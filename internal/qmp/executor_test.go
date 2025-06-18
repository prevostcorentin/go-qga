package qmp_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp"
)

const expectedCommandResponse = "777 holy trinity"

type fakeConnection struct{}

func (_ *fakeConnection) Connect(_ string) *QmpConnectionError {
	return nil
}

func (_ *fakeConnection) Send(_ []byte) ([]byte, *QmpConnectionError) {
	response := struct {
		Return *testCommandResponse `json:"return"`
	}{Return: &testCommandResponse{Value: expectedCommandResponse}}
	responseBytes, _ := json.Marshal(response)
	return responseBytes, nil
}

func (_ *fakeConnection) Close() error {
	return nil
}

type testCommand struct{}
type testCommandArguments struct {
	Argument int `json:"argument"`
}
type testCommandResponse struct {
	Value string `json:"value"`
}

func (_ *testCommand) Execute() string {
	return "test-command"
}

func (_ *testCommand) Arguments() any {
	return &testCommandArguments{}
}

func (_ *testCommand) Response() any {
	return &testCommandResponse{}
}

func TestRunCommand(t *testing.T) {
	executor := qmp.NewExecutor(&fakeConnection{})
	response, err := executor.Run(&testCommand{})
	if err != nil {
		t.Fatalf("while running command: %v", err)
	}
	if response.(*testCommandResponse).Value != expectedCommandResponse {
		t.Errorf(`wrong value "%v" for response. expected "777"`, response.(*testCommandResponse).Value)
	}
}

type testCommandRecursive struct{}
type testCommandArgumentsRecursive struct {
	Self *testCommandArgumentsRecursive `json:"self"`
}

func (_ *testCommandRecursive) Execute() string {
	return "test-command"
}

func (_ *testCommandRecursive) Arguments() any {
	arguments := &testCommandArgumentsRecursive{}
	arguments.Self = arguments
	return arguments
}

func (_ *testCommandRecursive) Response() any {
	return &testCommandResponse{}
}

func TestMarshalFailure(t *testing.T) {
	executor := qmp.NewExecutor(&fakeConnection{})
	failingCommand := &testCommandRecursive{}
	_, err := executor.Run(failingCommand)
	if err == nil {
		t.Fatal("should have raised an error")
	}
	if err.Domain() != CodecDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), CodecDomain)
	}
}

type unmarshableResponseConnection struct{}

func (_ *unmarshableResponseConnection) Connect(_ string) *QmpConnectionError {
	return nil
}

func (_ *unmarshableResponseConnection) Send(_ []byte) ([]byte, *QmpConnectionError) {
	return []byte("i am no object"), nil
}

func (_ *unmarshableResponseConnection) Close() error {
	return nil
}

func TestUnmarshalResponseTopLevelFailure(t *testing.T) {
	executor := qmp.NewExecutor(&unmarshableResponseConnection{})
	failingCommand := &testCommand{}
	var err QgaError
	if _, err = executor.Run(failingCommand); err == nil {
		t.Errorf("there should have been an error here")
	}
	if err.Domain() != CodecDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), CodecDomain)
	}
	if err.Kind() != string(Unmarshal) {
		t.Errorf(`wrong error kind "%v". expected "%s"`, err.Kind(), Unmarshal)
	}
}

type childUnmarshalFailureCommand struct{}

func (_ *childUnmarshalFailureCommand) Execute() string {
	return "unmarshal-failure"
}

func (_ *childUnmarshalFailureCommand) Arguments() any {
	return &testCommandArguments{}
}

func (_ *childUnmarshalFailureCommand) Response() any {
	return &childUnmarshableResponse{}
}

type childUnmarshableResponse struct {
	Value int `json:"value"`
}

func TestUnmarshalResponseReturnsFailure(t *testing.T) {
	executor := qmp.NewExecutor(&fakeConnection{})
	failingCommand := &childUnmarshalFailureCommand{}
	var err QgaError
	if _, err = executor.Run(failingCommand); err == nil {
		t.Errorf("there should have been an error here")
	}
	if err.Domain() != CodecDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), CodecDomain)
	}
	if err.Kind() != string(Unmarshal) {
		t.Errorf(`wrong error kind "%v". expected "%s"`, err.Kind(), Unmarshal)
	}
}

type sendFailureConnection struct{}

func (_ *sendFailureConnection) Connect(_ string) *QmpConnectionError {
	return nil
}

func (_ *sendFailureConnection) Send(_ []byte) ([]byte, *QmpConnectionError) {
	return nil, NewQmpConnectionError(fmt.Errorf("I am designed to fail"), SendErrorKind)
}

func (_ *sendFailureConnection) Close() error {
	return nil
}

func TestConnectionSendFailure(t *testing.T) {
	executor := qmp.NewExecutor(&sendFailureConnection{})
	failingCommand := &testCommand{}
	var err QgaError
	if _, err = executor.Run(failingCommand); err == nil {
		t.Errorf("there should have been an error here")
	}
	if err.Domain() != QmpConnectionDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), QmpConnectionDomain)
	}
	if err.Kind() != string(SendErrorKind) {
		t.Errorf(`wrong error kind "%v". expected "%s"`, err.Kind(), SendErrorKind)
	}
}

type missingReturnsConnection struct{}

func (_ *missingReturnsConnection) Connect(_ string) *QmpConnectionError {
	return nil
}

func (_ *missingReturnsConnection) Send(_ []byte) ([]byte, *QmpConnectionError) {
	response := struct {
		MissingReturn *testCommandResponse `json:"missing_return"`
	}{MissingReturn: &testCommandResponse{Value: expectedCommandResponse}}
	responseBytes, _ := json.Marshal(response)
	return responseBytes, nil
}

func (_ *missingReturnsConnection) Close() error {
	return nil
}

func TestNoReturnsFailure(t *testing.T) {
	executor := qmp.NewExecutor(&missingReturnsConnection{})
	failingCommand := &testCommand{}
	var err QgaError
	if _, err = executor.Run(failingCommand); err == nil {
		t.Errorf("there should have been an error here")
	}
	if err.Domain() != CodecDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), CodecDomain)
	}
	if err.Kind() != string(Key) {
		t.Errorf(`wrong error kind "%v". expected "%s"`, err.Kind(), Unmarshal)
	}
}
