package fail

import "github.com/mephistolie/chefbook-backend-common/responses/fail"

var (
	GrpcShoppingListNotFound = fail.CreateGrpcServer(fail.TypeNotFound, "shopping list not found")
)
