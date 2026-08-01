package kubehelpers

import (
	"fmt"
	"strings"

	"git.mkz.me/mycroft/k8s-home/imports/cnpg_cluster_postgresqlcnpgio"
	database "git.mkz.me/mycroft/k8s-home/imports/cnpg_database_postgresqlcnpgio"
	"git.mkz.me/mycroft/k8s-home/imports/externalsecrets_externalsecretsio"
	password "git.mkz.me/mycroft/k8s-home/imports/password_generatorsexternalsecretsio"
	pushsecret "git.mkz.me/mycroft/k8s-home/imports/pushsecret_externalsecretsio"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const (
	// Passwords are generated without symbols so they stay safe to embed in a DSN.
	cnpgPasswordLength = 32
	cnpgPasswordDigits = 6

	cnpgPasswordGeneratorAPIVersion = "generators.external-secrets.io/v1alpha1"
)

// cnpgCredentialsKeys is the ordered list of fields published for every database.
// The order is fixed so synthesized manifests stay stable across runs.
var cnpgCredentialsKeys = []string{"username", "password", "host", "dbname", "url"}

// CNPGDatabase declares one database hosted on the CloudNativePG cluster. Name is
// used for the database, its owning role and the generated secrets. Namespace is
// the application namespace the credentials are published to, and defaults to
// Name when left empty. VaultEntry is the secret name under that namespace and
// defaults to postgresql; point it elsewhere while migrating so the credentials
// of a still-running application are not overwritten.
type CNPGDatabase struct {
	Name            string
	Namespace       string
	VaultEntry      string
	PasswordAliases []string
	Extensions      []string
}

// AppNamespace returns the namespace whose Vault path receives the credentials.
func (db CNPGDatabase) AppNamespace() string {
	if db.Namespace != "" {
		return db.Namespace
	}

	return db.Name
}

// VaultKey is the Vault path the credentials are published to, relative to the
// KV mount. Unlike ExternalSecret keys, which are written with the leading
// secret/ mount prefix and have it stripped on read, PushSecret uses the remote
// ref verbatim: keeping the prefix here would nest the credentials under
// secret/secret/namespaces/, where the application never looks for them.
func (db CNPGDatabase) VaultKey() string {
	secret := db.VaultEntry
	if secret == "" {
		secret = "postgresql"
	}

	return fmt.Sprintf("namespaces/%s/%s", db.AppNamespace(), secret)
}

// resourceName is the shared name of the per-database resources, kept distinct
// from the database name itself so it cannot collide with the cluster's own
// secrets.
func (db CNPGDatabase) resourceName() string {
	return fmt.Sprintf("postgres-%s", strings.ReplaceAll(db.Name, "_", "-"))
}

// CNPGDatabaseConfig describes where a database is provisioned.
type CNPGDatabaseConfig struct {
	// Namespace holding the CloudNativePG cluster.
	Namespace string
	// ClusterName is the Cluster resource the database belongs to.
	ClusterName *string
	// Host is the read-write service applications connect to.
	Host string
	// Database is the database to provision.
	Database CNPGDatabase
}

// NewCNPGManagedRole returns the managed role backing a database. CloudNativePG
// never creates the referenced password secret, so it is generated separately by
// NewCNPGDatabase. Roles have to be known when the Cluster itself is created.
func NewCNPGManagedRole(db CNPGDatabase) *cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRoles {
	return &cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRoles{
		Name:            jsii.String(db.Name),
		Ensure:          cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRolesEnsure_PRESENT,
		Createdb:        jsii.Bool(false),
		Createrole:      jsii.Bool(false),
		DisablePassword: jsii.Bool(false),
		Login:           jsii.Bool(true),
		Superuser:       jsii.Bool(false),
		PasswordSecret: &cnpg_cluster_postgresqlcnpgio.ClusterSpecManagedRolesPasswordSecret{
			Name: jsii.String(db.resourceName()),
		},
	}
}

// NewCNPGDatabase provisions a database on the cluster along with the plumbing
// that gives its owner a password:
//
//  1. a Password generator produces a random password,
//  2. an ExternalSecret materializes it once into a basic-auth secret, which the
//     managed role returned by NewCNPGManagedRole consumes,
//  3. a PushSecret publishes the credentials to Vault under the application
//     namespace, where the application's own ExternalSecret already reads them.
func NewCNPGDatabase(chart constructs.Construct, cfg CNPGDatabaseConfig) {
	db := cfg.Database
	name := db.resourceName()

	password.NewPassword(
		chart,
		jsii.String(fmt.Sprintf("password-%s", db.Name)),
		&password.PasswordProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(cfg.Namespace),
				Name:      jsii.String(name),
			},
			Spec: &password.PasswordSpec{
				Length:      jsii.Number(cnpgPasswordLength),
				Digits:      jsii.Number(cnpgPasswordDigits),
				Symbols:     jsii.Number(0),
				NoUpper:     jsii.Bool(false),
				AllowRepeat: jsii.Bool(true),
			},
		},
	)

	// refreshInterval 0 pins the generated password: without it every refresh
	// would mint a new one and lock the application out until it restarts.
	// The secret is retained so deleting this ExternalSecret cannot drop it.
	externalsecrets_externalsecretsio.NewExternalSecret(
		chart,
		jsii.String(fmt.Sprintf("es-%s", db.Name)),
		&externalsecrets_externalsecretsio.ExternalSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(cfg.Namespace),
				Name:      jsii.String(name),
			},
			Spec: &externalsecrets_externalsecretsio.ExternalSecretSpec{
				RefreshInterval: jsii.String("0"),
				DataFrom: &[]*externalsecrets_externalsecretsio.ExternalSecretSpecDataFrom{
					{
						SourceRef: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRef{
							GeneratorRef: &externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRefGeneratorRef{
								ApiVersion: jsii.String(cnpgPasswordGeneratorAPIVersion),
								Kind:       externalsecrets_externalsecretsio.ExternalSecretSpecDataFromSourceRefGeneratorRefKind_PASSWORD,
								Name:       jsii.String(name),
							},
						},
					},
				},
				Target: &externalsecrets_externalsecretsio.ExternalSecretSpecTarget{
					CreationPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetCreationPolicy_OWNER,
					DeletionPolicy: externalsecrets_externalsecretsio.ExternalSecretSpecTargetDeletionPolicy_RETAIN,
					Immutable:      jsii.Bool(false),
					Name:           jsii.String(name),
					Template: &externalsecrets_externalsecretsio.ExternalSecretSpecTargetTemplate{
						EngineVersion: externalsecrets_externalsecretsio.ExternalSecretSpecTargetTemplateEngineVersion_V2,
						Type:          jsii.String("kubernetes.io/basic-auth"),
						Data:          cnpgCredentialsTemplate(cfg),
					},
				},
			},
		},
	)

	var extensions *[]*database.DatabaseSpecExtensions
	if len(db.Extensions) > 0 {
		items := make([]*database.DatabaseSpecExtensions, 0, len(db.Extensions))
		for _, extension := range db.Extensions {
			items = append(items, &database.DatabaseSpecExtensions{
				Name:   jsii.String(extension),
				Ensure: database.DatabaseSpecExtensionsEnsure_PRESENT,
			})
		}
		extensions = &items
	}

	database.NewDatabase(
		chart,
		jsii.String(fmt.Sprintf("database-%s", db.Name)),
		&database.DatabaseProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(cfg.Namespace),
				Name:      jsii.String(name),
			},
			Spec: &database.DatabaseSpec{
				Cluster: &database.DatabaseSpecCluster{
					Name: cfg.ClusterName,
				},
				Name:       jsii.String(db.Name),
				Owner:      jsii.String(db.Name),
				Extensions: extensions,
			},
		},
	)

	credentialKeys := append([]string{}, cnpgCredentialsKeys...)
	credentialKeys = append(credentialKeys, db.PasswordAliases...)
	pushSecretData := make([]*pushsecret.PushSecretSpecData, 0, len(credentialKeys))
	vaultKey := db.VaultKey()

	for _, key := range credentialKeys {
		pushSecretData = append(pushSecretData, &pushsecret.PushSecretSpecData{
			Match: &pushsecret.PushSecretSpecDataMatch{
				SecretKey: jsii.String(key),
				RemoteRef: &pushsecret.PushSecretSpecDataMatchRemoteRef{
					RemoteKey: jsii.String(vaultKey),
					Property:  jsii.String(key),
				},
			},
		})
	}

	// Vault keeps the credentials on deletion: the application namespace reads
	// this path and losing it would break a running application.
	pushsecret.NewPushSecret(
		chart,
		jsii.String(fmt.Sprintf("push-%s", db.Name)),
		&pushsecret.PushSecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Namespace: jsii.String(cfg.Namespace),
				Name:      jsii.String(name),
			},
			Spec: &pushsecret.PushSecretSpec{
				RefreshInterval: jsii.String("1h"),
				DeletionPolicy:  pushsecret.PushSecretSpecDeletionPolicy_NONE,
				UpdatePolicy:    pushsecret.PushSecretSpecUpdatePolicy_REPLACE,
				SecretStoreRefs: &[]*pushsecret.PushSecretSpecSecretStoreRefs{
					{
						Kind: pushsecret.PushSecretSpecSecretStoreRefsKind_SECRET_STORE,
						Name: jsii.String("secretstore-vault"),
					},
				},
				Selector: &pushsecret.PushSecretSpecSelector{
					Secret: &pushsecret.PushSecretSpecSelectorSecret{
						Name: jsii.String(name),
					},
				},
				Data: &pushSecretData,
			},
		},
	)
}

// cnpgCredentialsTemplate builds the fields published for a database. The url is
// assembled here so applications keep consuming a single ready-made DSN, exactly
// as they did with the Zalando operator.
func cnpgCredentialsTemplate(cfg CNPGDatabaseConfig) *map[string]*string {
	name := cfg.Database.Name

	credentials := map[string]*string{
		"username": jsii.String(name),
		"password": jsii.String("{{ .password }}"),
		"host":     jsii.String(cfg.Host),
		"dbname":   jsii.String(name),
		"url":      jsii.String(fmt.Sprintf("postgres://%s:{{ .password }}@%s/%s", name, cfg.Host, name)),
	}

	for _, alias := range cfg.Database.PasswordAliases {
		credentials[alias] = jsii.String("{{ .password }}")
	}

	return &credentials
}
