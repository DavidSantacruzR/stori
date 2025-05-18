
<h1>Stori: Coding Challenge</h1>

**Note: Before running the project please make sure to have set up your aws account with proper**
**permissions to build and deploy the challenge**

**services required:**
* aws-lambda
* aws-container registry
* aws-admin access
* aws-dynamodb
* aws-efs

<h2>Use of AWS EFS as a volume to store csv files:</h2>
Make use of the EFS to load the csv you're planning to act upon.

<h2>Building the lambda images:</h2>

To build the lambda images first use: 
```zsh
sh build.sh
```

<h2>Deploying the services to AWS:</h2>
To deploy the following services using terraform so run the following commands in order:

<h2>Running the tests:</h2>
To run the tests for each module please run ```go test``` inside both folders ```email```
and ```parser```.
