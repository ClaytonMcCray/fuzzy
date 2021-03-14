package generate

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	functionsToTest := []singleTest{
		{
			FunctionUnderTest: "MyFunction",
			CompositePtrsToInit: []argPair{
				{Type: "myType", Name: "one"},
			},

			SimpleAutosToInit: []argPair{
				{Type: "int", Name: "two"},
				{Type: "float64", Name: "three"},
			},

			FunctionArgumentNamesInOrder: []string{"one", "two", "three"},
			ReceiverName:                 argPair{Type: "myRecvrType", Name: "recvr"},
			ReceiverIsPointer:            true,
			ReceiverIsComposite:          true,
		},
		{
			FunctionUnderTest: "PlainFunction",
			SimplePtrsToInit: []argPair{
				{Type: "int", Name: "one"},
				{Type: "float64", Name: "two"},
			},
			CompositeAutosToInit: []argPair{
				{Type: "CompAuto", Name: "three"},
				{Type: "bytes.Buffer", Name: "four"},
			},

			FunctionArgumentNamesInOrder: []string{"one", "two", "four", "three"},
		},
	}

	in := metaData{
		PackageName:         "mypackage",
		CompletePackagePath: "gitter.com/me/mypackage",
		TestsToGenerate:     functionsToTest,
	}

	b, err := in.Do()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}
