#!/bin/bash

consul kv put gravity/config @"$(dirname $0)/config.hcl"
