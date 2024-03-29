## Schema Setup
Schema files are not checked in to the repo.

- `make download`to download schema from [AWS docs, us-east-1 region](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-type-schemas.html) and unzip in `cfn/` directory
- `make clean-downloads` to remove downloaded JSON schema files (if necessary)

## Use

Run `bb release`. It does whatever is necessary. Updated `models.rs` and `schema.graphql` are in `configuration` folder. 
For more options have a look at `Makefile`,

## Check results

Have a look at last log file in `logs` directory.
