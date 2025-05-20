
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


<h2>Running unit tests:</h2>
To run the unit tests for each module please run ```go test``` inside both folders ```summary```
and ```parser```.

<h2>Executing the project</h2>

After finishing the configuration file, test the solution accordingly.

1. Upload a csv file to the s3 bucket: stori-challenge-david-s.
2. invoke the lambda on the aws web console indicating the destination email, and the name of the csv
file you just uploaded.
3. wait for the email.

<h2>Explanation of the solution</h2>

For the challenge, and the technical requirements to fulfill I determined that the best
approximation was to use aws-lambdas, and step functions accordingly passing each output from the
function flow directly to the next function until the email is sent to the user.

The business logic is structured in four folders:

* email
* parser
* storage
* summary
