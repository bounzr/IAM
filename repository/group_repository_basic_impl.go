package repository

import "github.com/gofrs/uuid"

type GroupRepositoryBasic struct {
	groups map[uuid.UUID]*Group
}

func NewGroupRepositoryBasic() GroupRepository {
	gr := &GroupRepositoryBasic{}
	gr.init()
	return gr
}

func (g *GroupRepositoryBasic) init() {
	g.groups = make(map[uuid.UUID]*Group)
}

func (g *GroupRepositoryBasic) addGroupResource(group uuid.UUID, resource Resource) bool{
	g.groups[group].AddResource(resource)
	return true
}

func (g *GroupRepositoryBasic) deleteGroup(group uuid.UUID) {
	delete(g.groups, group)
}

func (g *GroupRepositoryBasic) deleteGroupResource(group uuid.UUID, resource uuid.UUID) {
	g.groups[group].DeleteResource(resource)
}

func (g *GroupRepositoryBasic) getGroup(group uuid.UUID) (*Group, bool) {
	grObj, ok := g.groups[group]
	return grObj, ok
}

func (g *GroupRepositoryBasic) setGroup(group *Group) {
	id := group.Metadata.ID
	g.groups[id] = group
}
