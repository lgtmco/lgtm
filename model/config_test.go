package model

import (
	"testing"
)

var deployFile =
`
[deploy]
tasks = [ "t1", "t2", "t3" ]
env = "qa"

[master]
tasks = [ "t1" ]
env = "prod"

[preprod]
env = "preprod"

[preprod2]
tasks = [ "t5" ]
`

func TestDeploymentToml(t *testing.T) {
	d, err := loadDeploymentMap(deployFile)
	if err != nil {
		panic(err)
	}
	if len(d) != 4 {
		t.Fatalf("Should have 4 entries in d, only have %d", len(d))
	}
	deploy, ok := d["deploy"]
	if !ok {
		t.Fatalf("should have an entry for deploy")
	}
	if *deploy.Environment != "qa" {
		t.Fatalf("should have qa for environment, had %s", *deploy.Environment)
	}
	if len(deploy.Tasks) != 3 {
		t.Fatalf("should have had 3 entries in tasks, had %d", len(deploy.Tasks))
	}

	master, ok := d["master"]
	if !ok {
		t.Fatalf("should have an entry for master")
	}
	if *master.Environment != "prod" {
		t.Fatalf("should have prod for environment, had %s", *master.Environment)
	}
	if len(master.Tasks) != 1 {
		t.Fatalf("should have had 1 entry in tasks, had %d", len(master.Tasks))
	}

	preprod, ok := d["preprod"]
	if !ok {
		t.Fatalf("should have an entry for preprod")
	}
	if *preprod.Environment != "preprod" {
		t.Fatalf("should have prod for environment, had %s", *preprod.Environment)
	}
	if len(preprod.Tasks) != 0 {
		t.Fatalf("should have had 0 entries in tasks, had %d", len(preprod.Tasks))
	}

	preprod2, ok := d["preprod2"]
	if !ok {
		t.Fatalf("should have an entry for preprod2")
	}
	if preprod2.Environment != nil {
		t.Fatalf("should have no  environment, had %s", *preprod2.Environment)
	}
	if len(preprod2.Tasks) != 1 {
		t.Fatalf("should have had 1 entry in tasks, had %d", len(preprod2.Tasks))
	}
}
