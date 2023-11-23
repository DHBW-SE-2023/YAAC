# tests

## setup tests
- create for each go file a new testfile called *testfile*_test.go in the same dictonary where the go file is
- create a function called Test*FunctionNameToTest*(t *testing.T)
- modul *testing* should be imported automaticaly
- test the function


## creating tests
# basic tests
- check if errors occur and handel it: 
> t.Error("Error Message!")
> t.Errorf("Got: %d; want 3 ", response)

# output tests
- call the function to test with the output "Hello World!"
> // Output:
> // Hello World

For more information about tests in go, see also https://pkg.go.dev/testing


# fyne tests
- test the UI with user input
- you can test the UI without 
- For more information abour fyne tests, see https://developer.fyne.io/started/testing




## run tests
navigate to the test directory
```shell
go test
```

Go tests are also executet by using make