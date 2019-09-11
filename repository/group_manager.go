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
	deleteResource(resource uuid.UUID)
	findGroups(filter map[string]interface{}) ([]Group, error)
	getGroup(groupID uuid.UUID) (*Group, bool)
	init()
	setGroup(group *Group)
}

var privateGroupsList = [...]string{"Admins", "Clients", "ProtectedResources"}
var privateGroups = make(map[string]uuid.UUID)

func initGroups() {
	implementation := config.IAM.Groups.Implementation
	switch implementation {
	case "leveldb":
		groupManager = &GroupManagerLeveldb{groupsPath: "./rep/group"}
	default:
		groupManager = &GroupManagerBasic{}
	}
	groupManager.init()
	addPrivateGroups()
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

func addPrivateGroups() {
	for _, group := range privateGroupsList {
		groupFilter := make(map[string]interface{})
		groupFilter["name"] = group
		groups, err := groupManager.findGroups(groupFilter)
		if err != nil {
			log.Debug("group not found. The server will try to add it", zap.String("name", group), zap.Error(err))
		}
		var groupID uuid.UUID
		if len(groups) == 0 {
			groupID, err = AddGroup(group)
			if err != nil {
				log.Error("can not add group", zap.String("name", group), zap.Error(err))
				panic("private group is required to run")
			}
		} else {
			groupID = groups[0].Metadata.ID
		}
		privateGroups[group] = groupID
	}
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
			resourceId := uuid.FromStringOrNil(res.Value)
			user, found := GetUser(resourceId)
			if found {
				userResource := user.GetResourceTag()
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

func SetResourceGroups(assigner scim2.GroupAssigner, resource ResourceTagger) {
	groupManager.deleteResource(resource.GetUUID())
	groupAssignment := assigner.GetGroups()
	for _, group := range groupAssignment {
		groupID, err := uuid.FromString(group.Value)
		if err != nil {
			log.Error("invalid group groupID in group assigner", zap.String("groupID", group.Value), zap.Error(err))
			continue
		}
		_, err = GetGroup(groupID)
		if err != nil {
			log.Error("unknown group groupID in group assigner", zap.String("groupID", group.Value), zap.Error(err))
			continue
		}
		AddGroupResource(groupID, resource)
	}
}

func ValidateResourceInGroup(resourceID uuid.UUID, groupName string) bool {
	filter := make(map[string]interface{})
	filter["name"] = groupName
	filter["member"] = resourceID
	groups, err := groupManager.findGroups(filter)
	if err != nil {
		log.Error("can not get groups from repository", zap.Error(err))
		return false
	}
	if len(groups) > 0 {
		return true
	}
	return false
}
