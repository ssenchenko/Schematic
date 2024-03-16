run: src/main.py
	python src/main.py

merge-schema:
	cat out/graphql/*.gql > out/graphql/schema.all.gql

show-failures: logs/failures.log
	less logs/failures.log

clean-logs:
	rm -f logs/*.log

clean-artifacts:
	rm -f out/graphql/*.gql
	rm -f out/model/*.rs
	rm -f out/map/*.json

clean-downloads:
	rm -f cfn/*.json
	rm -f cfn/*.zip

clean: clean-logs clean-artifacts clean-downloads

reset: clean-logs clean-artifacts

rerun: reset run

download:
	cd ./cfn ;\
	curl -# -o schema.zip "https://schema.cloudformation.us-east-1.amazonaws.com/CloudformationSchema.zip" ;\
	unzip schema.zip ;\
	rm -f *.zip
