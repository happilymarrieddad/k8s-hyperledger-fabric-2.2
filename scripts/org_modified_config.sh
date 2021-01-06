#!/bin/bash

jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"hp":.[1]}}}}}' config.json ./channels/hp.json > modified_config.json
