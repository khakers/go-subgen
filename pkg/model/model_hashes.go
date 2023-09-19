package model

import "errors"

func GetModelSha1Hash(model Model) (string, error) {
	switch model {
	case Tiny:
		return "bd577a113a864445d4c299885e0cb97d4ba92b5f", nil
	case Tiny_en:
		return "c78c86eb1a8faa21b369bcd33207cc90d64ae9df", nil
	case Base:
		return "465707469ff3a37a2b9b8d8f89f2f99de7299dac", nil
	case Base_en:
		return "137c40403d78fd54d454da0f9bd998f78703390c", nil
	case Small:
		return "55356645c2b361a969dfd0ef2c5a50d530afd8d5", nil
	case Small_en:
		return "db8a495a91d927739e50b3fc1cc4c6b8f6c2d022", nil
	case Medium:
		return "fd9727b6e1217c2f614f9b698455c4ffd82463b4", nil
	case Medium_en:
		return "8c30f0e44ce9560643ebd10bbe50cd20eafd3723", nil
	case Large_v1:
		return "b1caaf735c4cc1429223d5a74f0f4d0b9b59a299", nil
	case Large:
		return "0f4c8e34f21cf1a914c59d8b3ce882345ad349d6", nil
	}
	return "", errors.New("invalid model")
}

func GetModelSha256Hash(model Model) (string, error) {
	switch model {
	case Tiny:
		return "be07e048e1e599ad46341c8d2a135645097a538221678b7acdd1b1919c6e1b21", nil
	case Tiny_en:
		return "921e4cf8686fdd993dcd081a5da5b6c365bfde1162e72b08d75ac75289920b1f", nil
	case Base:
		return "60ed5bc3dd14eea856493d334349b405782ddcaf0028d4b5df4088345fba2efe", nil
	case Base_en:
		return "a03779c86df3323075f5e796cb2ce5029f00ec8869eee3fdfb897afe36c6d002", nil
	case Small:
		return "1be3a9b2063867b937e64e2ec7483364a79917e157fa98c5d94b5c1fffea987b", nil
	case Small_en:
		return "c6138d6d58ecc8322097e0f987c32f1be8bb0a18532a3f88f734d1bbf9c41e5d", nil
	case Medium:
		return "6c14d5adee5f86394037b4e4e8b59f1673b6cee10e3cf0b11bbdbee79c156208", nil
	case Medium_en:
		return "cc37e93478338ec7700281a7ac30a10128929eb8f427dda2e865faa8f6da4356", nil
	case Large_v1:
		return "7d99f41a10525d0206bddadd86760181fa920438b6b33237e3118ff6c83bb53d", nil
	case Large:
		return "9a423fe4d40c82774b6af34115b8b935f34152246eb19e80e376071d3f999487", nil
	}
	return "", errors.New("invalid model")
}
