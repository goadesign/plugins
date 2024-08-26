#! /usr/bin/make
#
# Makefile for goa v3 arnz plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
PLUGIN_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
include ../plugins.mk

gen:
	goa gen goa.design/plugins/v3/arnz/test/design -o "$(PLUGIN_DIR)/test"

build-examples:

clean: