package fail

import "github.com/mephistolie/chefbook-backend-common/responses/fail"

var (
	typeOutdatedVersion = "outdated_version"
)

var (
	GrpcOutdatedVersion      = fail.CreateGrpcConflict(typeOutdatedVersion, "shopping list version is outdated; process current version first")
	GrpcShoppingListNotFound = fail.CreateGrpcServer(fail.TypeNotFound, "shopping list not found")
)
