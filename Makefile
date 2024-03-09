run: main.py
	python main.py

merge-schema:
	cat out/graphql/*.gql > out/graphql/schema.all.gql

show-failures: logs/failures.log
	less logs/failures.log

reset:
	rm -f logs/failures.log
	rm -f out/graphql/*.all.gql

clean-logs:
	rm -f logs/*.log

clean-artifacts:
	rm -f out/graphql/*.gql
	rm -f out/model/*.rs

rerun: reset run
