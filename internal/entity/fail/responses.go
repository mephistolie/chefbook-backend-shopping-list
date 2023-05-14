package fail

import (
	"fmt"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
)

var (
	typeOutdatedVersion       = "outdated_version"
	typeMaxShoppingListsCount = "max_shopping_lists_count"
)

var (
	GrpcOutdatedVersion      = fail.CreateGrpcConflict(typeOutdatedVersion, "shopping list version is outdated; process current version first")
	GrpcShoppingListNotFound = fail.CreateGrpcServer(fail.TypeNotFound, "shopping list not found")
)

func GrpcMaxShoppingListsCount(count int) error {
	return fail.CreateGrpcAccessDenied(typeMaxShoppingListsCount, fmt.Sprintf("you can have maximum %d shopping lists", count))
}

func GrpcMaxShoppingListUsersCount(count int) error {
	return fail.CreateGrpcAccessDenied(typeMaxShoppingListsCount, fmt.Sprintf("you can have maximum %d users per shopping list", count))
}
