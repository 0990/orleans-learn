package catalog

//本来在core文件夹下的文件，但go不允许循环依赖，所以把这个文件移走了
import (
	"github.com/0990/orleans-learn/GrainDirectory"
)

type InternalGrainRuntime struct {
	GrainLocator  *GrainDirectory.GrainLocator
	Catalog       *Catalog
	MessageCenter *MessageCenter
}

//func (r *InternalGrainRuntime) GrainLocator() *GrainDirectory.GrainLocator {
//	return r.grainLocator
//}
//
//func (r *InternalGrainRuntime) Catalog() *catalog.Catalog {
//	return r.catalog
//}
