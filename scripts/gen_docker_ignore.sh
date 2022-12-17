#!/bin/bash

cp .gitignore ./.dockerignore
awk '{print "client/" $0}' client/.gitignore >> ./.dockerignore;
awk '{print "rewinder/" $0}' rewinder/.gitignore >> ./.dockerignore;
