// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcrypto

import (
	"hash"
	"math/big"
	"testing"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
)

type ecdsaFixture struct {
	name    string
	key     *ecdsaKey
	alg     func() hash.Hash
	message string
	r, s    string
}

type ecdsaKey struct {
	key      *ecdsa.PrivateKey
	subgroup int
}

var p224 = &ecdsaKey{
	key: &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P224(),
			X:     ecdsaLoadInt("00CF08DA5AD719E42707FA431292DEA11244D64FC51610D94B130D6C"),
			Y:     ecdsaLoadInt("EEAB6F3DEBE455E3DBF85416F7030CBD94F34F2D6F232C69F3C1385A"),
		},
		D: ecdsaLoadInt("F220266E1105BFE3083E03EC7A3A654651F45E37167E88600BF257C1"),
	},
	subgroup: 224,
}

var p256 = &ecdsaKey{
	key: &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     ecdsaLoadInt("60FED4BA255A9D31C961EB74C6356D68C049B8923B61FA6CE669622E60F29FB6"),
			Y:     ecdsaLoadInt("7903FE1008B8BC99A41AE9E95628BC64F2F1B20C2D7E9F5177A3C294D4462299"),
		},
		D: ecdsaLoadInt("C9AFA9D845BA75166B5C215767B1D6934E50C3DB36E89B127B8A622B120F6721"),
	},
	subgroup: 256,
}

var p384 = &ecdsaKey{
	key: &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P384(),
			X:     ecdsaLoadInt("EC3A4E415B4E19A4568618029F427FA5DA9A8BC4AE92E02E06AAE5286B300C64DEF8F0EA9055866064A254515480BC13"),
			Y:     ecdsaLoadInt("8015D9B72D7D57244EA8EF9AC0C621896708A59367F9DFB9F54CA84B3F1C9DB1288B231C3AE0D4FE7344FD2533264720"),
		},
		D: ecdsaLoadInt("6B9D3DAD2E1B8C1C05B19875B6659F4DE23C3B667BF297BA9AA47740787137D896D5724E4C70A825F872C9EA60D2EDF5"),
	},
	subgroup: 384,
}

var p521 = &ecdsaKey{
	key: &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P521(),
			X:     ecdsaLoadInt("1894550D0785932E00EAA23B694F213F8C3121F86DC97A04E5A7167DB4E5BCD371123D46E45DB6B5D5370A7F20FB633155D38FFA16D2BD761DCAC474B9A2F5023A4"),
			Y:     ecdsaLoadInt("0493101C962CD4D2FDDF782285E64584139C2F91B47F87FF82354D6630F746A28A0DB25741B5B34A828008B22ACC23F924FAAFBD4D33F81EA66956DFEAA2BFDFCF5"),
		},
		D: ecdsaLoadInt("0FAD06DAA62BA3B25D2FB40133DA757205DE67F5BB0018FEE8C86E1B68C7E75CAA896EB32F1F47C70855836A6D16FCC1466F6D8FBEC67DB89EC0C08B0E996B83538"),
	},
	subgroup: 521,
}

var fixtures = []ecdsaFixture{
	// ECDSA, 224 Bits (Prime Field)
	// https://tools.ietf.org/html/rfc6979#appendix-A.2.4
	{
		name:    "P224/SHA-256 #1",
		key:     p224,
		alg:     sha256.New,
		message: "sample",
		r:       "61AA3DA010E8E8406C656BC477A7A7189895E7E840CDFE8FF42307BA",
		s:       "437EBFAF254A2DC88F786B6B061E7022049DF927FA4B6B5EC9AB293C",
	},
	{
		name:    "P224/SHA-512 #1",
		key:     p224,
		alg:     sha512.New,
		message: "sample",
		r:       "074BD1D979D5F32BF958DDC61E4FB4872ADCAFEB2256497CDAC30397",
		s:       "5B3135E693C2A5E00CEFD84CCE793A13FC79C7B2F231F516FECD79B9",
	},
	{
		name:    "P224/SHA-256 #2",
		key:     p224,
		alg:     sha256.New,
		message: "test",
		r:       "AD04DDE87B84747A243A631EA47A1BA6D1FAA059149AD2440DE6FBA6",
		s:       "178D49B1AE90E3D8B629BE3DB5683915F4E8C99FDF6E666CF37ADCFD",
	},
	{
		name:    "P224/SHA-512 #2",
		key:     p224,
		alg:     sha512.New,
		message: "test",
		r:       "049F050477C5ADD858CAC56208394B5A55BAEBBE887FDF765047C17C",
		s:       "077EB13E7005929CEFA3CD0403C7CDCC077ADF4E44F3C41B2F60ECFF",
	},
	{
		name:    "P256/SHA-256 #1",
		key:     p256,
		alg:     sha256.New,
		message: "sample",
		r:       "EFD48B2AACB6A8FD1140DD9CD45E81D69D2C877B56AAF991C34D0EA84EAF3716",
		s:       "834E36AD29A83BF2BC9385E491D6099C8FDF9D1ED67AA7EA5F51F93782857A9",
	},
	{
		name:    "P256/SHA-384 #1",
		key:     p256,
		alg:     sha512.New384,
		message: "sample",
		r:       "0EAFEA039B20E9B42309FB1D89E213057CBF973DC0CFC8F129EDDDC800EF7719",
		s:       "4861F0491E6998B9455193E34E7B0D284DDD7149A74B95B9261F13ABDE940954",
	},
	{
		name:    "P256/SHA-256 #2",
		key:     p256,
		alg:     sha256.New,
		message: "test",
		r:       "F1ABB023518351CD71D881567B1EA663ED3EFCF6C5132B354F28D3B0B7D38367",
		s:       "019F4113742A2B14BD25926B49C649155F267E60D3814B4C0CC84250E46F0083",
	},
	// ECDSA, 384 Bits (Prime Field)
	// https://tools.ietf.org/html/rfc6979#appendix-A.2.6
	{
		name:    "P384/SHA-256 #1",
		key:     p384,
		alg:     sha256.New,
		message: "sample",
		r:       "21B13D1E013C7FA1392D03C5F99AF8B30C570C6F98D4EA8E354B63A21D3DAA33BDE1E888E63355D92FA2B3C36D8FB2CD",
		s:       "C55BBC04EF88BA40B428834C76E98B9CDF975EF35981C2B69B127124C652EF3683DA9C57B95E34C2C20930227CB1EC3",
	},
	{
		name:    "P384/SHA-256 #2",
		key:     p384,
		alg:     sha256.New,
		message: "test",
		r:       "6D6DEFAC9AB64DABAFE36C6BF510352A4CC27001263638E5B16D9BB51D451559F918EEDAF2293BE5B475CC8F0188636B",
		s:       "2D46F3BECBCC523D5F1A1256BF0C9B024D879BA9E838144C8BA6BAEB4B53B47D51AB373F9845C0514EEFB14024787265",
	},
	// ECDSA, 521 Bits (Prime Field)
	// https://tools.ietf.org/html/rfc6979#appendix-A.2.7
	{
		name:    "P521/SHA-224 #1",
		key:     p521,
		alg:     sha256.New224,
		message: "sample",
		r:       "1776331CFCDF927D666E032E00CF776187BC9FDD8E69D0DABB4109FFE1B5E2A30715F4CC923A4A5E94D2503E9ACFED92857B7F31D7152E0F8C00C15FF3D87E2ED2E",
		s:       "050CB5265417FE2320BBB5A122B8E1A32BD699089851128E360E620A30C7E17BA41A666AF126CE100E5799B153B60528D5300D08489CA9178FB610A2006C254B41F",
	},
	{
		name:    "P521/SHA-256 #1",
		key:     p521,
		alg:     sha256.New,
		message: "sample",
		r:       "1511BB4D675114FE266FC4372B87682BAECC01D3CC62CF2303C92B3526012659D16876E25C7C1E57648F23B73564D67F61C6F14D527D54972810421E7D87589E1A7",
		s:       "04A171143A83163D6DF460AAF61522695F207A58B95C0644D87E52AA1A347916E4F7A72930B1BC06DBE22CE3F58264AFD23704CBB63B29B931F7DE6C9D949A7ECFC",
	},
	{
		name:    "P521/SHA-256 #2",
		key:     p521,
		alg:     sha256.New,
		message: "test",
		r:       "00E871C4A14F993C6C7369501900C4BC1E9C7B0B4BA44E04868B30B41D8071042EB28C4C250411D0CE08CD197E4188EA4876F279F90B3D8D74A3C76E6F1E4656AA8",
		s:       "0CD52DBAA33B063C3A6CD8058A1FB0A46A4754B034FCC644766CA14DA8CA5CA9FDE00E88C1AD60CCBA759025299079D7A427EC3CC5B619BFBC828E7769BCD694E86",
	},
}

func TestECDSA(t *testing.T) {
	for _, f := range fixtures {
		testEcdsaFixture(&f, t)
	}
}

func ecdsaLoadInt(s string) (n *big.Int) {
	n, _ = new(big.Int).SetString(s, 16)
	return
}

func testEcdsaFixture(f *ecdsaFixture, t *testing.T) {
	t.Logf("Testing %s", f.name)

	h := f.alg()
	h.Write([]byte(f.message))
	digest := h.Sum(nil)

	g := f.key.subgroup / 8
	if len(digest) > g {
		digest = digest[0:g]
	}

	r, s, err := EcdsaSign(f.key.key, digest, f.alg)
	if err != nil {
		t.Error(err)
		return
	}

	expectedR := ecdsaLoadInt(f.r)
	expectedS := ecdsaLoadInt(f.s)

	if r.Cmp(expectedR) != 0 {
		t.Errorf("%s: Expected R of %X, got %X", f.name, expectedR, r)
	}

	if s.Cmp(expectedS) != 0 {
		t.Errorf("%s: Expected S of %X, got %X", f.name, expectedS, s)
	}
}

// BenchmarkEcdsaSign-4     	   10000	    150377 ns/op
func BenchmarkSignP256(b *testing.B) {
	f := fixtures[4]
	h := f.alg()
	h.Write([]byte(f.message))
	digest := h.Sum(nil)

	g := f.key.subgroup / 8
	if len(digest) > g {
		digest = digest[0:g]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := EcdsaSign(f.key.key, digest, f.alg); err != nil {
			panic(err)
		}
	}
}

// BenchmarkEcdsaVerify-4   	   10000	    154824 ns/op
func BenchmarkVerifyP256(b *testing.B) {
	f := fixtures[4]
	h := f.alg()
	h.Write([]byte(f.message))
	digest := h.Sum(nil)

	g := f.key.subgroup / 8
	if len(digest) > g {
		digest = digest[0:g]
	}
	r, s, _ := EcdsaSign(f.key.key, digest, f.alg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !ecdsa.Verify(&f.key.key.PublicKey, digest, r, s) {
			b.Fatal("failed to verify signature")
		}
	}
}
