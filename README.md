## Project Status

IA resources rust models for async-graphql library generated, compile, and can be used to generate graphql schema.

Results are in `configuration` folder, which will make them accessible in build time if this package is taken as a build-time dependency.

CFN schema downloads locally only from [AWS docs web-site, IAD region](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-type-schemas.html). 
Thus, no remote and dry-run builds are possible now without checking schema in to the repository. 
Considering size of CFN schema (~10M) checking it in is not wise. In the future schema package maintained by AppComposer team will be taken as a dependency.

## Project Set Up

No extra dependencies are necessary. You can run it with system Python interpreter, with Brazil 
or with a local Python virtual environment. Here is how to do it with the last option:

- Create Python virtual environment `python3 -m venv .venv`
- In your terminal session activate virtual environment `source .venv/bin/activate`

Tests for the codebase should have been run manually now (has to be fixed). 
Compilation of rust library and generation of graphql sdl (schema definition language) is automated. 

## Schema Setup
Schema files are not checked in to the repo.

- `make download`to download schema from [AWS docs, us-east-1 region](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-type-schemas.html) and unzip in `cfn/` directory
- `make clean-downloads` to remove downloaded JSON schema files (if necessary)

## Use

Run `bb release`. It does whatever is necessary. Updated `models.rs` and `schema.graphql` are in `configuration` folder. 
For more options have a look at `Makefile`,

## Check results

Have a look at last log file in `logs` directory.

