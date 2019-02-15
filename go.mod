module intel/isecl/threat-detection-service

require (
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/sirupsen/logrus v1.3.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/crypto v0.0.0-20180904163835-0709b304e793
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/common v0.0.0
)

replace intel/isecl/lib/common => gitlab.devtools.intel.com/sst/isecl/lib/common.git v0.0.0-20190208035330-09f2616d9eb0
