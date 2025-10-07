#!/bin/sh

echo "There are two ways to substitute variables - environment variables or properties. This script uses both"

echo "Note that properties override environment variables"

echo 
echo

echo "RAYGUN_DELEGATE=bob raygun execute -D RAYGUN_NAME=ray"

echo
echo

RAYGUN_DELEGATE=bob raygun execute -D RAYGUN_NAME=ray

