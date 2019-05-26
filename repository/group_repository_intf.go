package repository

import(
	"github.com/gofrs/uuid"
	"../scim2"
	"go.uber.org/zap"
)

type GroupRepository interface {
	init()
	addGroupResource(groupID uuid.UUID, resource Resource) bool
	deleteGroup(groupID uuid.UUID)
	deleteGroupResource(groupID uuid.UUID, resource uuid.UUID)
	getGroup(groupID uuid.UUID) (*Group, bool)
	setGroup(group *Group)
}

func AddScimGroup(scimGroup *scim2.Group) (uuid.UUID,error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not generate uuid", zap.Error(err))
		return id, err
	}
	group := NewGroup(id,scimGroup.DisplayName)
	if len(scimGroup.Members) > 0 {
		for _, res := range scimGroup.Members {
			resourceId, err1 := uuid.FromString(res.Value)
			userResource, err2 := GetUserResourceMetadata(resourceId)
			if err1 == nil && err2 == nil{
				group.AddResource(userResource)
			}
		}
	}
	groupManager.setGroup(group)
	return id, nil
}

func GetGroup(group uuid.UUID) (*Group, error){
	retGroup, ok := groupManager.getGroup(group)
	if ok {
		return retGroup, nil
	}
	return nil, ErrGroupNotFound
}