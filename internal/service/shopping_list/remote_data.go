package shopping_list

import (
	"context"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/utils/slices"
	profileApi "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	"time"
)

func (s *Service) getProfilesInfo(authorIds []string) map[string]*profileApi.ProfileMinInfo {
	uniqueAuthorIds := slices.RemoveDuplicates(authorIds)
	infos := make(map[string]*profileApi.ProfileMinInfo)
	if len(uniqueAuthorIds) == 0 {
		return infos
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := s.grpc.Profile.GetProfilesMinInfo(ctx, &profileApi.GetProfilesMinInfoRequest{ProfileIds: uniqueAuthorIds})
	cancelCtx()

	if err == nil {
		for _, authorId := range uniqueAuthorIds {
			if info, ok := res.Infos[authorId]; ok {
				infos[authorId] = info
			}
		}
	} else {
		log.Warn("unable to get recipe authors data: %s", err)
	}

	return infos
}
