package api

import (
	"database/sql"
	"fmt"
	"github.com/datatug/datatug/packages/execute"
	"github.com/datatug/datatug/packages/models"
	"github.com/datatug/datatug/packages/parallel"
	"github.com/datatug/datatug/packages/schemer"
	"github.com/datatug/datatug/packages/slice"
	"github.com/strongo/random"
	"log"
	"strings"
	"time"
)

// ProjectLoader defines an interface to load project info
type ProjectLoader interface {
	GetProjectSummary(id string) (projectSummary models.ProjectSummary, err error)
	GetProject(id string) (project *models.DataTugProject, err error)
}

// UpdateDbSchema updates DB schema
func UpdateDbSchema(loader ProjectLoader, projectID, environment, driver, dbModelID string, connectionString execute.ConnectionString) (project *models.DataTugProject, err error) {
	var (
		latestDb *models.Database
	)
	var projSummaryErr error
	getProjectSummaryWorker := func() error {
		_, projSummaryErr = loader.GetProjectSummary(projectID)
		if models.ProjectDoesNotExist(projSummaryErr) {
			return nil
		}
		return projSummaryErr
	}
	dbServer := models.DbServer{
		Driver: driver,
		Host:   connectionString.Server(),
		Port:   connectionString.Port(),
	}
	scanDbWorker := func() error {
		var scanErr error
		if latestDb, scanErr = scanDbSchema(dbServer, connectionString); err != nil {
			return scanErr
		}
		return scanErr
	}
	if err = parallel.Run(getProjectSummaryWorker, scanDbWorker); err != nil {
		return project, err
	}
	if dbModelID != "" {
		latestDb.DbModel = dbModelID
	} else if latestDb.DbModel == "" {
		latestDb.DbModel = latestDb.ID
	}
	if models.ProjectDoesNotExist(projSummaryErr) {
		log.Println("Creating a new DataTug project...")
		if project, err = newProjectWithDatabase(environment, dbServer, latestDb); err != nil {
			return project, err
		}
	} else {
		log.Printf("Loading existign project...")
		if project, err = loader.GetProject(projectID); err != nil {
			err = fmt.Errorf("failed to load DataTug project: %w", err)
			return
		}
		log.Println("Updating project with latest database info...", environment)
		if err = updateProjectWithDatabase(project, environment, dbServer, latestDb); err != nil {
			return project, err
		}
	}
	dbModel := project.DbModels.GetDbModelByID(latestDb.DbModel)
	if dbModel == nil {
		err = fmt.Errorf("db model not found by ID: %v. there is %v db models in the project: %v", latestDb.DbModel, len(project.DbModels), strings.Join(project.DbModels.IDs(), ", "))
		return
	}
	if err = updateDbModelWithDatabase(environment, dbModel, latestDb); err != nil {
		err = fmt.Errorf("failed to update dbModel with database: %w", err)
		return
	}
	return project, err
}

func updateProjectWithDatabase(project *models.DataTugProject, envID string, dbServer models.DbServer, latestDb *models.Database) (err error) {
	// Update environment
	{
		if environment := project.Environments.GetEnvByID(envID); environment == nil {
			environment = &models.Environment{
				ProjectEntity: models.ProjectEntity{
					ID: envID,
				},
				DbServers: models.EnvDbServers{
					{
						DbServer:  dbServer,
						Databases: []string{latestDb.ID},
					},
				},
			}
			project.Environments = append(project.Environments, environment)
		} else if envDbServer := environment.DbServers.GetByID(dbServer.ID()); envDbServer == nil {
			environment.DbServers = append(environment.DbServers, &models.EnvDbServer{
				DbServer:  dbServer,
				Databases: []string{latestDb.ID},
			})
		} else if i := slice.IndexOfString(envDbServer.Databases, latestDb.ID); i < 0 {
			envDbServer.Databases = append(envDbServer.Databases, latestDb.ID)
		}
	}
	// Update DB server
	{
		for i, projDbServer := range project.DbServers {
			if projDbServer.ProjectEntity.ID == dbServer.ID() {
				for j, db := range project.DbServers[i].Databases {
					if db.ID == latestDb.ID {
						project.DbServers[i].Databases[j] = latestDb
						goto ProjectDbServerDatabaseUpdated
					}
				}
				project.DbServers[i].Databases = append(project.DbServers[i].Databases, latestDb)
			ProjectDbServerDatabaseUpdated:
				goto ProjDbServerUpdate
			}
		}
		project.DbServers = append(project.DbServers, &models.ProjDbServer{
			ProjectEntity: models.ProjectEntity{ID: dbServer.ID()},
			DbServer:      dbServer,
			Databases:     models.Databases{latestDb},
		})
	ProjDbServerUpdate:
	}
	return nil
}

func newProjectWithDatabase(environment string, dbServer models.DbServer, database *models.Database) (project *models.DataTugProject, err error) {
	//var currentUser *user.User
	//if currentUser, err = user.Current(); err != nil {
	//	err = fmt.Errorf("failed to get current OS user")
	//	return
	//}
	project = &models.DataTugProject{
		Access: "private",
		Created: &models.ProjectCreated{
			//ByUsername: currentUser.Username,
			//ByName:     currentUser.Name,
			At: time.Now(),
		},
		DbModels: models.DbModels{
			&models.DbModel{
				ProjectEntity: models.ProjectEntity{ID: database.DbModel},
			},
		},
		Environments: models.Environments{
			{
				ProjectEntity: models.ProjectEntity{ID: environment},
				DbServers: []*models.EnvDbServer{
					{
						DbServer:  dbServer,
						Databases: []string{database.ID},
					},
				},
			},
		},
		DbServers: models.ProjDbServers{
			{
				ProjectEntity: models.ProjectEntity{ID: dbServer.ID()},
				DbServer:      dbServer,
				Databases: models.Databases{
					database,
				},
			},
		},
	}
	project.ID = random.ID(9)
	log.Println("project.ID:", project.ID)
	return project, err
}

func scanDbSchema(server models.DbServer, connectionString execute.ConnectionString) (database *models.Database, err error) {
	var db *sql.DB

	if db, err = sql.Open(server.Driver, connectionString.String()); err != nil {
		return nil, fmt.Errorf("failed to open SQL db: %w", err)
	}

	// Close the database connection pool after command executes
	defer func() { _ = db.Close() }()

	informationSchema := schemer.NewInformationSchema(server, db)

	if database, err = informationSchema.GetDatabase(connectionString.Database()); err != nil {
		return nil, fmt.Errorf("failed to get database metadata: %w", err)
	}
	return
}

func updateDbModelWithDatabase(envID string, dbModel *models.DbModel, database *models.Database) (err error) {
	{ // Update dbmodel environments
		environment := dbModel.Environments.GetByID(envID)
		if environment == nil {
			environment = &models.DbModelEnv{ID: envID}
			dbModel.Environments = append(dbModel.Environments, environment)
		}
		dbModelDb := environment.Databases.GetByID(database.ID)
		if dbModelDb == nil {
			environment.Databases = append(environment.Databases, &models.DbModelDb{
				ID: database.ID,
			})
		}
	}

	for _, schema := range database.Schemas {
		var schemaModel *models.SchemaModel
		for _, sm := range dbModel.Schemas {
			if sm.ID == schema.ID {
				schemaModel = sm
				goto UpdateSchemaModel
			}
		}
		schemaModel = &models.SchemaModel{
			ProjectEntity: schema.ProjectEntity,
		}
		dbModel.Schemas = append(dbModel.Schemas, schemaModel)
	UpdateSchemaModel:
		if err = updateSchemaModel(envID, schemaModel, schema); err != nil {
			return fmt.Errorf("faild to update DB schema model: %w", err)
		}
	}
	return nil
}

func updateSchemaModel(envID string, schemaModel *models.SchemaModel, dbSchema *models.DbSchema) (err error) {
	updateTables := func(tables []*models.Table) (result models.TableModels) {
		for _, table := range tables {
			tableModel := schemaModel.Tables.GetByKey(table.TableKey)
			if tableModel == nil {
				tableModel = &models.TableModel{
					TableKey: table.TableKey,
					ByEnv:    make(models.StateByEnv),
				}
				tableModel.ByEnv[envID] = &models.EnvState{
					Status: "exists",
				}
				tableModel.Columns = make(models.ColumnModels, len(table.Columns))
				for i, c := range table.Columns {
					tableModel.Columns[i] = &models.ColumnModel{
						Column: *c,
						ByEnv:  make(models.StateByEnv),
					}
					tableModel.Columns[i].ByEnv[envID] = &models.EnvState{
						Status: "exists",
					}
				}
				result = append(result, tableModel)
			} else {
				panic("not implemented")
			}
		}
		return
	}
	schemaModel.Tables = updateTables(dbSchema.Tables)
	schemaModel.Views = updateTables(dbSchema.Views)
	return nil
}
