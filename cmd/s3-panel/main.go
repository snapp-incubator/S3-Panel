package main

import (
	_ "github.com/snapp-incubator/S3-Panel/docs"
	"github.com/snapp-incubator/S3-Panel/internal/app"
	_ "go.uber.org/automaxprocs"
)

//	@title			S3 Panel Backend
//	@version		1.0
//	@description	Serves the S3 object storage panel backend APIs
//	@termsOfService	https://swagger.io/terms/

//	@contact.name	Cloud-Platform
//	@contact.url	https://docs.snappcloud.io/docs/servicedesk
//	@contact.email	cloud-platform@snappcloud.io

// @license.name	Apache 2.0
// @license.url	https://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	app.Execute()
}
