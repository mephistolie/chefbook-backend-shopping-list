package fail

import (
	"fmt"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
)

var (
	typeOutdatedVersion       = "outdated_version"
	typeMaxShoppingListsCount = "max_shopping_lists_count"
	typePersonalShoppingList  = "personal_shopping_list"
)

var (
	GrpcOutdatedVersion        = fail.CreateGrpcConflict(typeOutdatedVersion, "shopping list version is outdated; process current version first")
	GrpcShoppingListNotFound   = fail.CreateGrpcClient(fail.TypeNotFound, "shopping list not found")
	GrpcPersonalShoppingList   = fail.CreateGrpcClient(typePersonalShoppingList, "personal shopping list can't be shared")
	GrpcInvalidShoppingListKey = fail.CreateGrpcAccessDenied(fail.TypeAccessDenied, "invalid shopping list key")
)

func GrpcMaxShoppingListsCount(count int) error {
	return fail.CreateGrpcAccessDenied(typeMaxShoppingListsCount, fmt.Sprintf("you can have maximum %d shopping lists", count))
}

func GrpcMaxShoppingListUsersCount(count int) error {
	return fail.CreateGrpcAccessDenied(typeMaxShoppingListsCount, fmt.Sprintf("you can have maximum %d users per shopping list", count))
}
