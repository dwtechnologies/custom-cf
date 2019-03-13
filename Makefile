.PHONY: package-cf deploy-cf build clean

PWD          := $(shell pwd)
PROJECT      ?= custom-cf
OWNER        ?=
ENVIRONMENT  ?= dev
AWS_REGION   ?=
AWS_PROFILE  ?=
S3_BUCKET    ?=
FUNCTIONNAME := $(shell echo $(FUNCTION) | sed -e 's/\//-/g')

# Check vars inside targets by calling "@:$(call check_var, VAR)"
check_var = $(strip $(foreach 1,$1,$(call __check_var,$1,$(strip $(value 2)))))
__check_var = $(if $(value $1),,$(error $1 variable not set))

deploy: build package-cf deploy-cf

package-cf:
	@:$(call check_var, ENVIRONMENT)
	@:$(call check_var, AWS_REGION)
	@:$(call check_var, AWS_PROFILE)
	@:$(call check_var, S3_BUCKET)
	@:$(call check_var, PROJECT)
	@:$(call check_var, FUNCTION)
	@aws cloudformation package \
		--template-file ./$(FUNCTION)/template.yaml \
		--output-template-file ./build/$(FUNCTIONNAME)/template-$(ENVIRONMENT).yaml \
		--profile $(AWS_PROFILE) --region $(AWS_REGION) \
		--s3-bucket $(S3_BUCKET) --s3-prefix lambda/$(PROJECT)/$(FUNCTIONNAME)-$(ENVIRONMENT)

deploy-cf:
	@:$(call check_var, ENVIRONMENT)
	@:$(call check_var, AWS_REGION)
	@:$(call check_var, AWS_PROFILE)
	@:$(call check_var, PROJECT)
	@:$(call check_var, OWNER)
	@:$(call check_var, FUNCTION)
	@aws cloudformation deploy \
		--template-file ./build/$(FUNCTIONNAME)/template-$(ENVIRONMENT).yaml \
		--stack-name $(FUNCTIONNAME)-$(ENVIRONMENT) \
		--profile $(AWS_PROFILE) --region $(AWS_REGION) \
		--capabilities CAPABILITY_NAMED_IAM \
		--no-fail-on-empty-changeset \
		--parameter-overrides \
			Environment=$(ENVIRONMENT) \
			CodeUri=../build/$(FUNCTIONNAME)/handler.zip \
			FunctionName=$(FUNCTIONNAME) \
		--tags \
			Environment=$(ENVIRONMENT) \
			Project=$(PROJECT) \
			Owner=$(OWNER)
	@echo "\nProject successfully deployed"


build:
	@:$(call check_var, PWD)
	@:$(call check_var, FUNCTION)
	@mkdir -p ./build/$(FUNCTIONNAME)
	@docker run --rm \
		-v $(PWD)/build:/build \
		-v $(PWD):/src \
		-w /src \
		golang:1.12.0-stretch \
		sh -c "apt-get update && apt-get install -y zip && \
		cd /src/$(FUNCTION) && go build -o handler && \
		zip handler.zip handler && rm handler && mv handler.zip /build/$(FUNCTIONNAME)"
	@echo "\nProject successfully built"


clean:
	rm -rf build
	