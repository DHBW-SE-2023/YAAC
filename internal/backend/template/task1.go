package yaac_backend_template

func (b *BackendTemplate) Foo() {
	Bar()
}

// "public" function
func Bar() {
	bar()
}

/*
Packages in the backend should be grouped by purpose, however, those can be
seperated into multiple files if needed.
*/
