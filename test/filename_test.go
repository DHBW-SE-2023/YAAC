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

func ExamplTest() { // Test if the following Output is printed
	fmt.Println("hello, and") // Print statements actually in the function to be tested
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}
