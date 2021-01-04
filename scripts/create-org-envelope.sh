#!/bin/bash

echo '{"payload":{"header":{"channel_header":{"channel_id":"mainchannel", "type":2}},"data":{"config_update":'$(cat org_update.json)'}}}' | jq . > org_update_in_envelope.json
