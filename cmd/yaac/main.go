package yaac

import yaac_mvvm "github.com/DHBW-SE-2023/YAAC/internal/mvvm"

func Run() {
	mvvm := yaac_mvvm.New()
	mvvm.StartApplication()
}
