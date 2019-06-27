package repository

type ResourceTagProvider interface {
	GetResourceTag() *ResourceTag
}