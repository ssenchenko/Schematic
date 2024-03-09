# Project Status
## Schema from [AWS docs web-site, IAD region](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-type-schemas.html)

- ### Files (resources) - 1206
- ### Size - 9.8M

## Progress

- ### Currently processed - 811 or 2/3 of resources
- ### 1/3 or 395 resources will be supported soon

## Schema size (67% of resources)

- ### With filters - 1.2M
- ### Without filters - 618K

## Full schema size (rough estimation)
- ### With filters - 1.8M
- ### Without filters - 922K

# Set Up

## Python
No extra dependencies are necessary. You can run it with system Python interpreter, with Brazil 
or with a local Python virtual environment. Here is how to do it with the last option:

- Create Python virtual environment `python3 -m venv .venv`
- In your terminal session activate virtual environment `source .venv/bin/activate`

## Schema
Schema files are not checked in to the repo

- Download schema from [AWS docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/resource-type-schemas.html)
- Unzip and move schema files to the `cfn/` directory

# Use

## Run
- `make run`
  - first it runs on all files in `cfn/` directory
  - on the next run, if there were failed files from the last time, script will run only on them
- `make rerun` to run again on all files from `cfn/` even if there are errors from the last run
- `make merge-schema` to merge schema in a single file `out/graphql/schema.all.gql`
- `make clean-logs` to remove log files
- `make clean-artifacts` to remove generated types

## Check results
Generated types are not checked in to the repo

- Rust types are in `out/model/`
- GraphQL files are in `out/graphql`
- Files which failed to mapping are listed in `logs/out/failures.log`
- To see how many files are mapped successfully `ls out/graphql/*.gql | wc -l`
- To see how many files failed mapping `wc -l < logs/failures.log`
- To see mapping logs, check the log file with the latest time in the file name in `logs`
