package provisionerbeta

import (
	"time"

	"github.com/pkg/errors"
	nebula "github.com/slackhq/nebula/cert"
	"github.com/smallstep/cli/utils"
	"github.com/urfave/cli"
	"go.step.sm/cli-utils/errs"
)

// Command returns the jwk subcommand.
func Command() cli.Command {
	return cli.Command{
		Name:      "provisioner",
		Usage:     "create and manage the certificate authority provisioners",
		UsageText: "**step beta ca provisioner** <subcommand> [arguments] [global-flags] [subcommand-flags]",
		Subcommands: cli.Commands{
			//listCommand(),
			addCommand(),
			removeCommand(),
			getCommand(),
			updateCommand(),
		},
		Description: `**step beta ca provisioner** command group provides facilities for managing the
certificate authority provisioners.

A provisioner is an entity that controls provisioning credentials, which are
used to generate provisioning tokens.

Provisioning credentials are simple JWK key pairs using public-key cryptography.
The public key is used to verify a provisioning token while the private key is
used to sign the provisioning token.

Provisioning tokens are JWT tokens signed by the JWK private key. These JWT
tokens are used to get a valid TLS certificate from the certificate authority.
Each provisioner is able to manage a different set of rules that can be used to
configure the bounds of the certificate.

In the certificate authority, a provisioner is configured with a JSON object
with the following properties:

* **name**: the provisioner name, it will become the JWT issuer and a good
  practice is to use an email address for this.
* **type**: the provisioner type, currently only "jwk" is supported.
* **key**: the JWK public key used to verify the provisioning tokens.
* **encryptedKey** (optional): the JWE compact serialization of the private key
  used to sign the provisioning tokens.
* **claims** (optional): an object with custom options for each provisioner.
  Options supported are:
  * **minTLSCertDuration**: minimum duration of a certificate, set to 5m by
    default.
  * **maxTLSCertDuration**: maximum duration of a certificate, set to 24h by
    default.
  * **defaultTLSCertDuration**: default duration of the certificate, set to 24h
    by default.
  * **disableRenewal**: whether or not to disable certificate renewal, set to false
    by default.

## EXAMPLES

Add a single provisioner:
'''
$ step beta ca provisioner add max@smallstep.com --type JWK --create
'''

Remove a provisioner:
'''
$ step beta ca provisioner remove max@smallstep.com
'''`,
	}
}

func parseIntaceAge(ctx *cli.Context) (age string, err error) {
	if !ctx.IsSet("instance-age") {
		return
	}
	age = ctx.String("instance-age")
	dur, err := time.ParseDuration(age)
	if err != nil {
		return "", err
	}
	if dur < 0 {
		return "", errs.MinSizeFlag(ctx, "instance-age", "0s")
	}
	return
}

func removeElements(list, rems []string) []string {
	if len(list) == 0 {
		return list
	}
	for _, rem := range rems {
		for i, elem := range list {
			if elem == rem {
				list[i] = list[len(list)-1]
				list = list[:len(list)-1]
				break
			}
		}
	}
	return list
}

var (
	x509TemplateFlag = cli.StringFlag{
		Name:  "x509-template",
		Usage: `The x509 certificate template <file>, a JSON representation of the certificate to create.`,
	}
	x509TemplateDataFlag = cli.StringFlag{
		Name:  "x509-template-data",
		Usage: `The x509 certificate template data <file>, a JSON map of data that can be used by the certificate template.`,
	}
	sshTemplateFlag = cli.StringFlag{
		Name:  "ssh-template",
		Usage: `The x509 certificate template <file>, a JSON representation of the certificate to create.`,
	}
	sshTemplateDataFlag = cli.StringFlag{
		Name:  "ssh-template-data",
		Usage: `The ssh certificate template data <file>, a JSON map of data that can be used by the certificate template.`,
	}
	x509MinDurFlag = cli.StringFlag{
		Name:  "x509-min-dur",
		Usage: `The minimum <duration> for an x509 certificate generated by this provisioner.`,
	}
	x509MaxDurFlag = cli.StringFlag{
		Name:  "x509-max-dur",
		Usage: `The maximum <duration> for an x509 certificate generated by this provisioner.`,
	}
	x509DefaultDurFlag = cli.StringFlag{
		Name:  "x509-default-dur",
		Usage: `The default <duration> for an x509 certificate generated by this provisioner.`,
	}
	sshUserMinDurFlag = cli.StringFlag{
		Name:  "ssh-user-min-dur",
		Usage: `The minimum <duration> for an ssh user certificate generated by this provisioner.`,
	}
	sshUserMaxDurFlag = cli.StringFlag{
		Name:  "ssh-user-max-dur",
		Usage: `The maximum <duration> for an ssh user certificate generated by this provisioner.`,
	}
	sshUserDefaultDurFlag = cli.StringFlag{
		Name:  "ssh-user-default-dur",
		Usage: `The maximum <duration> for an ssh user certificate generated by this provisioner.`,
	}
	sshHostMinDurFlag = cli.StringFlag{
		Name:  "ssh-host-min-dur",
		Usage: `The minimum <duration> for an ssh host certificate generated by this provisioner.`,
	}
	sshHostMaxDurFlag = cli.StringFlag{
		Name:  "ssh-host-max-dur",
		Usage: `The maximum <duration> for an ssh host certificate generated by this provisioner.`,
	}
	sshHostDefaultDurFlag = cli.StringFlag{
		Name:  "ssh-host-default-dur",
		Usage: `The maximum <duration> for an ssh host certificate generated by this provisioner.`,
	}
	disableRenewalFlag = cli.BoolFlag{
		Name:  "disable-renewal",
		Usage: `Disable renewal for all certificates generated by this provisioner`,
	}
	enableX509Flag = cli.BoolFlag{
		Name:  "x509",
		Usage: `Enable provisioning of x509 certificates.`,
	}
	enableSSHFlag = cli.BoolFlag{
		Name:  "ssh",
		Usage: `Enable provisioning of ssh certificates.`,
	}
	forceCNFlag = cli.BoolFlag{
		Name:  "force-cn",
		Usage: `Always set the common name in provisioned certificates.`,
	}

	// Cloud provisioner flags
	awsAccountFlag = cli.StringSliceFlag{
		Name: "aws-account",
		Usage: `The AWS account <id> used to validate the identity documents.
Use the flag multiple times to configure multiple accounts.`,
	}
	removeAWSAccountFlag = cli.StringSliceFlag{
		Name: "remove-aws-account",
		Usage: `Remove an AWS account <id> used to validate the identity documents.
Use the flag multiple times to remove multiple accounts.`,
	}
	azureTenantFlag = cli.StringFlag{
		Name:  "azure-tenant",
		Usage: `The Microsoft Azure tenant <id> used to validate the identity tokens.`,
	}
	azureResourceGroupFlag = cli.StringSliceFlag{
		Name: "azure-resource-group",
		Usage: `The Microsoft Azure resource group <name> used to validate the identity tokens.
Use the flag multipl etimes to configure multiple resource groups`,
	}
	removeAzureResourceGroupFlag = cli.StringSliceFlag{
		Name: "remove-azure-resource-group",
		Usage: `Remove a Microsoft Azure resource group <name> used to validate the identity tokens.
Use the flag multipl etimes to configure multiple resource groups`,
	}
	gcpServiceAccountFlag = cli.StringSliceFlag{
		Name: "gcp-service-account",
		Usage: `The Google service account <email> or <id> used to validate the identity tokens.
Use the flag multiple times to configure multiple service accounts.`,
	}
	removeGCPServiceAccountFlag = cli.StringSliceFlag{
		Name: "remove-gcp-service-account",
		Usage: `Remove a Google service account <email> or <id> used to validate the identity tokens.
Use the flag multiple times to configure multiple service accounts.`,
	}
	gcpProjectFlag = cli.StringSliceFlag{
		Name: "gcp-project",
		Usage: `The Google project <id> used to validate the identity tokens.
Use the flag multipl etimes to configure multiple projects`,
	}
	removeGCPProjectFlag = cli.StringSliceFlag{
		Name: "remove-gcp-project",
		Usage: `Remove a Google project <id> used to validate the identity tokens.
Use the flag multipl etimes to configure multiple projects`,
	}
	instanceAgeFlag = cli.DurationFlag{
		Name: "instance-age",
		Usage: `The maximum <duration> to grant a certificate in AWS and GCP provisioners.
A <duration> is sequence of decimal numbers, each with optional fraction and a
unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns",
"us" (or "µs"), "ms", "s", "m", "h".`,
	}
	iidRootsFlag = cli.StringFlag{
		Name: "iid-roots",
		Usage: `The <file> containing the certificates used to validate the
instance identity documents in AWS.`,
	}
	disableCustomSANsFlag = cli.BoolFlag{
		Name: "disable-custom-sans",
		Usage: `On cloud provisioners, if enabled only the internal DNS and IP will be added as a SAN.
By default it will accept any SAN in the CSR.`,
	}
	disableTOFUFlag = cli.BoolFlag{
		Name: "disable-trust-on-first-use,disable-tofu",
		Usage: `On cloud provisioners, if enabled multiple sign request for this provisioner
with the same instance will be accepted. By default only the first request
will be accepted.`,
	}

	// Nebula provisioner flags
	nebulaRootFlag = cli.StringFlag{
		Name: "nebula-root",
		Usage: `Root certificate (chain) <file> used to validate the signature on Nebula
provisioning tokens.`,
	}
)

func readNebulaRoots(rootFile string) ([][]byte, error) {
	b, err := utils.ReadFile(rootFile)
	if err != nil {
		return nil, err
	}

	var crt *nebula.NebulaCertificate
	var certs []*nebula.NebulaCertificate
	for len(b) > 0 {
		crt, b, err = nebula.UnmarshalNebulaCertificateFromPEM(b)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading %s", rootFile)
		}
		if crt.Details.IsCA {
			certs = append(certs, crt)
		}
	}
	if len(certs) == 0 {
		return nil, errors.Errorf("error reading %s: no certificates found", rootFile)
	}

	rootBytes := make([][]byte, len(certs))
	for i, crt := range certs {
		b, err = crt.MarshalToPEM()
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling certificate")
		}
		rootBytes[i] = b
	}

	return rootBytes, nil
}
