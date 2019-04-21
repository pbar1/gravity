#!/bin/bash

consul kv put gravity/config @"$(dirname $0)/consul-config.hcl"
