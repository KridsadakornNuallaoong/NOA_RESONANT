package sensitive

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/pquerna/otp"
)

// GenerateOpts provides options for Generate().  The default values
// are compatible with Google-Authenticator.
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint
	// Size in size of the generated Secret. Defaults to 20 bytes.
	SecretSize uint
	// Secret to store. Defaults to a randomly generated secret of SecretSize.  You should generally leave this empty.
	Secret []byte
	// Digits to request. Defaults to 6.
	Digits otp.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm otp.Algorithm
	// Reader to use for generating TOTP Key.
	Rand io.Reader
}

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

// Generate a new TOTP Key.
func Generate(opts GenerateOpts) (*otp.Key, error) {
	// url encode the Issuer/AccountName
	if opts.Issuer == "" {
		return nil, fmt.Errorf("issuer is required")
	}

	if opts.AccountName == "" {
		return nil, fmt.Errorf("account name is required")
	}

	if opts.Period == 0 {
		opts.Period = 30
	}

	if opts.SecretSize == 0 {
		opts.SecretSize = 20
	}

	if opts.Digits == 0 {
		opts.Digits = otp.DigitsSix
	}

	if opts.Rand == nil {
		opts.Rand = rand.Reader
	}

	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example

	v := url.Values{}
	if len(opts.Secret) != 0 {
		v.Set("secret", b32NoPadding.EncodeToString(opts.Secret))
	} else {
		secret := make([]byte, opts.SecretSize)
		if _, err := io.ReadFull(opts.Rand, secret); err != nil {
			return nil, err
		}
		v.Set("secret", b32NoPadding.EncodeToString(secret))
	}

	v.Set("issuer", opts.Issuer)
	v.Set("period", strconv.FormatUint(uint64(opts.Period), 10))
	v.Set("algorithm", opts.Algorithm.String())
	v.Set("digits", opts.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + url.PathEscape(opts.Issuer) + ":" + url.PathEscape(opts.AccountName),
		RawQuery: v.Encode(),
	}

	// display everything in the console for now
	fmt.Println("URL:", u.String())
	fmt.Println("Secret:", v.Get("secret"))

	return otp.NewKeyFromURL(u.String())
}
