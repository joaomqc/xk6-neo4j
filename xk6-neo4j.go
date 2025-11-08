package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/neo4j", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
		// comparator is the exported type
		comparator *Neo4j
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:         vu,
		comparator: &Neo4j{vu: vu},
	}
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.comparator,
	}
}

type Neo4j struct {
	vu modules.VU
}

type Neo4jDriver struct {
	driver neo4j.Driver
	vu     modules.VU
}

type DriverConfig struct {
	Uri      string
	User     string
	Password string
	Realm    string
}

func (d *Neo4j) NewDriver(conf DriverConfig) Neo4jDriver {
	driver, err := neo4j.NewDriver(conf.Uri, neo4j.BasicAuth(conf.User, conf.Password, conf.Realm))
	if err != nil {
		panic(fmt.Errorf("connecting to neo4j: %w", err))
	}
	return Neo4jDriver{
		driver: driver,
		vu:     d.vu,
	}
}

func (d *Neo4jDriver) Read(query string, parameters map[string]any) []neo4j.Record {
	return d.ExecuteQuery(neo4j.AccessModeRead, query, parameters)
}

func (d *Neo4jDriver) Write(query string, parameters map[string]any) []neo4j.Record {
	return d.ExecuteQuery(neo4j.AccessModeWrite, query, parameters)
}

func (d *Neo4jDriver) ExecuteQuery(accessMode neo4j.AccessMode, query string, parameters map[string]any) []neo4j.Record {
	ctx := d.vu.Context()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: accessMode,
	})
	defer func() {
		if err := session.Close(ctx); err != nil {
			panic(fmt.Errorf("failed to close session: %w", err))
		}
	}()
	result, err := session.Run(ctx, query, parameters)
	if err != nil {
		panic(fmt.Errorf("failed to execute query: %w", err))
	}
	records := []neo4j.Record{}
	for result.Next(ctx) {
		records = append(records, *result.Record())
	}
	return records
}

func (d *Neo4jDriver) Close() {
	ctx := d.vu.Context()
	err := d.driver.Close(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to close connection: %w", err))
	}
}
