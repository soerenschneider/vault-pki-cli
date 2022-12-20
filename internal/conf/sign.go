package conf

type SignArguments struct {
}

func (c *SignArguments) Validate() []error {
	errs := make([]error, 0)

	/*
		if len(c.CertificateFile) == 0 {
			errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CERTIFICATE_FILE))
		}

		if len(c.CsrFile) == 0 {
			errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_CSR_FILE))
		}

		if len(c.CommonName) == 0 {
			errs = append(errs, fmt.Errorf("empty '%s' provided", FLAG_ISSUE_COMMON_NAME))
		}

		ownerDefined := len(c.FileOwner) > 0
		groupDefined := len(c.FileGroup) > 0
		if !ownerDefined && groupDefined {
			errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_GROUP, FLAG_FILE_OWNER))
		}
		if ownerDefined && !groupDefined {
			errs = append(errs, fmt.Errorf("only '%s' defined but not '%s'", FLAG_FILE_OWNER, FLAG_FILE_GROUP))
		}

	*/

	return errs
}
