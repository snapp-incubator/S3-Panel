package main

import (
	_ "go.uber.org/automaxprocs"

	_ "github.com/snapp-incubator/S3-Panel/docs"
	"github.com/snapp-incubator/S3-Panel/internal/app"
)

//	@title			S3 Panel Backend
//	@version		1.0
//	@description	Serves the S3 object storage panel backend APIs
//	@termsOfService	https://swagger.io/terms/

//	@contact.name	snapp-incubator
//	@contact.url	https://github.com/snapp-incubator/S3-Panel

// @license.name	GPL-3.0
// @license.url	https://www.gnu.org/licenses/gpl-3.0.html
func main() {
	app.Execute()
}
