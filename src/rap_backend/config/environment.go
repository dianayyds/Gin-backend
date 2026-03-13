package config

import "os"

var PRIEST = os.Getenv("PRIEST")
var CB_HOST = os.Getenv("CB_HOST")
var CB_GROUP = os.Getenv("CB_GROUP")
var HOSTNAME = os.Getenv("HOSTNAME")
var SERVICENAME = os.Getenv("SERVICENAME")
var CUSTOM_RUNTIME_ENV = os.Getenv("CUSTOM_RUNTIME_ENV")
var CONTAINER_HEADER_NAME = os.Getenv("CONTAINER_HEADER_NAME")
