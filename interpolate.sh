#!/bin/bash

sed -i 's#PEGASUS_AWX_DOCKER_REGISTRY_AUTH_SECRETS_JINJA2_VAR#PEGASUS_AWX_DOCKER_REGISTRY_AUTH_SECRETS#g' ./.env
