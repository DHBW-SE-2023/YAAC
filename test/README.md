# tests

## setup tests
- create for each go file a new testfile called *testfile*_test.go
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


# fyne tests




## run tests
navigate to the test directory
```shell
go test
```