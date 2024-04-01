// Copyright 2022 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"crypto/x509/pkix"
	"fmt"
	"net/url"
	"reflect"

	"github.com/gobwas/glob"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
	"github.com/sigstore/fulcio/pkg/certificate"
)

const (
	AllowAllConstraint = "*"
)

// +kubebuilder:object:generate=true
type CertConstraint struct {
	CommonName    string                 `json:"commonname"`
	DNSNames      []string               `json:"dnsnames"`
	Emails        []string               `json:"emails"`
	Organizations []string               `json:"organizations"`
	URIs          []string               `json:"uris"`
	Roots         []string               `json:"roots"`
	Extensions    certificate.Extensions `json:"extensions"`
}

func (cc CertConstraint) Check(verifier *cryptoutil.X509Verifier, trustBundles map[string]TrustBundle) error {
	errs := make([]error, 0)
	cert := verifier.Certificate()

	if err := checkCertConstraint("common name", []string{cc.CommonName}, []string{cert.Subject.CommonName}); err != nil {
		errs = append(errs, err)
	}

	if err := checkCertConstraint("dns name", cc.DNSNames, cert.DNSNames); err != nil {
		errs = append(errs, err)
	}

	if err := checkCertConstraint("email", cc.Emails, cert.EmailAddresses); err != nil {
		errs = append(errs, err)
	}

	if err := checkCertConstraint("organization", cc.Organizations, cert.Subject.Organization); err != nil {
		errs = append(errs, err)
	}

	if err := checkCertConstraint("uri", cc.URIs, urisToStrings(cert.URIs)); err != nil {
		errs = append(errs, err)
	}

	if err := cc.checkTrustBundles(verifier, trustBundles); err != nil {
		errs = append(errs, err)
	}

	if err := cc.checkExtensions(cert.Extensions); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return ErrConstraintCheckFailed{errs}
	}

	return nil
}

func (cc CertConstraint) checkTrustBundles(verifier *cryptoutil.X509Verifier, trustBundles map[string]TrustBundle) error {
	if len(cc.Roots) == 1 && cc.Roots[0] == AllowAllConstraint {
		for _, bundle := range trustBundles {
			if err := verifier.BelongsToRoot(bundle.Root); err == nil {
				return nil
			}
		}
	} else {
		for _, rootID := range cc.Roots {
			if bundle, ok := trustBundles[rootID]; ok {
				if err := verifier.BelongsToRoot(bundle.Root); err == nil {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("cert doesn't belong to any root specified by constraint %+q", cc.Roots)
}

func (cc CertConstraint) checkExtensions(ext []pkix.Extension) error {
	extensions, err := certificate.ParseExtensions(ext)
	if err != nil {
		return fmt.Errorf("error parsing fulcio cert extensions: %w", err)
	}

	fields := reflect.VisibleFields(reflect.TypeOf(cc.Extensions))
	for _, field := range fields {
		constraintField := reflect.ValueOf(cc.Extensions).FieldByName(field.Name)
		if constraintField.String() == "" {
			log.Infof("No constraint for field %s, allowing all values", field.Name)
			continue
		}
		extensionsField := reflect.ValueOf(extensions).FieldByName(field.Name)

		fieldGlob := glob.MustCompile(constraintField.String())
		if !fieldGlob.Match(extensionsField.String()) {
			return fmt.Errorf("cert field %s doesn't match constraint %+q", field.Name, constraintField.String())
		}
	}

	return nil
}

func urisToStrings(uris []*url.URL) []string {
	res := make([]string, 0)
	for _, uri := range uris {
		res = append(res, uri.String())
	}

	return res
}

func checkCertConstraint(attribute string, constraints, values []string) error {
	// If our only constraint is the AllowAllConstraint it's a pass
	if len(constraints) == 1 && constraints[0] == AllowAllConstraint {
		return nil
	}

	// treat a single empty string the same as a constraint on an empty attribute
	if len(constraints) == 1 && constraints[0] == "" {
		constraints = []string{}
	}

	if len(values) == 1 && values[0] == "" {
		values = []string{}
	}

	if len(constraints) == 0 && len(values) > 0 {
		return fmt.Errorf("not expecting any %s(s), but cert has %d %s(s)", attribute, len(values), attribute)
	}

	unmet := make(map[string]struct{})
	for _, constraint := range constraints {
		unmet[constraint] = struct{}{}
	}

	for _, value := range values {
		if _, ok := unmet[value]; !ok {
			return fmt.Errorf("cert has an unexpected %s %s given constraints %+q", attribute, value, constraints)
		}

		delete(unmet, value)
	}

	if len(unmet) > 0 {
		return fmt.Errorf("cert with %s(s) %+qDid not pass all constraints %+q", attribute, values, constraints)
	}

	return nil
}
