#!/bin/sh

echo "The fourth test (PUT) should fail, because the rego changed to support PUT after this decision was made"

raygun backtest --opa-bundle-url=example6-bundle.tar.gz backtest.json


