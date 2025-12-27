package api

import (
	"fmt"

	"github.com/datatug/datatug-core/pkg/datatug"
	"github.com/datatug/datatug/pkg/sqlexecute"
)

// ExecuteCommands executes command
func ExecuteCommands(storeID string, request sqlexecute.Request) (response sqlexecute.Response, err error) {

	//ctx := context.Background()
	//var store datatug.ProjectStore
	//store, err = storage.GetProjectStore(ctx, storeID, request.Project)
	//if err != nil {
	//	return
	//}

	dbs := make(map[string]*datatug.EnvDb)

	var getEnvDbByID = func(envID, dbID string) (envDb *datatug.EnvDb, err error) {
		key := fmt.Sprintf("%v/%v", envDb, dbID)
		if db, cached := dbs[key]; cached {
			return db, err
		}
		panic("not implemented yet")
		//envServerStore := store.Environments().Environment(envID).Servers().Server(dbID)
		//if envDb, err = envServerStore.Catalogs().Catalog(dbID).LoadEnvironmentCatalog(); err != nil {
		//	return
		//}
		//dbs[key] = envDb
		//return
	}

	var getCatalog = func(server datatug.ServerReference, catalogID string) (*datatug.DbCatalogSummary, error) {
		panic("not implemented yet")
		//serverStore := dal.GetProjectStore(request.Project).DbServers().DbServer(server)
		//return serverStore.Catalogs().DbCatalog(catalogID).LoadDbCatalogSummary(context.Background())
	}

	executor := sqlexecute.NewExecutor(getEnvDbByID, getCatalog)
	return executor.Execute(request)
}
