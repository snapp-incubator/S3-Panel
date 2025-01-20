package main

import (
	_ "gitlab.snapp.ir/platform/snapp_object_store/docs"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

//	@title			ObjectStorage Backend Swagger
//	@version		1.0
//	@description	Serves the S3 backend APIs
//	@termsOfService	https://swagger.io/terms/

//	@contact.name	Cloud-Platform
//	@contact.url	https://docs.snappcloud.io/docs/servicedesk
//	@contact.email	cloud-platform@snappcloud.io

// @license.name	Apache 2.0
// @license.url	https://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	cmd.Execute()
}
