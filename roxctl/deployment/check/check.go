package check

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/retry"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/roxctl/common/environment"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/common/printer"
	"github.com/stackrox/rox/roxctl/common/report"
	"github.com/stackrox/rox/roxctl/summaries/policy"
)

const (
	// Default JSON path expression which retrieves the data from the policy.PolicyJSONResult
	defaultDeploymentCheckJSONPathExpression = "{" +
		"results.#.violatedPolicies.#.name," +
		"results.#.violatedPolicies.#.severity," +
		"results.#.violatedPolicies.#.failingCheck.@boolReplace:{\"true\":\"X\",\"false\":\"-\"}," +
		"results.#.metadata.additionalInfo.name," +
		"results.#.violatedPolicies.#.description," +
		"results.#.violatedPolicies.#.violation.@list," +
		"results.#.violatedPolicies.#.remediation}"
)

var (
	// Default headers to use when printing tabular output
	defaultDeploymentCheckHeaders = []string{
		"POLICY", "SEVERITY", "BREAKS DEPLOY", "DEPLOYMENT", "DESCRIPTION", "VIOLATION", "REMEDIATION",
	}

	// supported output formats with default values
	supportedObjectPrinters = []printer.CustomPrinterFactory{
		printer.NewTabularPrinterFactory(false, defaultDeploymentCheckHeaders,
			defaultDeploymentCheckJSONPathExpression, false, false),
		printer.NewJSONPrinterFactory(false, false),
	}
)

// Command checks the deployment against deploy time system policies
func Command(cliEnvironment environment.Environment) *cobra.Command {
	deploymentCheckCmd := &deploymentCheckCommand{env: cliEnvironment}

	objectPrinterFactory, err := printer.NewObjectPrinterFactory("table", supportedObjectPrinters...)
	// this error should never occur, it would only occur if default values are invalid
	utils.Must(err)

	c := &cobra.Command{
		Use:  "check",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := deploymentCheckCmd.Construct(args, cmd, objectPrinterFactory); err != nil {
				return err
			}
			if err := deploymentCheckCmd.Validate(); err != nil {
				return err
			}

			return deploymentCheckCmd.Check()
		},
	}

	// Add all printer related flags
	objectPrinterFactory.AddFlags(c)

	c.Flags().StringVarP(&deploymentCheckCmd.file, "file", "f", "", "yaml file to send to Central to evaluate policies against")
	c.Flags().BoolVar(&deploymentCheckCmd.json, "json", false, "output policy results as json.")
	c.Flags().IntVarP(&deploymentCheckCmd.retryDelay, "retry-delay", "d", 3, "set time to wait between retries in seconds")
	c.Flags().IntVarP(&deploymentCheckCmd.retryCount, "retries", "r", 3, "Number of retries before exiting as error")
	c.Flags().StringSliceVarP(&deploymentCheckCmd.policyCategories, "categories", "c", nil, "optional comma separated list of policy categories to run.  Defaults to all policy categories.")
	c.Flags().BoolVar(&deploymentCheckCmd.printAllViolations, "print-all-violations", false, "whether to print all violations per alert or truncate violations for readability")
	utils.Must(c.MarkFlagRequired("file"))

	// mark legacy output format specific flags as deprecated
	utils.Must(c.Flags().MarkDeprecated("json", "use the new output format which also offers JSON. NOTE: "+
		"The new output format's structure has changed in a non-backward compatible way."))
	utils.Must(c.Flags().MarkDeprecated("print-all-violations", "use the new output format where all "+
		"violations are printed by default. This flag will only be relevant in combination with the --json flag"))

	return c
}

type deploymentCheckCommand struct {
	// properties bound to cobra flags
	file               string
	json               bool
	retryDelay         int
	retryCount         int
	policyCategories   []string
	printAllViolations bool
	timeout            time.Duration

	// injected or constructed values by Construct
	env                environment.Environment
	printer            printer.ObjectPrinter
	standardizedFormat bool
}

func (d *deploymentCheckCommand) Construct(args []string, cmd *cobra.Command, f *printer.ObjectPrinterFactory) error {
	d.timeout = flags.Timeout(cmd)

	// Only create a printer if legacy json output format is not used
	// TODO(ROX-8303): Remove this once we have fully deprecated the old output format
	if !d.json {
		p, err := f.CreatePrinter()
		if err != nil {
			return err
		}
		d.printer = p
		d.standardizedFormat = f.IsStandardizedFormat()
	}

	return nil
}

func (d *deploymentCheckCommand) Validate() error {
	if _, err := os.Open(d.file); err != nil {
		return errorhelpers.NewErrInvalidArgs(err.Error())
	}

	return nil
}

func (d *deploymentCheckCommand) Check() error {
	err := retry.WithRetry(func() error {
		return d.checkDeployment()
	},
		retry.Tries(d.retryCount+1),
		retry.OnFailedAttempts(func(err error) {
			d.env.Logger().ErrfLn("Scanning image failed: %v. Retrying after %d seconds...",
				err, d.retryDelay)
			time.Sleep(time.Duration(d.retryDelay) * time.Second)
		}))
	if err != nil {
		return err
	}
	return nil
}

func (d *deploymentCheckCommand) checkDeployment() error {
	deploymentFileContents, err := os.ReadFile(d.file)
	if err != nil {
		return err
	}

	alerts, err := d.getAlerts(string(deploymentFileContents))
	if err != nil {
		return err
	}

	return d.printResults(alerts)
}

func (d *deploymentCheckCommand) getAlerts(deploymentYaml string) ([]*storage.Alert, error) {
	conn, err := d.env.GRPCConnection()
	if err != nil {
		return nil, err
	}
	defer utils.IgnoreError(conn.Close)

	svc := v1.NewDetectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	response, err := svc.DetectDeployTimeFromYAML(ctx, &v1.DeployYAMLDetectionRequest{
		Yaml:             deploymentYaml,
		PolicyCategories: d.policyCategories,
	})
	if err != nil {
		return nil, err
	}

	var alerts []*storage.Alert
	for _, r := range response.GetRuns() {
		alerts = append(alerts, r.GetAlerts()...)
	}
	return alerts, nil
}

func (d *deploymentCheckCommand) printResults(alerts []*storage.Alert) error {
	if d.json {
		return report.JSON(d.env.InputOutput().Out, alerts)
	}

	// TODO: Need to refactor this to include additional summary info for non-standardized formats
	// as well as multiple results for each deployment
	policySummary := policy.NewPolicySummaryForPrinting(alerts, storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT)

	if !d.standardizedFormat {
		printDeploymentPolicySummary(policySummary.Summary, d.env.Logger(), policySummary.GetResultNames()...)
	}

	if err := d.printer.Print(policySummary, d.env.InputOutput().Out); err != nil {
		return err
	}

	amountBreakingPolicies := policySummary.GetTotalAmountOfBreakingPolicies()

	if !d.standardizedFormat {
		printAdditionalWarnsAndErrs(policySummary.Summary[policy.TotalPolicyAmountKey], amountBreakingPolicies,
			policySummary.Results, d.env.Logger())
	}

	if amountBreakingPolicies != 0 {
		return policy.NewErrBreakingPolicies(amountBreakingPolicies)
	}
	return nil
}

func printDeploymentPolicySummary(numOfPolicyViolations map[string]int, out environment.Logger, deployments ...string) {
	out.PrintfLn("Policy check results for deployments: %v", deployments)
	out.PrintfLn("(%s: %d, %s: %d, %s: %d, %s: %d, %s: %d)\n",
		policy.TotalPolicyAmountKey, numOfPolicyViolations[policy.TotalPolicyAmountKey],
		policy.LowSeverity, numOfPolicyViolations[policy.LowSeverity.String()],
		policy.MediumSeverity, numOfPolicyViolations[policy.MediumSeverity.String()],
		policy.HighSeverity, numOfPolicyViolations[policy.HighSeverity.String()],
		policy.CriticalSeverity, numOfPolicyViolations[policy.CriticalSeverity.String()])
}

func printAdditionalWarnsAndErrs(amountViolatedPolicies, amountBreakingPolicies int, results []policy.EntityResult,
	out environment.Logger) {
	if amountViolatedPolicies == 0 {
		return
	}

	out.WarnfLn("A total of %d policies have been violated", amountViolatedPolicies)

	if amountBreakingPolicies == 0 {
		return
	}

	out.ErrfLn("%s", policy.NewErrBreakingPolicies(amountBreakingPolicies))

	for _, res := range results {
		for _, breakingPolicy := range res.GetBreakingPolicies() {
			out.ErrfLn("Policy %q within Deployment %q - Possible remediation: %q",
				breakingPolicy.Name, res.Metadata.GetName(), breakingPolicy.Name)
		}
	}
}
