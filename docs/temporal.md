# Temporal

Notes about the deployed temporal in my home-lab.

## Upgrading Temporal

It seems like Temporal does not automatically upgrade its schema. Therefore it might be required to proceed to manual updates, prior upgrading the server in the cluster:

* Fetch sources
* Build temporal-sql-tool
* Run temporal-sql-tool against the `temporal` db then the `temporal_visibility` db

```sh
$ git clone https://github.com/temporalio/temporal.git
$ cd temporalio/
$ make temporal-sql-tool
$ ./temporal-sql-tool --ep localhost -p 5432 -u temporal-admin -pw (kubectl get secret -n postgres temporal-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq .data.password | base64 -d) --pl postgres12 --db temporal update-schema -d ./schema/postgresql/v12/temporal/versioned/
2024-06-06T16:23:58.529+0200	INFO	UpdateSchemaTask started	{"config": {"DBName":"","TargetVersion":"","SchemaDir":"./schema/postgresql/v12/temporal/versioned/","SchemaName":"","IsDryRun":false}, "logging-call-at": "updatetask.go:105"}
2024-06-06T16:23:58.538+0200	DEBUG	Schema Dirs: []	{"logging-call-at": "updatetask.go:213"}
2024-06-06T16:23:58.538+0200	DEBUG	found zero updates from current version 1.12	{"logging-call-at": "updatetask.go:135"}
2024-06-06T16:23:58.538+0200	INFO	UpdateSchemaTask done	{"logging-call-at": "updatetask.go:128"}

$ ./temporal-sql-tool --ep localhost -p 5432 -u temporal_visibility-admin -pw (kubectl get secret -n postgres temporal-visibility-admin.postgres-instance.credentials.postgresql.acid.zalan.do -o yaml | yq .data.password | base64 -d) --pl postgres12 --db temporal_visibility update-schema -d ./schema/postgresql/v12/visibility/versioned/
2024-06-06T16:24:47.090+0200	INFO	UpdateSchemaTask started	{"config": {"DBName":"","TargetVersion":"","SchemaDir":"./schema/postgresql/v12/visibility/versioned/","SchemaName":"","IsDryRun":false}, "logging-call-at": "updatetask.go:105"}
2024-06-06T16:24:47.100+0200	DEBUG	Schema Dirs: [v1.6]	{"logging-call-at": "updatetask.go:213"}
2024-06-06T16:24:47.100+0200	INFO	Processing schema file: v1.6/fix_root_workflow_info.sql	{"logging-call-at": "updatetask.go:257"}
2024-06-06T16:24:47.101+0200	DEBUG	running 1 updates for current version 1.5	{"logging-call-at": "updatetask.go:139"}
2024-06-06T16:24:47.101+0200	DEBUG	---- Executing updates for version 1.6 ----	{"logging-call-at": "updatetask.go:158"}
2024-06-06T16:24:47.101+0200	DEBUG	DROP INDEX by_root_workflow_id;	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.126+0200	DEBUG	DROP INDEX by_root_run_id;	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.133+0200	DEBUG	ALTER TABLE executions_visibility DROP COLUMN root_workflow_id;	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.147+0200	DEBUG	ALTER TABLE executions_visibility DROP COLUMN root_run_id;	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.153+0200	DEBUG	ALTER TABLE executions_visibility ADD COLUMN root_workflow_id VARCHAR(255) NOT NULL DEFAULT '';	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.178+0200	DEBUG	ALTER TABLE executions_visibility ADD COLUMN root_run_id VARCHAR(255) NOT NULL DEFAULT '';	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.185+0200	DEBUG	CREATE INDEX by_root_workflow_id ON executions_visibility (namespace_id, root_workflow_id, (COALESCE(close_time, '9999-12-31 23:59:59')) DESC, start_time DESC, run_id);	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.236+0200	DEBUG	CREATE INDEX by_root_run_id ON executions_visibility (namespace_id, root_run_id, (COALESCE(close_time, '9999-12-31 23:59:59')) DESC, start_time DESC, run_id);	{"logging-call-at": "updatetask.go:160"}
2024-06-06T16:24:47.288+0200	DEBUG	---- Done ----	{"logging-call-at": "updatetask.go:177"}
2024-06-06T16:24:47.322+0200	DEBUG	Schema updated from 1.5 to 1.6	{"logging-call-at": "updatetask.go:150"}
2024-06-06T16:24:47.322+0200	INFO	UpdateSchemaTask done	{"logging-call-at": "updatetask.go:128"}

```
