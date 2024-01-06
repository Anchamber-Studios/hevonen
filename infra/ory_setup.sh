#! /bin/bash
source .env
echo "Setup for project $ORY_PROJECT_ID"

JWKS_FILE=.es256.jwks.json
if [ -f "$JWKS_FILE" ]; then
    echo "jwks exists"
else 
    echo "jwks does not exist exists. generate a new one"
    ory create jwk some-example-set \
  		--alg ES256 --project $ORY_PROJECT_ID --format json-pretty \
  		> $JWKS_FILE
fi

JWKS=$(cat $JWKS_FILE | base64 | tr -d '\n')
CLAIMS=$(cat claims.jsonnet | base64 | tr -d '\n')
ory patch identity-config $ORY_PROJECT_ID \
	--add "/session/whoami/tokenizer/templates/jwt_example_template={\"jwks_url\":\"base64://$JWKS\",\"claims_mapper_url\":\"base64://$CLAIMS\",\"ttl\":\"10m\"}" \
	--format yaml

