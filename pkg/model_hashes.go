package pkg

func getHash(model Model) string {
	switch model {
	case tiny:
		return "bd577a113a864445d4c299885e0cb97d4ba92b5f"
	case tiny_en:
		return "c78c86eb1a8faa21b369bcd33207cc90d64ae9df"
	case base:
		return "465707469ff3a37a2b9b8d8f89f2f99de7299dac"
	case base_en:
		return "137c40403d78fd54d454da0f9bd998f78703390c"
	case small:
		return "55356645c2b361a969dfd0ef2c5a50d530afd8d5"
	case small_en:
		return "db8a495a91d927739e50b3fc1cc4c6b8f6c2d022"
	case medium:
		return "fd9727b6e1217c2f614f9b698455c4ffd82463b4"
	case medium_en:
		return "8c30f0e44ce9560643ebd10bbe50cd20eafd3723"
	case large_v1:
		return "b1caaf735c4cc1429223d5a74f0f4d0b9b59a299"
	case large:
		return "0f4c8e34f21cf1a914c59d8b3ce882345ad349d6"
	}
	return ""
}
