package repository

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
)

type GroupManagerBasic struct {
	groups map[uuid.UUID]*Group
}

func (g *GroupManagerBasic) init() {
	g.groups = make(map[uuid.UUID]*Group)
}

func (g *GroupManagerBasic) setGroupResource(group uuid.UUID, resource ResourceTagger) bool {
	g.groups[group].AddResource(resource)
	return true
}

func (g *GroupManagerBasic) close() {
	//nothing
}

func (g *GroupManagerBasic) deleteGroup(group uuid.UUID) {
	delete(g.groups, group)
}

func (g *GroupManagerBasic) deleteGroupResource(group uuid.UUID, resource uuid.UUID) {
	g.groups[group].DeleteResource(resource)
}

func (g *GroupManagerBasic) getGroup(group uuid.UUID) (*Group, bool) {
	grObj, ok := g.groups[group]
	return grObj, ok
}

//todo atomic operation including len(g.groups) and range(g.groups)
func (g *GroupManagerBasic) findGroups(conditions map[string]interface{}) ([]Group, error) {
	groups := make([]Group, len(g.groups))
	idx := 0
	for _, group := range g.groups {
		nameCondition := conditions["name"]
		var nameToFind string
		if nameCondition != nil {
			nameToFind = nameCondition.(string)
		}
		if len(nameToFind) > 0 {
			log.Debug("searching with name condition", zap.String("name", conditions["name"].(string)))
			if strings.Compare(nameToFind, group.Metadata.Name) != 0 {
				continue
			}
		}
		groups[idx] = *group
		idx++
	}
	return groups, nil
}

func (g *GroupManagerBasic) setGroup(group *Group) {
	id := group.Metadata.ID
	g.groups[id] = group
}
