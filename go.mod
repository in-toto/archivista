module github.com/testifysec/go-witness

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/digitorus/pkcs7 v0.0.0-20230220124406-51331ccfc40f
	github.com/digitorus/timestamp v0.0.0-20230220124323-d542479a2425
	github.com/edwarnicke/gitoid v0.0.0-20220710194850-1be5bfda1f9d
	github.com/go-git/go-git/v5 v5.5.2
	github.com/mattn/go-isatty v0.0.17
	github.com/open-policy-agent/opa v0.49.1
	github.com/owenrumney/go-sarif v1.1.1
	github.com/spiffe/go-spiffe/v2 v2.1.2
	github.com/stretchr/testify v1.8.1
	github.com/testifysec/archivista-api v0.0.0-20230220215059-632b84b82b76
	go.step.sm/crypto v0.25.0
	golang.org/x/sys v0.5.0
	google.golang.org/grpc v1.53.0
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cloudflare/circl v1.3.2 // indirect
	github.com/coreos/go-oidc/v3 v3.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/google/go-containerregistry v0.13.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/letsencrypt/boulder v0.0.0-20221109233200-85aa52084eaf // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pjbgf/sha1cd v0.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	github.com/skeema/knownhosts v1.1.0 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/titanous/rocacheck v0.0.0-20171023193734-afe73141d399 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/oauth2 v0.5.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/aws/aws-sdk-go v1.44.207
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.4.1 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/sigstore/fulcio v1.1.0
	github.com/sigstore/sigstore v1.5.1
	github.com/theupdateframework/go-tuf v0.5.2-0.20220930112810-3890c1e7ace4 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	github.com/zeebo/errs v1.3.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230222225845-10f96fb3dbec // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sigstore/rekor => github.com/testifysec/rekor v0.4.0-dsse-intermediates-2

replace github.com/gin-gonic/gin v1.5.0 => github.com/gin-gonic/gin v1.7.7

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.3-0.20220303224323-02efb9a75ee1
