package repository

import (
	"../config"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type GroupManager interface {
	setGroupResource(groupID uuid.UUID, resource ResourceTagger) bool
	close()
	deleteGroup(groupID uuid.UUID)
	deleteGroupResource(groupID uuid.UUID, resource uuid.UUID)
	findGroups(filter map[string]interface{}) ([]Group, error)
	getGroup(groupID uuid.UUID) (*Group, bool)
	init()
	setGroup(group *Group)
}

func initGroups() {
	implementation := config.IAM.Groups.Implementation
	switch implementation {
	case "leveldb":
		groupManager = &GroupManagerLeveldb{groupsPath: "./rep/group"}
	default:
		groupManager = &GroupManagerBasic{}
	}
	groupManager.init()
}

func AddGroup(name string) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not generate uuid", zap.Error(err))
		return id, err
	}
	group := NewGroup(id, name)
	groupManager.setGroup(group)
	return id, nil
}

func AddGroupResource(group uuid.UUID, resource ResourceTagger) {
	groupManager.setGroupResource(group, resource)
}

func AddScimGroup(scimGroup *scim2.Group) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not generate uuid", zap.Error(err))
		return id, err
	}
	group := NewGroup(id, scimGroup.DisplayName)
	if len(scimGroup.Members) > 0 {
		for _, res := range scimGroup.Members {
			resourceId, err1 := uuid.FromString(res.Value)
			userResource, ok := GetResourceMetadata(resourceId)
			if err1 == nil && ok {
				group.AddResource(userResource)
			}
		}
	}
	groupManager.setGroup(group)
	return id, nil
}

func FindGroups(conditions map[string]interface{}) []scim2.Group {
	var groups []scim2.Group
	repGroups, err := groupManager.findGroups(conditions)
	if err != nil {
		log.Error("can not get groups from repository", zap.Error(err))
	}
	for _, group := range repGroups {
		groups = append(groups, *group.GetScim())
	}
	return groups
}

func FindGroupAssignments(conditions map[string]interface{}) []scim2.GroupAssignment {
	var groupAssignments []scim2.GroupAssignment
	repGroups, err := groupManager.findGroups(conditions)
	if err != nil {
		log.Error("can not get groups from repository", zap.Error(err))
	}
	for _, group := range repGroups {
		groupAssignments = append(groupAssignments, *group.GetScim().GetGroupAssignment())
	}
	return groupAssignments
}

func GetGroup(group uuid.UUID) (*Group, error) {
	retGroup, ok := groupManager.getGroup(group)
	if ok {
		return retGroup, nil
	}
	return nil, ErrGroupNotFound
}
