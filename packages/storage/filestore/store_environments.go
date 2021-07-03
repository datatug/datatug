package filestore

import (
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/storage"
	"path"
)

var _ storage.EnvironmentsStore = (*fsEnvironmentsStore)(nil)

type fsEnvironmentsStore struct {
	envsDirPath string
	fsProjectStoreRef
}

func (store fsEnvironmentsStore) Environment(id string) storage.EnvironmentStore {
	return store.environment(id)
}

func (store fsEnvironmentsStore) environment(id string) fsEnvironmentStore {
	return newFsEnvironmentStore(id, store)
}

func newFsEnvironmentsStore(fsProjectStore fsProjectStore) fsEnvironmentsStore {
	return fsEnvironmentsStore{
		envsDirPath:       path.Join(fsProjectStore.projectPath, DatatugFolder, EnvironmentsFolder),
		fsProjectStoreRef: fsProjectStoreRef{fsProjectStore},
	}
}

func (store fsEnvironmentsStore) saveEnvironments(project models.DatatugProject) (err error) {
	return saveItems("environments", len(project.Environments), func(i int) func() error {
		return func() error {
			env := *project.Environments[i]
			envStore := store.environment(env.ID)
			return envStore.saveEnvironment(env)
		}
	})
}
