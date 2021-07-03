package filestore

import (
	"context"
	"fmt"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/parallel"
	"github.com/datatug/datatug/packages/storage"
	"log"
	"os"
	"path"
)

var _ storage.EnvironmentStore = (*fsEnvironmentStore)(nil)

type fsEnvironmentStore struct {
	envID   string
	envPath string
	fsEnvironmentsStore
}

func (store fsEnvironmentStore) Project() storage.ProjectStore {
	return store.fsProjectStore
}

func (store fsEnvironmentStore) ID() string {
	return store.envID
}

func (store fsEnvironmentStore) Servers() storage.EnvServersStore {
	return newFsEnvServersStore(store)
}

func (store fsEnvironmentStore) DeleteEnvironment() (err error) {
	panic("implement me")
}

func (store fsEnvironmentStore) SaveEnvironment(environment *models.Environment) (err error) {
	panic("implement me")
}

func newFsEnvironmentStore(id string, fsEnvironmentsStore fsEnvironmentsStore) fsEnvironmentStore {
	return fsEnvironmentStore{
		envID:               id,
		fsEnvironmentsStore: fsEnvironmentsStore,
	}
}

// GetEnvironmentSummary loads environment summary
func (store fsEnvironmentStore) LoadEnvironmentSummary() (*models.EnvironmentSummary, error) {
	envSummary, err := loadEnvFile(store.envsDirPath, store.envID)
	if err != nil {
		err = fmt.Errorf("failed to load environment [%v] from project [%v]: %w", store.envID, store.projectID, err)
		return nil, err
	}
	return &envSummary, err
}

// GetEnvironmentDbSummary return DB summary for specific environment
func (store fsEnvironmentStore) LoadEnvironmentDbSummary(databaseID string) (models.DbCatalogSummary, error) {
	panic(fmt.Sprintf("implement me: %v, %v, %v", store.projectID, store.envID, databaseID))
}

func (store fsEnvironmentStore) saveEnvironment(ctx context.Context, env models.Environment) (err error) {
	dirPath := path.Join(store.projectPath, DatatugFolder, EnvironmentsFolder, env.ID)
	log.Printf("Saving environment [%v]: %v", env.ID, dirPath)
	if err = os.MkdirAll(dirPath, 0777); err != nil {
		return fmt.Errorf("failed to create environemtn folder: %w", err)
	}
	return parallel.Run(
		func() error {
			if err = saveJSONFile(dirPath, jsonFileName(env.ID, environmentFileSuffix), models.EnvironmentFile{ID: env.ID}); err != nil {
				return fmt.Errorf("failed to write environment json to file: %w", err)
			}
			return nil
		},
		func() error {
			envServers := newFsEnvServersStore(store)
			if err = envServers.saveEnvServers(env.DbServers); err != nil {
				return fmt.Errorf("failed to save environment servers: %w", err)
			}
			return nil
		},
	)
}
