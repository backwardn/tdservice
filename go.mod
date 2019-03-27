module intel/isecl/tdservice

require (
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/sirupsen/logrus v1.3.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190219172222-a4c6cb3142f2
	gopkg.in/yaml.v2 v2.2.2
	intel/isecl/lib/common v0.0.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
)

replace intel/isecl/lib/common => gitlab.devtools.intel.com/sst/isecl/lib/common v0.0.0-20190319012411-272076869a86
