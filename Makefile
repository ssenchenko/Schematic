run: src/main.py
	python src/main.py

clean-logs:
	rm -f logs/*.log

clean-artifacts:
	rm -f out/*.rs

clean-downloads:
	rm -f data/cfn/*.json
	rm -f data/cfn/*.zip

clean: clean-logs clean-artifacts clean-downloads

reset: clean-logs clean-artifacts

rerun: reset run

download:
	cd ./data/cfn ;\
	curl -# -o schema.zip "https://schema.cloudformation.us-east-1.amazonaws.com/CloudformationSchema.zip" ;\
	unzip schema.zip ;\
	rm -f *.zip
