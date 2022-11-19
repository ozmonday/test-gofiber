#!/bin/bash


set -m 
testfiber &
redis-server 

fg %1