
<h1>Stori: Coding Challenge</h1>

**Note: Before running the project please make sure to have set up your aws account with proper**
**permissions to build and deploy the challenge**

**services required:**
* aws-lambda
* aws-container registry
* aws-admin access
* aws-s3
* aws SES (simple email service)

<h2>Setting up the project</h2>

To build the project, and upload the lambdas to AWS please follow these instructions.

There are four main files to build, and deploy the services:
* ``build_images.sh``
* ``push_images_to_ecr.sh``
* ``build_base_resources.sh``
* ``build_deploy_lambdas.sh``


<h3>1. Building the container images</h3>
In the project root folder first run the sh building file to construct the lambda images:
```zsh
sh build_images.sh
```

<h3>2. Spinning up base resources</h3>
Once the building process finished then go to the **infrastructure** folder and
run the following script to create the base resources to spin up: security groups, roles, SES, ecr, and s3:

```zsh
sh build_base_resources.sh
```

<h3>3. Pushing images to ecr</h3>
Once the previous script is done, now it's the time to push the built images to ecr. For that go back to the root folder
and use the following script.

**IMPORTANT:**
Before running this script check the repository uri and in your `.env` file set a variable ``REGISTRY_URL`` which is the variable
that's going to be used to upload the images to the registry.

**An example**:
``
REGISTRY_URL=123456789.dkr.ecr.us-east-1.amazonaws.com
``

After making sure the .env file exists in the project root folder, and that's the variable has the correct ECR uri, run
the following script.

```zsh
sh push_images_to_ecr.sh
```

<h3>4. Deploying the lambdas to AWS</h3>
At this stage, please very that all resources have been created successfully in your AWS account.

Finally, as a last step, move back to the infrastructure folder and run the last script to deploy the lambdas to the 
cloud provider:
```zsh
sh build_deploy_lambdas.sh
```

<h2>Testing the lambdas in AWS web console</h2>

To test the code, please follow these instructions:

1. In your recently created bucket called ``stori-challenge-david-s``, upload a .csv file according to the 
challenge specifications.
2. Go to the lambda section, and you should see three functions: lambda-parser, lambda-summary, lambda-email.
3. Execute the lambda functions in the following order:
   * lambda-parser
   * lambda-summary
   * lambda-email

<h3>Manual execution</h3>

<h4>Executing lambda-parser:</h4>

The parser function is the first one to be called is this will read the file from s3, parse it, and generate
a detailed transaction output, to be passed onto the next function: lambda-summary.

The function expects an input like the following, the filename, must match the name of the one you've uploaded
on the previous step to the s3 bucket:
```json
{
  "email": "myverifiedemail@gmail.com",
  "sender": "myverifiedemail@gmail.com",
   "filename": "transactions.csv"
}
```

<h4>Executing lambda-summary</h4>

Please use the output of the previous step to run this lambda function. The following function will take a list of
detailed records, and calculate total amounts, averages per transaction type: debit/credit.

The input should look something like this:
```json
{
   "summary":[{
      "transaction_id": 0,
      "month": "July",
      "day": 15,
      "move_type": "credit",
      "amount": 60.5
   }],
   "email": "myverifiedemail@gmail.com",
   "sender": "myverifiedemail@gmail.com",
   "filename": "transactions.csv"
}
```

<h4>Executing lambda-email</h4>
As stated in previous steps, please use the output from lambda-summary for this specific function.

This specific function receives the summarised data, parse it into an HTML template, download the file from s3, and
include it in the email as an attachment.

The expected input should look something like following:

```json
{
"summary": {
   "total_balance": 60.5,
   "average_credit_amount": 60.5,
   "average_debit_amount": 0,
   "monthly_summary": [
      {
         "month": "July",
         "number_of_transactions": 1
      }
   ]
},
"email": "myverifiedemail@gmail.com",
"sender": "myverifiedemail@gmail.com",
"filename": "transactions.csv"
}
```

Once that's done you'll see a confirmation that an email was sent to the email address defined at the beginning
of the flow, and an email to said address, that should look like the following:

![Expected email](https://github.com/DavidSantacruzR/stori/blob/master/output_stori.jpeg)

<h3>Alternative Execution Method</h3>
As an alternative and more automated execution method you can use a step function called
**lambda-step-function**. As usual run it from the aws console.

The expected input follows the same contract as in manual execution:

```json
{
   "email": "myverifiedemail@gmail.com",
   "sender": "myverifiedemail@gmail.com",
   "filename": "transactions.csv"
}
```

For execution remember that the email account you use has to be verified, and the file uploaded in the bucket
as described in previous steps.

Once the execution is finished, you should see an output like the following:

![Execution Runner]()


<h2>Running unit tests:</h2>
To run the unit tests for each module please run ``go test`` inside both folders ``summary``
and `parser``.

<h2>Project explanation</h2>

In this section there is a brief explanation of the code, the technical decisions, implementations, 
and areas for improvement.

<h3>Technical decisions</h3>

The decision to use lambdas derives from the ease of use, configuration, and that it's easier to map a responsibility to
each one of them, controlling the execution flow, verifying easily the output at every stage. 

The use of docker allow to
seamlessly build images, and integrated them with ecr to build the lambdas, instead of using .zip files. 

The use of terraform helps in making sure that every step needed for execution is always the same, and don't get lost in
cloud configuration.

The use of go workspaces derives from the need to use docker, and have the code clearly separated for each build, and
it was easier to automate with a script.

<h3>Implementations</h3>

As mentioned earlier, the project is separated into workspaces: parser, summary, and email. Every function makes use of
the AWS SDK to handle interactions between different services.

In the parser module, you’ll find two Go files — the `main` entrypoint and a `utils` file — along with two CSV 
files used for unit testing.


The primary objective of this lambda function is to read the csv file from a s3 bucket and transform it into
a list of transactions following the DTO:
```golang
type Transaction struct {
Move          string  `json:"move_type"`
TransactionId int     `json:"transaction_id"`
Month         string  `json:"month"`
Day           int     `json:"day"`
Amount        float64 `json:"amount"`
}
```
Splitting transactions into debit and credit types enables simpler grouping and transformation logic downstream.

This parsed structure ensures the summary Lambda’s complexity remains O(n), avoiding multiple passes 
(e.g., one for debits, one for credits, one for monthly grouping).

In the summary function two primary DTOs are implemented, one to hold all the summary information related to the account
and one to hold the summarised data from each month.
```golang
type AccountSummary struct {
	TotalBalance        float64        `json:"total_balance"`
	AverageCreditAmount float64        `json:"average_credit_amount"`
	AverageDebitAmount  float64        `json:"average_debit_amount"`
	Transactions        []MonthSummary `json:"monthly_summary"`
}

type MonthSummary struct {
	Month                string `json:"month"`
	NumberOfTransactions int    `json:"number_of_transactions"`
}
```

For the email lambda function there are two important pieces related to the implementation, and are that the
template is defined as a constant in the module.

This tight coupling simplifies rendering but trades off flexibility — 
template customization requires a code change and redeployment.

The HTML email template is embedded as a string constant in the module.
It uses Go’s `html/template` package and follows the `AccountSummary` structure:
```golang
const TEMPLATE = `
<h1 style="color:darkgreen;">Transaction Summary</h1>
<br/>
<div>
	<ul>
		<li><strong>Total Balance:</strong> {{.TotalBalance}}</li>
		<li><strong>Average Debit:</strong> {{.AverageDebitAmount}}</li>
		<li><strong>Average Credit:</strong> {{.AverageCreditAmount}}</li>
	</ul>
	<h3>Monthly Summary:</h3>
	<ul>
	{{range .Transactions}}
		<li><strong>{{.Month}}:</strong> {{.NumberOfTransactions}} transactions</li>
	{{end}}
	</ul>
</div>
<img src="https://www.storicard.com/_next/static/media/stori_s_color.90dc745f.svg" alt="Stori Logo"/>
`
```

The second important piece is the use of the AWS **SendRawEmail** functionality. 
The function reads the original CSV file from S3, encodes it in base64, and attaches it to the outgoing 
email using the SendRawEmail SES API.

<h3>Areas for improvement</h3>

From the challenge, I've identified some areas for improvement, outlined below.

* Clear area for improvement in the current project is the tight coupling between lambda contracts. For the challenge
the easiest way to pass on information related to the sender, and the name attached file (instead of sending a key with base64 info)
so any lambda could access that resource.
* roles, and permissions, each lambda has access to all resources including SES, and S3. Separating also roles according
to responsibilities is ideal.
* For easier access to resources, and specifically to have a single entrypoint for execution using an AWS api-gateway
allows for a single http request calling the lambda-parser function, and executing all lambdas in steps.
* Implement the storage function using a dynamodb instance.
* Implement lambda versioning.