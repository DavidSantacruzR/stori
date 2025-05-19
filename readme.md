
<h1>Stori: Coding Challenge</h1>

**Note: Before running the project please make sure to have set up your aws account with proper**
**permissions to build and deploy the challenge**

**services required:**
* aws-lambda
* aws-container registry
* aws-admin access
* aws-dynamodb
* aws-efs

<h2>Spinning up the initial infra configuration</h2>

<h3>1. Building the container images</h3>
In the root folder first run the building file to construct the lambda images:
```zsh
sh build_images.sh
```

<h3>2. Spinning up base resources</h3>
Once the building process finished then go to the **infrastructure** folder and
run the following script to create the base resources to use: dynamo, SES, ecr, and efs:

```zsh
sh build_base_resources.sh
```

<h3>3. Pushing images to ecr</h3>
Once the previous script is done, now it's the time to push the built images to ecr. For that go back to the root folder
and use the following script.

**IMPORTANT:**
Before running this script check the repository uri and export the variable as ``REGISTRY_URL`` which is the variable
that's going to be used to upload the images to the registry.

**An example**: 123456789.dkr.ecr.us-east-1.amazonaws.com

```zsh
sh push_images_to_ecr.sh
```

<h3>4. Deploying the lambdas to AWS</h3>
Finally, as a last step, move back to the infrastructure folder and run the last script to deploy the lambdas to the cloud provider:
```zsh
sh build_deploy_lambdas.sh
```

<h2>Running the tests:</h2>
To run the tests for each module please run ```go test``` inside both folders ```email```
and ```parser```.
