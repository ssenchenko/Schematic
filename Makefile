test-types: out/model.rs test/src/tools/sdl.rs
	cp out/model.rs test/src/lib.rs ;\
	cd test/ ;\
	cargo build ;\
	cargo run --bin generate-sdl > schema.graphql

run: src/main.py
	python src/main.py

clean-logs:
	rm -f logs/*.log

clean-artifacts:
	rm -f out/*.rs

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
