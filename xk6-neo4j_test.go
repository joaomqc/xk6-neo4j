package neo4j

import (
	"context"
	"errors"
	"testing"

	"github.com/grafana/sobek"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

func TestRead(t *testing.T) {
	n := Neo4j{
		vu: mockVU{},
	}
	d := n.NewDriver(DriverConfig{
		Uri:      "bolt://localhost:7687",
		User:     "neo4j",
		Password: "neo4jpass",
		Realm:    "",
	})
	r := d.Read("MATCH (p:Person) WHERE p.name = 'Errico' RETURN p;", nil)
	if len(r) == 0 {
		t.Error("empty return")
		return
	}
	item, ok := r[0].Get("p")
	if !ok {
		t.Error("key 'p' found")
	}
	node := item.(neo4j.Node)
	v, ok := node.Props["name"]
	if !ok {
		t.Error("field 'name' not found")
		return
	}
	name, ok := v.(string)
	if !ok {
		t.Errorf("field 'name' is not of type string")
		return
	}
	if name != "Errico" {
		t.Errorf("expected field 'name' to have value 'Errico', found value '%s'", name)
		return
	}
}

func TestWrite(t *testing.T) {
	n := Neo4j{
		vu: mockVU{},
	}
	d := n.NewDriver(DriverConfig{
		Uri:      "bolt://localhost:7687",
		User:     "neo4j",
		Password: "neo4jpass",
		Realm:    "",
	})
	params := map[string]any{
		"name":    "Errico",
		"country": "Italia",
		"year":    1853,
	}
	// function write(query string, params object)
	r := d.Write("CREATE (p:Person {name: $name, country: $country, year: $year}) RETURN p;", params)
	if len(r) == 0 {
		t.Error("empty return")
		return
	}
	item, ok := r[0].Get("p")
	if !ok {
		t.Error("key 'p' found")
	}
	node := item.(neo4j.Node)
	v, ok := node.Props["name"]
	if !ok {
		t.Error("field 'name' not found")
		return
	}
	name, ok := v.(string)
	if !ok {
		t.Errorf("field 'name' is not of type string")
		return
	}
	if name != "Errico" {
		t.Errorf("expected field 'name' to have value 'Errico', found value '%s'", name)
		return
	}
}

func TestExecuteQuery(t *testing.T) {
	n := Neo4j{
		vu: mockVU{},
	}
	d := n.NewDriver(DriverConfig{
		Uri:      "bolt://localhost:7687",
		User:     "neo4j",
		Password: "neo4jpass",
		Realm:    "",
	})
	params := map[string]any{
		"name":    "Errico",
		"country": "Italia",
		"year":    1853,
	}
	// function write(query string, params object)
	r := d.ExecuteQuery(0, "CREATE (p:Person {name: $name, country: $country, year: $year}) RETURN p;", params)
	if len(r) == 0 {
		t.Error("empty return")
		return
	}
	item, ok := r[0].Get("p")
	if !ok {
		t.Error("key 'p' found")
	}
	node := item.(neo4j.Node)
	v, ok := node.Props["name"]
	if !ok {
		t.Error("field 'name' not found")
		return
	}
	name, ok := v.(string)
	if !ok {
		t.Errorf("field 'name' is not of type string")
		return
	}
	if name != "Errico" {
		t.Errorf("expected field 'name' to have value 'Errico', found value '%s'", name)
		return
	}
}

func TestExecuteQuery_IncorrectAccessMode(t *testing.T) {
	n := Neo4j{
		vu: mockVU{},
	}
	d := n.NewDriver(DriverConfig{
		Uri:      "bolt://localhost:7687",
		User:     "neo4j",
		Password: "neo4jpass",
		Realm:    "",
	})
	params := map[string]any{
		"name":    "Errico",
		"country": "Italia",
		"year":    1853,
	}
	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected panic due to incorrect access mode, but there was none")
		}
		var neo4jError *neo4j.Neo4jError
		if !errors.As(err.(error), &neo4jError) {
			t.Errorf("expected error of type 'neo4j.Neo4jError', but found type '%T': '%s'", err, err)
			return
		}
		if neo4jError.Code != "Neo.ClientError.Statement.AccessMode" {
			t.Errorf("expected error code to be 'Neo.ClientError.Statement.AccessMode', but found code '%s'", neo4jError.Code)
			return
		}
	}()
	// function write(query string, params object)
	d.ExecuteQuery(1, "CREATE (p:Person {name: $name, country: $country, year: $year}) RETURN p;", params)
}

type mockVU struct{}

func (vu mockVU) Context() context.Context {
	return context.TODO()
}

func (vu mockVU) Events() common.Events {
	return common.Events{}
}

func (vu mockVU) InitEnv() *common.InitEnvironment {
	return nil
}

func (vu mockVU) State() *lib.State {
	return nil
}

func (vu mockVU) Runtime() *sobek.Runtime {
	return nil
}

func (vu mockVU) RegisterCallback() (enqueueCallback func(func() error)) {
	return nil
}
