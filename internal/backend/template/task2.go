package yaac_backend_template

import "fmt"

func (b *BackendTemplate) BackTemp() {
	b.MVVM.TemplateFuncRecive(true)
}

// "Prvate" function
func bar() {
	fmt.Println("Hello from the other side")
}
