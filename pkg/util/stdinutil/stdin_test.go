package stdinutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/juju/errors"
)

func TestGetFromStdin_NoChangeQuestion_Default(t *testing.T) {

	params := GetFromStdin_params{
		Question: "Is this a test?",
		DefaultValue: "Yes",
		ValidationRegexPattern: "No|Yes",
	}

	err := mockStdin("invalid\ninvalid\n")
	if err != nil {
		t.Error(errors.ErrorStack(err))
	}
	defer cleanUpMockedStdin()

	answer := GetFromStdin(&params)

	if answer != params.DefaultValue {
		t.Error("Wrong Answer.\nExpected default answer: " + params.DefaultValue + "\nBut Got: " + answer)
	}
}

func TestGetFromStdin_NoChangeQuestion_NonDefault(t *testing.T) {

	params := GetFromStdin_params{
		Question: "Is this a test?",
		DefaultValue: "No",
		ValidationRegexPattern: "No|Yes",
	}

	err := mockStdin("invalid\nYes\n")
	if err != nil {
		t.Error(errors.ErrorStack(err))
	}
	defer cleanUpMockedStdin()

	answer := GetFromStdin(&params)

	if answer != "Yes" {
		t.Error("Wrong Answer.\nExpected: Yes\nBut Got: " + answer)
	}
}

func TestGetFromStdin_ChangeQuestion(t *testing.T) {

	params := GetFromStdin_params{
		Question: "Is this a test?",
		DefaultValue: "No",
		ValidationRegexPattern: "No|Yes",
		InputTerminationString: " ",
	}

	err := mockStdin("invalid\nYes\n")
	if err != nil {
		t.Error(errors.ErrorStack(err))
	}
	defer cleanUpMockedStdin()

	answer := GetFromStdin(&params)

	if answer != "Yes" {
		t.Error("Wrong Answer.\nExpected: Yes\nBut Got: " + answer)
	}
}

var tmpfile *os.File
var oldStdin *os.File

func mockStdin(inputString string) error{
	//Code from https://stackoverflow.com/a/46365584 (modified)
	input := []byte(inputString)
	var err error
	tmpfile, err = ioutil.TempFile("", "testGetFromStdin")
    if err != nil {
        return errors.Trace(err)
    }

    if _, err := tmpfile.Write(input); err != nil {
        return errors.Trace(err)
    }

    if _, err := tmpfile.Seek(0, 0); err != nil {
        return errors.Trace(err)
    }

    oldStdin = os.Stdin
	os.Stdin = tmpfile

	return nil
}

func cleanUpMockedStdin() {
	os.Remove(tmpfile.Name())
	os.Stdin = oldStdin
}
