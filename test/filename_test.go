package test

import (
	"fmt"
	"testing"
)

func TestFunctionName(t *testing.T) {
	/*
		test function
		//if error occurrs, t.Error("Error Message") or t.Errorf("Error with value: %s", value)

		//simple Example
		a := math.getSquare(9)
		if a != 3 {
			t.Errorf("Got: %d; want 3 ", response)
		}
	*/
}

func ExampleTest() { // Test if the following Output is printed
	fmt.Println("hello, and") // Print statements actually in the function to be tested
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}

/*
func TestThirdWindowUI(t *testing.T) {
	out_label, in_entry := makeUI_w3()

	// initial state
	if out_label.Text != "Hello World!" {
		t.Error("Incorrect initial value")
	}

	// user input
	test.Type(in_entry, "Hi")
	if out_label.Text != "Hi to you as well!" {
		t.Error("Incorrect user input handling")
	}
}

*/
